package snapshot

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"strconv"
	"testing"
)

//TRANSACTION
func TestTransaction(t *testing.T) {
	totalTests := 1000
	failedTests := 0

	for i := 0; i < totalTests; i++ {
		code := int32(i * 10)
		reward := float64(i) / 10
		exchange := float64(i) / 2
		tx := createTransaction(code, reward, exchange, "ID1", "ID2")

		raw, err := tx.Marshal()
		if err != nil {
			panic(err)
		}

		tx2 := SimpleTransaction{}
		err = tx2.Unmarshal(raw)
		if err != nil {
			panic(err)
		}
		raw2, err := tx2.Marshal()
		if err != nil {
			panic(err)
		}

		if !bytes.Equal(raw, raw2) {
			failedTests++
			fmt.Printf("Failed with inputs: code(%d) reward(%f) exchange(%f)\n", code, reward, exchange)
		}
	}
	fmt.Printf("Passed %d/%d tests\n", totalTests-failedTests, totalTests)
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

func createTransaction(a int32, r float64, e float64, g string, l string) *SimpleTransaction {
	tx := NewSimpleTransaction()
	tx.SetActionCode(a)
	tx.SetBystanderReward(r)
	tx.SetValueExchange(e)
	tx.SetGainingParty(g)
	tx.SetLosingParty(l)
	return tx
}

//SNAPSHOT
func TestSnapshot(t *testing.T) {
	totalKeys := 8
	keys := make(map[string]*rsa.PrivateKey)
	for i := 0; i < totalKeys; i++ {
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			panic(err)
		}
		keys[strconv.Itoa(i)] = key
	}

	totalTests := 100
	totalFails := 0
	passReq := 0.666
	for i := 0; i < totalTests; i++ {
		code := int32(i * 10)
		reward := float64(i) / 10
		exchange := float64(i) / 2
		tx := createTransaction(code, reward, exchange, "ID1", "ID2")

		snapshot := NewSimpleSnapshot(tx)

		count := 0
		for id, key := range keys {
			tup, _ := NewSimpleProofTuple(tx, id, 1, key)
			if count%3 == 0 {
				tup.protoProofTuple.EpochSign = ""
			}
			snapshot.AddProof(tup)
			count++
		}

		pubKeys := make(map[string]crypto.PublicKey)
		for id, key := range keys {
			pubKeys[id] = &key.PublicKey
		}
		err := VerifySnapshot(passReq, snapshot, pubKeys, pkcsVerifier)
		if err != nil {
			totalFails++
			fmt.Printf("Failed with inputs: code(%d) reward(%f) exchange(%f)\n", code, reward, exchange)
			fmt.Printf("Error was: %v\n", err)
		}
	}

	fmt.Printf("Passed %d/%d tests\n", totalTests-totalFails, totalTests)
}

func pkcsVerifier(key crypto.PublicKey, hash crypto.Hash, digest []byte, sig []byte) error {
	rsaKey := (key).(*rsa.PublicKey)
	return rsa.VerifyPKCS1v15(rsaKey, hash, digest, sig)
}
