package snapshot

import (
	"crypto"
	"crypto/rand"
	_ "crypto/sha256"
	"encoding/base64"

	"github.com/gogo/protobuf/proto"
)

/* ProofHashFunc stores the crypto.Hash that is used for transaction
proofing and verification */
var ProofHashFunc crypto.Hash = crypto.SHA256

// Interface for any datatype that can be Marshaled
type marshaler interface {
	Marshal() ([]byte, error)
}

// SimpleProofTuple implements a nodes transaction proof record
type SimpleProofTuple struct {
	protoProofTuple *Snapshot_ProofTuple
}

/* NewSimpleProofTuple instantiates a new SimpleProofTuple with
the given attributes */
func NewSimpleProofTuple(tx *SimpleTransaction, id string, epoch int32, balance float64, signer crypto.Signer) (*SimpleProofTuple, error) {
	tHashed, err := digestMarshaler(tx)
	if err != nil {
		return nil, &DigestErr{simpleErr{err: err, msg: "NewSimpleProofTuple() on Transaction"}}
	}
	transactionSign, err := signer.Sign(rand.Reader, tHashed, ProofHashFunc)
	if err != nil {
		return nil, &SignatureErr{simpleErr{err: err, msg: "NewSimpleProofTuple() on Transaction"}}
	}

	EpochTriplet := NewSimpleEpochTriplet(id, epoch, balance)
	eHashed, err := digestMarshaler(EpochTriplet)
	if err != nil {
		return nil, &DigestErr{simpleErr{err: err, msg: "NewSimpleProofTuple() on Epoch"}}
	}
	epochSign, err := signer.Sign(rand.Reader, eHashed, ProofHashFunc)
	if err != nil {
		return nil, &SignatureErr{simpleErr{err: err, msg: "NewSimpleProofTuple() on Epoch"}}
	}

	b64TransactionSign := base64.StdEncoding.EncodeToString(transactionSign)
	b64EpochSign := base64.StdEncoding.EncodeToString(epochSign)
	return &SimpleProofTuple{
		protoProofTuple: &Snapshot_ProofTuple{
			Epoch:           EpochTriplet.protoEpochTriplet,
			TransactionSign: b64TransactionSign,
			EpochSign:       b64EpochSign,
		},
	}, nil
}

/* digestMarshaler hashes the marshaler with the ProofHashFunc and
slices the digest down to the size specified by the hash function.
Cutting the digest down shouldn't have to be done but for some reason
the crypto.Hash SHA256 function returns 42 bytes instead of 32 like
using the crypto/256 library directly does */
func digestMarshaler(m marshaler) ([]byte, error) {
	serial, err := m.Marshal()
	if err != nil {
		return nil, &MarshalErr{simpleErr{err: err, msg: "NewSimpleProofTuple()"}}
	}

	hasher := ProofHashFunc.New()
	hasher.Reset()
	tHashed := hasher.Sum(serial)
	return tHashed[:hasher.Size()], nil
}

/* GetTransactionSignature returns the signature of a transaction
created using a nodes private key */
func (sp *SimpleProofTuple) GetTransactionSignature() string {
	return sp.protoProofTuple.TransactionSign
}

/* GetEpochSignature returns the signature of an epoch created using
a nodes private key */
func (sp *SimpleProofTuple) GetEpochSignature() string {
	return sp.protoProofTuple.EpochSign
}

// GetEpoch returns the SimpleEpochTriplet embedded in SimpleProofTuple
func (sp *SimpleProofTuple) GetEpoch() *SimpleEpochTriplet {
	return &SimpleEpochTriplet{
		protoEpochTriplet: sp.protoProofTuple.Epoch,
	}
}

/* SimpleEpochTriplet implements a (Node Id, Epoch Number) pair used to
describe the number of transactions `Node Id` has been involved with */
type SimpleEpochTriplet struct {
	protoEpochTriplet *Snapshot_ProofTuple_EpochTriplet
}

/* NewSimpleEpochTriplet instantiates a new SimpleEpochTriplet with the
given attributes */
func NewSimpleEpochTriplet(id string, epoch int32, balance float64) *SimpleEpochTriplet {
	return &SimpleEpochTriplet{
		protoEpochTriplet: &Snapshot_ProofTuple_EpochTriplet{
			Id:      id,
			Epoch:   epoch,
			Balance: balance,
		},
	}
}

func (se *SimpleEpochTriplet) Sign(signer crypto.Signer) ([]byte, error) {
	digest, err := digestMarshaler(se)
	if err != nil {
		return nil, &DigestErr{simpleErr{err: err, msg: "SimpleEpochTriplet.Sign()"}}
	}
	sig, err := signer.Sign(rand.Reader, digest, ProofHashFunc)
	if err != nil {
		return nil, &SignatureErr{simpleErr{err: err, msg: "SimpleEpochTriplet.Sign()"}}
	}
	return sig, nil
}

// Marhsal serializes SimpleEpochTriplet into a slice of bytes
func (se *SimpleEpochTriplet) Marshal() ([]byte, error) {
	out, err := proto.Marshal(se.protoEpochTriplet)
	if err != nil {
		return out, &MarshalErr{simpleErr{err: err, msg: "SimpleEpochTriplet.Marshal()"}}
	}
	return out, nil
}

// Unmarshal deserializes SimpleEpochTriplet from a slice of bytes
func (se *SimpleEpochTriplet) Unmarshal(serial []byte) error {
	se.protoEpochTriplet = &Snapshot_ProofTuple_EpochTriplet{}
	if err := proto.Unmarshal(serial, se.protoEpochTriplet); err != nil {
		return &MarshalErr{simpleErr{err: err, msg: "SimpleEpochTriplet.Unmarshal()"}}
	}
	return nil
}

// GetId returns a Node Id portion of the epoch pair
func (se *SimpleEpochTriplet) GetId() string {
	return se.protoEpochTriplet.GetId()
}

// GetEpochNumber returns the epoch number as an int32
func (se *SimpleEpochTriplet) GetEpochNumber() int32 {
	return se.protoEpochTriplet.GetEpoch()
}

// GetBalance returns the balance of the node associated with this epoch
func (se *SimpleEpochTriplet) GetBalance() float64 {
	return se.protoEpochTriplet.GetBalance()
}
