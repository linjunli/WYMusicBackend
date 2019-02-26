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
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	PresentKey = "0CoJUm6Qyw8W8jud"
	iv = []byte("0102030405060708")
	base62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	linuxapikey = "rFgB&h#%2?^eDg:Q"
	publicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDgtQn2JZ34ZC28NWYpAUd98iZ37BUrX/aKzmFbt7clFSs6sXqHauqKWqdtLkF2KexO40H1YTX8z2lSgBBOAxLsvaklV8k4cBFK9snQXE9/DDaFt6Rr7iVZMldczhC0JNgTz+SHXT6CBHuX3e9SdB1Ua44oncaTWz7OBGLbCiK45wIDAQAB
-----END PUBLIC KEY-----`
)

type Request struct {
	Method string
	Url string
	Data string
}

func NewRequest() *Request {
	return &Request{}
}

type LinuxApiParams struct {
	Url string `json:"url"`
	Method string `json:"method"`
	Params map[string]interface{} `json:"params"`
}

func PostReq(Url string,Data []byte , mode string) ([]byte,error) {
	client := http.Client{}

	req,err := http.NewRequest("POST",Url, nil)
	if err != nil {
		fmt.Println("new req err :" + err.Error())
		return nil,err
	}
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	req.Header.Set("User-Agent",SelectOneUa())
	if strings.Contains(Url,"music.163.com") {
		req.Header.Set("Referer","https://music.163.com")
	}

	q := req.URL.Query()

	// 处理url
	var matchUrl string
	reg, _ := regexp.Compile("\\w*api")

	if mode == "weapi" {
		params,encSecKey := WeApi(string(Data))

		q.Add("params",params)
		q.Add("encSecKey",encSecKey)


		if reg.MatchString(Url) {
			matchUrl = reg.ReplaceAllString(Url,"weapi")
		} else {
			matchUrl = Url
		}

		fmt.Println("we api raw query : "+q.Encode())
	} else if mode == "linuxapi" {

		if reg.MatchString(Url) {
			matchUrl = reg.ReplaceAllString(Url,"api")
		} else {
			matchUrl = Url
		}

		var v *map[string]interface{}

		err = json.Unmarshal(Data,&v)

		if err != nil {
			fmt.Println("linux api json unmarshal err :" + err.Error())
			return nil,err
		}

		linuxApiParams := &LinuxApiParams{
			Url:matchUrl,
			Method:"POST",
			Params: *v, // params 是一个key value 对象
		}
		// 需要设置linux下的ua
		req.Header.Set("User-Agent","Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36")
		matchUrl = "https://music.163.com/api/linux/forward"
		paramStr, _ := json.Marshal(linuxApiParams)
		params := LinuxApi(string(paramStr))

		q.Add("eparams",params)

		fmt.Println("linux api raw query : "+q.Encode())
	}

	req.URL, _ = url.Parse(matchUrl)
	req.URL.RawQuery = q.Encode()


	fmt.Printf("req url : %s",req.URL.String())

	res,err := client.Do(req)
	if err != nil {
		fmt.Println("client err :" + err.Error())
		return nil,err
	}

	defer res.Body.Close()

	body,err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("read body err :" + err.Error())
		return nil,err
	}

	fmt.Println("res body : " + string(body))
	return body,nil
}

func WeApi(data string)(string, string) {
	secretKey := "1234567890123456"
	fmt.Println(Reverse(secretKey))

	params := base64.StdEncoding.EncodeToString(AesEncrypt(base64.StdEncoding.EncodeToString(AesEncrypt(data,PresentKey)),secretKey))
	encSecKey := Bin2hex(RsaEncrypt(Reverse(secretKey)))

	result := fmt.Sprintf(`params=%s,"encSecKey":%s}`,params,encSecKey)


	fmt.Println("we api encrypt data :" + result)
	return params,encSecKey
}

func LinuxApi(data string) string  {
	fmt.Printf("\nlinux data: %s\n",data)
	result := AesEcbEncrypt(data, linuxapikey)
	// 转为16进制
	hexResult := hex.EncodeToString(result)

	upperResult := strings.ToUpper(hexResult)
	return upperResult
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

	c := new(big.Int).SetBytes(enc.Bytes())
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

