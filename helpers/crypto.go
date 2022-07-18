package helpers

import (
	"crypto/md5"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func GetMD5Hash(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

func GetCryptoPassword(text string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hashedPassword), nil
}

func CheckPassword(hash string, password string) bool {
	h, err := hex.DecodeString(hash)
	if err == nil {
		err = bcrypt.CompareHashAndPassword(h, []byte(password))
		if err != nil {
			return false
		}
	}
	return true
}
