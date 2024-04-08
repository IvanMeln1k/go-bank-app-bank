package hasher

import (
	"crypto/sha256"
	"fmt"
)

type HasherInterface interface {
	Hash(password string) string
	Check(password string, hash string) bool
}

type Hasher struct {
	salt string
}

func NewHasher(salt string) *Hasher {
	return &Hasher{
		salt: salt,
	}
}

func (h *Hasher) Hash(password string) string {
	bytes := sha256.Sum256([]byte(password + h.salt))
	return fmt.Sprintf("%x", bytes)
}

func (h *Hasher) Check(password string, hash string) bool {
	return h.Hash(password) == hash
}
