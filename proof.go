package snapshot

import (
	"crypto"
	"crypto/rand"
	_ "crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/gogo/protobuf/proto"
)

/* ProofHashFunc stores the crypto.Hash that is used for transaction
proofing and verification */
var ProofHashFunc crypto.Hash = crypto.SHA256

// Interface for any datatype that can be Marshaled
type Marshaler interface {
	Marshal() ([]byte, error)
}

// SimpleProofTuple implements a nodes transaction proof record
type SimpleProofTuple struct {
	protoProofTuple *Snapshot_ProofTuple
}

/* NewSimpleProofTuple instantiates a new SimpleProofTuple with
the given attributes */
func NewSimpleProofTuple(tx *SimpleTransaction, id string, epoch int32, signer crypto.Signer) (*SimpleProofTuple, error) {
	tHashed, err := digestMarshaler(tx)
	if err != nil {
		return nil, fmt.Errorf("Digest fail in NewSimpleProofTuple(): %v", err)
	}
	transactionSign, err := signer.Sign(rand.Reader, tHashed, ProofHashFunc)
	if err != nil {
		return nil, fmt.Errorf("Transaction Sign failed in NewSimpleProofTuple(): %v", err)
	}

	epochPair := NewSimpleEpochPair(id, epoch)
	eHashed, err := digestMarshaler(epochPair)
	if err != nil {
		return nil, fmt.Errorf("Epoch digest fail in NewSimpleProofTuple(): %v", err)
	}
	epochSign, err := signer.Sign(rand.Reader, eHashed, ProofHashFunc)
	if err != nil {
		return nil, fmt.Errorf("Epoch Sign failed in NewSimpleProofTuple(): %v", err)
	}

	b64TransactionSign := base64.StdEncoding.EncodeToString(transactionSign)
	b64EpochSign := base64.StdEncoding.EncodeToString(epochSign)
	return &SimpleProofTuple{
		protoProofTuple: &Snapshot_ProofTuple{
			Epoch:           epochPair.protoEpochPair,
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
func digestMarshaler(m Marshaler) ([]byte, error) {
	serial, err := m.Marshal()
	if err != nil {
		return nil, fmt.Errorf("Transaction Marshal fail in NewSimpleProofTuple(): %v", err)
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

// GetEpoch returns the SimpleEpochPair embedded in SimpleProofTuple
func (sp *SimpleProofTuple) GetEpoch() *SimpleEpochPair {
	return &SimpleEpochPair{
		protoEpochPair: sp.protoProofTuple.Epoch,
	}
}

/* SimpleEpochPair implements a (Node Id, Epoch Number) pair used to
describe the number of transactions `Node Id` has been involved with */
type SimpleEpochPair struct {
	protoEpochPair *Snapshot_ProofTuple_EpochPair
}

/* NewSimpleEpochPair instantiates a new SimpleEpochPair with the
given attributes */
func NewSimpleEpochPair(id string, epoch int32) *SimpleEpochPair {
	return &SimpleEpochPair{
		protoEpochPair: &Snapshot_ProofTuple_EpochPair{
			Id:    id,
			Epoch: epoch,
		},
	}
}

// Marhsal serializes SimpleEpochPair into a slice of bytes
func (se *SimpleEpochPair) Marshal() ([]byte, error) {
	out, err := proto.Marshal(se.protoEpochPair)
	if err != nil {
		return out, fmt.Errorf("Marshal fail in SimpleEpochPair.Marshal(): %v", err)
	}
	return out, nil
}

// Unmarshal deserializes SimpleEpochPair from a slice of bytes
func (se *SimpleEpochPair) Unmarshal(serial []byte) error {
	se.protoEpochPair = &Snapshot_ProofTuple_EpochPair{}
	if err := proto.Unmarshal(serial, se.protoEpochPair); err != nil {
		return fmt.Errorf("Unmarshal fail in SimpleEpochPair.Unmarshal(): %v", err)
	}
	return nil
}

// GetId returns a Node Id portion of the epoch pair
func (se *SimpleEpochPair) GetId() string {
	return se.protoEpochPair.GetId()
}

// GetEpochNumber returns the epoch number as an int32
func (se *SimpleEpochPair) GetEpochNumber() int32 {
	return se.protoEpochPair.GetEpoch()
}
