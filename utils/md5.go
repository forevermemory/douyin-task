package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encrypt(a string) string {
	data := []byte(a)
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
