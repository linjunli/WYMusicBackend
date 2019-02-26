package top

import (
	"WYMusicBackend/src/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PlayList struct {
}

type PlayListResponse struct {
	Code int
	Msg string
	PlayLists []string
}

type RequestPrams struct {
	Cat string `json:"cat"`
	Order string `json:"order"`
	Limit int `json:"limit"`
	Offset int `json:"offset"`
	Total bool `json:"total"`
	CsrfToken string `json:"csrf_token"`
}

func (p *PlayList)ServeHTTP(w http.ResponseWriter, r *http.Request)  {

	w.Header().Set("Content-Type","application/json")

	_ = r.ParseForm()

	b, err := ioutil.ReadAll(r.Body)
	if err!=nil {
		fmt.Println("ioutil err"+err.Error())
		_,_ = w.Write([]byte(err.Error()))
		return
	}
	fmt.Printf("rcv data:%s",string(b))
	req := &RequestPrams{
		Cat: "全部",
		Order: "hot",
		Limit: 50,
		Offset:0,
		Total:true,
		CsrfToken:"",
	}

	err = json.Unmarshal(b,req)
	if err != nil {
		fmt.Println("json unmarshal err "+err.Error())
		_,_ = w.Write([]byte(err.Error()))
		return
	}

	reqStr,err := json.Marshal(req)
	if err != nil {
		fmt.Println(err.Error())
		_,_ = w.Write([]byte(err.Error()))
		return
	}

	res,err := util.PostReq("https://music.163.com/weapi/playlist/list",[]byte(reqStr),"weapi")
	if err != nil {
		fmt.Println("post err "+err.Error())
		_,_ = w.Write([]byte(err.Error()))
		return
	}

	_,err = w.Write(res)
	if err != nil {
		fmt.Println("write err "+err.Error())
		_,_ = w.Write([]byte(err.Error()))
		return
	}
}