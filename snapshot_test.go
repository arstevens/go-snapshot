package snapshot

import (
	"fmt"
	"testing"
)

func TestTransaction(t *testing.T) {
	tx := createTransaction(1, 0.01, 1.2, "ID1", "ID2")
	fmt.Println("Original ->")
	printTransaction(tx)

	raw, err := tx.Marshal()
	if err != nil {
		panic(err)
	}

	tx2 := SimpleTransaction{}
	err = tx2.Unmarshal(raw)
	if err != nil {
		panic(err)
	}

	fmt.Println("Unmarshaled ->")
	printTransaction(tx2)
}

func printTransaction(tx SimpleTransaction) {
	fmt.Println("TRANSACTION START")
	fmt.Printf("%d %f %f %s %s\n", tx.GetActionCode(), tx.GetBystanderReward(),
		tx.GetValueExchange(), tx.GetGainingParty(), tx.GetLosingParty())
	for _, bystander := range tx.GetBystanders() {
		fmt.Printf("%s ", bystander)
	}
	fmt.Println("\nTRANSACTION END")
}

func createTransaction(a int32, r float64, e float64, g string, l string) SimpleTransaction {
	tx := SimpleTransaction{}
	tx.SetActionCode(a)
	tx.SetBystanderReward(r)
	tx.SetValueExchange(e)
	tx.SetGainingParty(g)
	tx.SetLosingParty(l)
	return tx
}
