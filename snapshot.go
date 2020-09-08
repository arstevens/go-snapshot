package snapshot

import (
	"google.golang.org/protobuf/proto"
)

/* SimpleSnapshot implements a collection of post-transaction
verification data */
type SimpleSnapshot struct {
	protoSnapshot *Snapshot
}

// NewSimpleSnapshot returns an empty instance of a SimpleSnapshot
func NewSimpleSnapshot(tx *SimpleTransaction) *SimpleSnapshot {
	return &SimpleSnapshot{
		protoSnapshot: &Snapshot{Transaction: tx.protoTransaction},
	}
}

// Marshal serializes a SimpleSnapshot into a slice of bytes
func (ss *SimpleSnapshot) Marshal() ([]byte, error) {
	out, err := proto.Marshal(ss.protoSnapshot)
	if err != nil {
		return out, &MarshalErr{simpleErr{err: err, msg: "SimpleSnapshot.Marshal()"}}
	}
	return out, nil
}

// Unmarshal deserializes a slice of bytes into a SimpleSnapshot
func (ss *SimpleSnapshot) Unmarshal(serial []byte) error {
	ss.protoSnapshot = &Snapshot{}
	if err := proto.Unmarshal(serial, ss.protoSnapshot); err != nil {
		return &MarshalErr{simpleErr{err: err, msg: "SimpleSnapshot.Unmarshal()"}}
	}
	return nil
}

/* GetTransaction returns a pointer to a SimpleTransaction object.
All calls to GetTransaction() will return an object that points to
the same underlying data. A change to the SimpleTransaction object
will change all SimpleTransaction objects created with this method */
func (ss *SimpleSnapshot) GetTransaction() *SimpleTransaction {
	return &SimpleTransaction{
		protoTransaction: ss.protoSnapshot.GetTransaction(),
	}
}

/* AddProof adds a SimpleProofTuple to the SimpleSnapshot for later
verification */
func (ss *SimpleSnapshot) AddProof(proof *SimpleProofTuple) {
	proofs := ss.protoSnapshot.GetProofs()
	proofs = append(proofs, proof.protoProofTuple)
	ss.protoSnapshot.Proofs = proofs
}

/* GetProofs returns all SimpleProofTuples currently held by
SimpleSnapshot in the form of a slice */
func (ss *SimpleSnapshot) GetProofs() []*SimpleProofTuple {
	proofs := ss.protoSnapshot.GetProofs()
	simpleProofs := make([]*SimpleProofTuple, 0, len(proofs))
	for _, proof := range proofs {
		simpleProofs = append(simpleProofs, &SimpleProofTuple{protoProofTuple: proof})
	}
	return simpleProofs
}
