package wallit

import (
	"github.com/mit-dci/lit/btcutil"
	"github.com/mit-dci/lit/crypto/koblitz"
	"github.com/mit-dci/lit/logging"
	"github.com/mit-dci/lit/portxo"
)

/*
Key derivation for a TxStore has 3 levels: use case, peer index, and keyindex.
Regular wallet addresses are use 0, peer 0, and then a linear index.
The identity key is use 11, peer 0, index 0.
Channel multisig keys are use 2, peer and index per peer and channel.
Channel refund keys are use 3, peer and index per peer / channel.
*/

// =====================================================================
// OK only use these now

// PathPrivkey returns a private key by descending the given path
// Returns nil if there's an error.
func (w *Wallit) PathPrivkey(kg portxo.KeyGen) *koblitz.PrivateKey {
	// in uspv, we require path depth of 5
	if kg.Depth != 5 {
		return nil
	}
	priv, err := kg.DerivePrivateKey(w.rootPrivKey)
	if err != nil {
		logging.Errorf("PathPrivkey err %s", err.Error())
		return nil
	}
	return priv
}

// PathPubkey returns a public key by descending the given path.
// Returns nil if there's an error.
func (w *Wallit) PathPubkey(kg portxo.KeyGen) *koblitz.PublicKey {
	priv := w.PathPrivkey(kg)
	if priv == nil {
		return nil
	}
	return w.PathPrivkey(kg).PubKey()
}

// PathPubHash160 returns a 20 byte pubkey hash for the given path
// It'll always return 20 bytes, or a nil if there's an error.
func (w *Wallit) PathPubHash160(kg portxo.KeyGen) [20]byte {
	var pkh [20]byte
	pub := w.PathPubkey(kg)
	if pub == nil {
		return pkh
	}
	copy(pkh[:], btcutil.Hash160(pub.SerializeCompressed()))

	return pkh
}

// ------------- end of 2 main key deriv functions

// get a private key from the regular wallet
func (w *Wallit) GetWalletPrivkey(idx uint32) *koblitz.PrivateKey {
	var kg portxo.KeyGen
	kg.Depth = 5
	kg.Step[0] = 44 | 1<<31
	kg.Step[1] = w.Param.HDCoinType | 1<<31
	kg.Step[2] = 0 | 1<<31
	kg.Step[3] = 0 | 1<<31
	kg.Step[4] = idx | 1<<31
	return w.PathPrivkey(kg)
}

// GetWalletKeygen returns the keygen for a standard wallet address
func GetWalletKeygen(idx, cointype uint32) portxo.KeyGen {
	var kg portxo.KeyGen
	kg.Depth = 5
	kg.Step[0] = 44 | 1<<31
	kg.Step[1] = cointype | 1<<31
	kg.Step[2] = 0 | 1<<31
	kg.Step[3] = 0 | 1<<31
	kg.Step[4] = idx | 1<<31
	return kg
}

// GetUsePrive generates a private key for the given use case & keypath
func (w *Wallit) GetUsePriv(kg portxo.KeyGen, use uint32) *koblitz.PrivateKey {
	kg.Step[2] = use
	return w.PathPrivkey(kg)
}

// GetUsePub generates a pubkey for the given use case & keypath
func (w *Wallit) GetUsePub(kg portxo.KeyGen, use uint32) [33]byte {
	var b [33]byte
	pub := w.GetUsePriv(kg, use).PubKey()
	if pub != nil {
		copy(b[:], pub.SerializeCompressed())
	}
	return b
}
