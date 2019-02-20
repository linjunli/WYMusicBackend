package main

import (
	"WYMusicBackend/util"
	"encoding/json"
	"fmt"
	"math/big"
)

type Text struct {
	A int `json:"a"`
}

func main() {
	a := &Text{A:1}

	b,_ := json.Marshal(a)

	text := string(b)

	params,encSec := util.WeApi(text)

	fmt.Printf("params: %s | encSec: %s\n",params,encSec)
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
