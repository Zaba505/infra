package backend

import "fmt"

type ObjectReadError struct {
	Cause error
}

func (e ObjectReadError) Error() string {
	return fmt.Sprintf("failed to read object: %s", e.Cause)
}

func (e ObjectReadError) Unwrap() error {
	return e.Cause
}

type ChecksumMismatchError struct{}

func (e ChecksumMismatchError) Error() string {
	return "checksum mismatch"
}
