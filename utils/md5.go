package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Md5Encrypt(a string) string {
	data := []byte(a)
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func TokenDecrypt(token string) int {
	plaintext, _ := hex.DecodeString(token)
	key, _ := hex.DecodeString("e2f068d799c1112c529e140df713329a")
	iv := key[:aes.BlockSize]
	str1, _ := AesDecrypt(plaintext, key, iv)
	str := string(str1)
	if str == "" {
		return -1
	}
	userid, _ := strconv.Atoi(GetBetweenStr(str, "id=", "&"))
	return userid
}

func DyidDecrypt(dysign string, dyid string, user string) int {
	plaintext, _ := hex.DecodeString(dysign)
	key, _ := hex.DecodeString("69f4b74a1fc88d9a44f753bdddf3cda9")
	iv := key[:aes.BlockSize]
	str1, _ := AesDecrypt(plaintext, key, iv)
	str := string(str1)
	fmt.Println(str)
	if str == "" {
		return -1
	}

	yztime, _ := strconv.Atoi(str[:13])
	nowtime := time.Now().UnixNano() / 1e6
	if nowtime-int64(yztime) > 86400000 {
		return -101
	}
	yz := str[15+yztime%10 : 15+yztime%10+1]
	if yz != "1" {
		return -102
	}

	bs, _ := strconv.Atoi(str[13:15])
	if str[25:len(str)-bs] != dyid {
		return -103
	}
	if str[len(str)-bs:] != user {
		return -104
	}

	return 1
}

func AesDecrypt(ciphertext []byte, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(ciphertext))
	defer func() {
		if err := recover(); err != nil {
		}

	}()
	blockMode.CryptBlocks(origData, ciphertext)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AesEncrypt 加密函数
func AesEncrypt(plaintext []byte, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	plaintext = PKCS7Padding(plaintext, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(plaintext))
	blockMode.CryptBlocks(crypted, plaintext)
	return crypted, nil
}

func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n+len(start):])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

// GetToken 获取token
func GetToken(userid int) (string, error) {
	return "TODO", nil
}
