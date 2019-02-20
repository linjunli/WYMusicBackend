package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
)

var (
	PresentKey = "0CoJUm6Qyw8W8jud"
	iv = []byte("0102030405060708")
	base62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	publicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDgtQn2JZ34ZC28NWYpAUd98iZ37BUrX/aKzmFbt7clFSs6sXqHauqKWqdtLkF2KexO40H1YTX8z2lSgBBOAxLsvaklV8k4cBFK9snQXE9/DDaFt6Rr7iVZMldczhC0JNgTz+SHXT6CBHuX3e9SdB1Ua44oncaTWz7OBGLbCiK45wIDAQAB
-----END PUBLIC KEY-----`
)

func WeApi(data string) (string, string) {
	secretKey := "1234567890123456"
	fmt.Println(Reverse(secretKey))

	params := base64.StdEncoding.EncodeToString(AesEncrypt(base64.StdEncoding.EncodeToString(AesEncrypt(data,PresentKey)),secretKey))
	encSecKey := Bin2hex(RsaEncrypt(Reverse(secretKey)))

	return params,encSecKey
}

func AesEncrypt(msg string, key string) []byte {
	theKey := []byte(key)
	block, _ := aes.NewCipher(theKey)

	mode := cipher.NewCBCEncrypter(block,iv)

	theRaw := PKCS5Padding([]byte(msg),block.BlockSize())
	result := make([]byte,len(theRaw))
	mode.CryptBlocks(result, theRaw)

	return result
}

func RsaEncrypt(msg string) []byte  {
	block, _ := pem.Decode([]byte(publicKey))
	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)

	pub := pubInterface.(*rsa.PublicKey)

	msgLen :=len([]byte(msg))

	enc := new(bytes.Buffer)
	_ = binary.Write(enc, binary.BigEndian, make([]byte, 128 - msgLen))
	_ = binary.Write(enc, binary.BigEndian, []byte(msg))

	com := BytesCombine(make([]byte,128 - msgLen),[]byte(msg))

	c := new(big.Int).SetBytes(com)
	return c.Exp(c, big.NewInt(int64(pub.E)), pub.N).Bytes()
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Unpadding(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func Hex2bin(raw string) []byte  {
	result, _ := hex.DecodeString(raw)
	return result
}

func Bin2hex(raw []byte) string {
	return hex.EncodeToString(raw)
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func BytesCombine(Bytes ...[]byte) []byte  {
	return bytes.Join(Bytes, []byte(""))
}