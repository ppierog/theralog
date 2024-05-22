package cred

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

func GenerateSalt() string {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10) // 1436773875771421417
	h := sha256.New()
	h.Write([]byte(timestamp))
	return hex.EncodeToString(h.Sum(nil))[:8]
}
func CalcSha256(password string, salt string) string {
	passwdSalt := password + salt
	h := sha256.New()
	h.Write([]byte(passwdSalt))
	paswdSaltSha256 := hex.EncodeToString(h.Sum(nil))
	return paswdSaltSha256
}
