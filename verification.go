package snapshot

import (
	"crypto"
	_ "crypto/sha256"
	"encoding/base64"
	"fmt"
)

/* Verifier is a function type that can verify data was signed by the
owner of the provided public key */
type Verifier func(key crypto.PublicKey, hash crypto.Hash, digest []byte, sig []byte) error

/* VerifySnapshot returns whether or not the provided SimpleSnapshot is
valid or not. If the percentage of valid SimpleProofTuples is greater
than the pass parameter then VerifySnapshot returns nil, otherwise it
returns an error */
func VerifySnapshot(pass float64, snapshot *SimpleSnapshot, keys map[string]crypto.PublicKey,
	verf Verifier) error {
	tx := snapshot.GetTransaction()
	tDigest, err := digestMarshaler(tx)
	if err != nil {
		return fmt.Errorf("Transaction digest fail in VerifySnapshot(): %v", err)
	}

	totalPasses := 0
	proofs := snapshot.GetProofs()
	for _, proof := range proofs {
		pk := keys[proof.GetEpoch().GetId()]
		err := verifyProofComponents(proof, pk, verf, tDigest)

		if err == nil {
			totalPasses++
		}
	}
	return didPass(pass, totalPasses, len(proofs))
}

/* verifyProofComponents does the heavy lifting for VerifySnapshot by
verifying the individual SimpleProofTuples */
func verifyProofComponents(proof *SimpleProofTuple, pk crypto.PublicKey, verf Verifier,
	transDigest []byte) error {

	/* Don't need error because verification will fail anyway if the signature
	is empty */
	tSig, _ := base64.StdEncoding.DecodeString(proof.GetTransactionSignature())
	err := verf(pk, ProofHashFunc, transDigest, tSig)
	if err != nil {
		return fmt.Errorf("Unable to verify transaction in verifyProofComponents(): %v", err)
	}

	eDigest, err := digestMarshaler(proof.GetEpoch())
	if err != nil {
		return fmt.Errorf("Epoch digest fail in VerifySnapshot(): %v", err)
	}

	eSig, _ := base64.StdEncoding.DecodeString(proof.GetEpochSignature())
	err = verf(pk, ProofHashFunc, eDigest, eSig)
	if err != nil {
		return fmt.Errorf("Unable to verify epoch in verifyProofComponents(): %v", err)
	}
	return nil
}

/* didPass runs the final check for VerifySnapshot to see whether the
percentage of valid proofs was greater than the pass parameter */
func didPass(pass float64, totalPass int, total int) error {
	passStat := float64(totalPass) / float64(total)
	if passStat < pass {
		return fmt.Errorf("Not enough passes in VerifySnapshot. Passes %f", passStat)
	}
	return nil
}
