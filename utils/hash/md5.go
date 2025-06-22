package hasher

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

type md5Hash struct{}

func NewMd5Hash() *md5Hash {
	return &md5Hash{}
}

func (h *md5Hash) Hash(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (h *md5Hash) CheckPassword(hashedPassword, password string) error {
	hashedInput := h.Hash(password)
	if hashedInput != hashedPassword {
		return errors.New("password mismatch")
	}
	return nil
}
