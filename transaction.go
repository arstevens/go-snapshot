package snapshot

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

// SimpleTransaction implements a network transaction record
type SimpleTransaction struct {
	protoTransaction *Transaction
}

/* Marshal serializes the SimpleTransaction into a
slice of bytes and returns an error if unable to marshal */
func (st *SimpleTransaction) Marshal() ([]byte, error) {
	out, err := proto.Marshal(st.protoTransaction)
	if err != nil {
		return out, fmt.Errorf("Marshal Fail in SimpleTransaction.Marshal(): %v", err)
	}

	return out, nil
}

/* Unmarshal deserializes a slice of bytes into a
SimpleTransaction and returns an error if unable to unmarshal */
func (st *SimpleTransaction) Unmarshal(serial []byte) error {
	st.protoTransaction = &Transaction{}
	if err := proto.Unmarshal(serial, st.protoTransaction); err != nil {
		return fmt.Errorf("Unmarshal fail in SimpleTransaction.Unmarshal(): %v", err)
	}
	return nil
}

// Getter for action code
func (st *SimpleTransaction) GetActionCode() int32 {
	return st.protoTransaction.GetAction()
}

// Setter for action code
func (st *SimpleTransaction) SetActionCode(code int32) {
	createIfNil(st)
	st.protoTransaction.Action = code
}

// Getter for bystander reward
func (st *SimpleTransaction) GetBystanderReward() float64 {
	return st.protoTransaction.GetReward()
}

// Setter for bystander reward
func (st *SimpleTransaction) SetBystanderReward(reward float64) {
	createIfNil(st)
	st.protoTransaction.Reward = reward
}

// Getter for value exchange
func (st *SimpleTransaction) GetValueExchange() float64 {
	return st.protoTransaction.GetExchange()
}

// Setter for value exchange
func (st *SimpleTransaction) SetValueExchange(exchange float64) {
	createIfNil(st)
	st.protoTransaction.Exchange = exchange
}

// Getter for gaining party ID
func (st *SimpleTransaction) GetGainingParty() string {
	return st.protoTransaction.GetGainer()
}

// Setter for gaining party ID
func (st *SimpleTransaction) SetGainingParty(gainer string) {
	createIfNil(st)
	st.protoTransaction.Gainer = gainer
}

// Getter for losing party ID
func (st *SimpleTransaction) GetLosingParty() string {
	return st.protoTransaction.GetLoser()
}

// Setter for losing party ID
func (st *SimpleTransaction) SetLosingParty(loser string) {
	createIfNil(st)
	st.protoTransaction.Loser = loser
}

// Getter for bystander IDs
func (st *SimpleTransaction) GetBystanders() []string {
	return st.protoTransaction.GetBystanders()
}

// Setter for bystander IDs
func (st *SimpleTransaction) SetBystanders(bystanders []string) {
	createIfNil(st)
	st.protoTransaction.Bystanders = bystanders
}

func createIfNil(st *SimpleTransaction) {
	if st == nil {
		st = &SimpleTransaction{}
	}
	if st.protoTransaction == nil {
		st.protoTransaction = &Transaction{}
	}
}
