package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"
)


func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// 自定义大数运算
func pow(a *big.Int, b int) *big.Int {
	var i int
	val := big.NewInt(1)
	for i = 0; i < b; i++ {
		val.Mul(val, a)
		fmt.Println(val)
	}
	return val.Rem(val, big.NewInt(21))
}

func gcd(a int64, b int64) int64 {
	if b > 0 {
		return gcd(b, a%b)
	} else {
		return a
	}
}

/**
通用预处理请求
 */
func GetRequestBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	_ = r.ParseForm()
	w.Header().Set("Content-Type","application/json")

	if strings.ToUpper(r.Method) != "POST" {
		ResultError(w,"request method must be post",-100)
		return nil, errors.New("request method error")
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		ResultError(w,err.Error(),-2000)
		return nil,err
	}

	log.Printf("rcv req : %s",string(b))

	return b,nil
}


/**
通用返回错误方式
 */
func ResultError(w http.ResponseWriter, msg string, code int)  {
	type Result struct {
		Code int `json:"code"`
		Msg string `json:"msg"`
	}
	w.Header().Set("Content-Type","application/json")
	result := &Result{
		Code:code,
		Msg:msg,
	}

	b, _ := json.Marshal(result)

	log.Printf("fatal err occurd:%s", string(b))

	_,_ = w.Write(b)
}