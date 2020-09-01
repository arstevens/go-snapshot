package snapshot

import (
	"fmt"
)

// Shared attributes between module error types
type simpleErr struct {
	err error
	msg string
}

// Prints error message
func (se *simpleErr) Error() string {
	return fmt.Errorf(se.msg+": %v", se.err).Error()
}

// Returns the underlying error
func (se *simpleErr) Unwrap() error {
	return se.err
}

// MarshalErr is returned if there was a problem when marshaling
type MarshalErr struct {
	simpleErr
}

// DigestErr is returned if there was a problem when hashing data
type DigestErr struct {
	simpleErr
}

// SignatureErr is returned if there was a problem signing data
type SignatureErr struct {
	simpleErr
}

// VerificationErr is returned if there was a problem verifying a signature
type VerificationErr struct {
	simpleErr
}

// PassErr is returned if VerifySnapshot fails to meet the required success rate
type PassErr struct {
	simpleErr
}
