package sha

import (
	"crypto/sha256"
	"crypto/sha512"
)

func ComputeSHA256(i []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(i)
	return h.Sum(nil), nil
}

func ComputeSHA512(i []byte) ([]byte, error) {
	h := sha512.New()
	h.Write(i)
	return h.Sum(nil), nil
}
