package transaction

import (
	"errors"
	"fmt"
	"log"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/kittycash/kittiverse/src/kitty"
)

type ID cipher.SHA256

func EmptyID() ID {
	return ID(cipher.SHA256{})
}

func (h ID) Hex() string {
	return cipher.SHA256(h).Hex()
}

type IDs []ID

func (ids IDs) ToStringArray() []string {
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = cipher.SHA256(id).Hex()
	}
	return out
}

type Action func(tx *Transaction) error

// Transaction represents a kitty transaction.
// For IKO, transaction and block are combined to formed one entity.
type Transaction struct {
	KittyID kitty.ID
	In      ID
	Out     cipher.Address
	Sig     cipher.Sig
}

// Meta records meta information for a transaction.
type Meta struct {
	Seq uint64
	TS  int64
}

// Wrapper wraps a transaction with it's meta information.
type Wrapper struct {
	Tx   Transaction
	Meta Meta
}

// NewGenTx creates a "generation" transaction.
// This is where a kitty is created on the blockchain.
func NewGenTx(kittyID kitty.ID, sk cipher.SecKey) *Transaction {
	var (
		address = cipher.AddressFromSecKey(sk)
		tx      = &Transaction{
			KittyID: kittyID,
			In:      EmptyID(),
			Out:     address,
		}
	)
	tx.Sig = tx.Sign(sk)
	return tx
}

// NewTransferTx creates a normal transaction where a kitty is transferred from
// one address to another.
// It returns error when provided secret key does not own input address.
func NewTransferTx(in *Transaction, out cipher.Address, sk cipher.SecKey) (*Transaction, error) {

	// Check input with secret key.
	if expAddr := cipher.AddressFromSecKey(sk); in.Out != expAddr {
		return nil, errors.New("secret key does not own input tx address")
	}

	tx := &Transaction{
		KittyID: in.KittyID,
		In:      in.Hash(),
		Out:     out,
	}
	tx.Sig = tx.Sign(sk)
	return tx, nil
}

func (tx Transaction) Serialize() []byte {
	return encoder.Serialize(tx)
}

func (tx Transaction) Hash() ID {
	return ID(cipher.SumSHA256(tx.Serialize()))
}

func (tx Transaction) HashInner() cipher.SHA256 {
	tx.Sig = cipher.Sig{}
	return cipher.SHA256(tx.Hash())
}

func (tx Transaction) Sign(sk cipher.SecKey) cipher.Sig {
	e := cipher.
		AddressFromSecKey(sk).
		Verify(cipher.PubKeyFromSecKey(sk))
	if e != nil {
		log.Panic(e)
	}
	return cipher.SignHash(tx.HashInner(), sk)
}

// VerifyWith checks the input and signature of the transaction.
//		- Input tx hash.
//		- Tx signature.
// VerifyWith does not check:
//		- Double spending of kitties.
//      - True ownership (as 'Verify' does not know current state).
func (tx Transaction) VerifyWith(in *Transaction, genPK cipher.PubKey) error {
	var isGen = in == nil

	// Check input.
	if isGen == true {
		if exp := EmptyID(); tx.In != exp {
			return fmt.Errorf("generation tx expected 'in:%s', but we got 'in:%s'",
				exp.Hex(), tx.In.Hex())
		}

		// Check signature based on trusted generation public key 'genPK'.
		return cipher.VerifySignature(genPK, tx.Sig, tx.HashInner())

	} else {
		if exp := in.Hash(); tx.In != exp {
			return fmt.Errorf("transfer tx expected 'in:%s', but we got 'in:%s'",
				exp.Hex(), tx.In.Hex())
		}
		// Check kitty.
		if exp := in.KittyID; tx.KittyID != exp {
			return fmt.Errorf("tx expected 'kitty_id:%d', but we got 'kitty_id:%d'",
				exp, tx.KittyID)
		}

		// Check signature based on previous unspent output.
		return cipher.ChkSig(in.Out, tx.HashInner(), tx.Sig)
	}
}

// IsKittyGen returns true if tx is a generation tx:
//		- Tx is of the correct structure to create a new kitty.
//		- Tx is of the right address to create a new kitty.
func (tx Transaction) IsKittyGen(pk cipher.PubKey) bool {
	finish := func(is bool, reason string) bool {
		fmt.Printf("[GENERATION TX?] (%v) %s\n", is, reason)
		return is
	}
	// Check input tx hash is empty.
	if tx.In != EmptyID() {
		return finish(false, "hash is empty")
	}
	// Check output address.
	if e := tx.Out.Verify(pk); e != nil {
		return finish(false, fmt.Sprintf("error: %s", e.Error()))
	}
	// Accept.
	return finish(true, "")
}

// String returns human readable string of transaction.
func (tx Transaction) String() string {
	return fmt.Sprintf("kitty_id:%d|in:%s|out:%s|sig:%s",
		tx.KittyID, tx.In.Hex(), tx.Out.String(), tx.Sig.Hex())
}
