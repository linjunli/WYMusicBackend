package playlist

import (
	"WYMusicBackend/src/netnease"
	"WYMusicBackend/src/util"
	"encoding/json"
	"net/http"
)

type ApiGetPlayListDetail struct {

}

type request struct {
	N int `json:"n"`
	S int `json:"s"`
	Id int `json:"id"`
}

type response struct {
	Ret int `json:"ret"`
	Msg string `json:"msg"`
	Data *netnease.PlayListDetailRes `json:"data"`
}

func (i *ApiGetPlayListDetail)ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	b, err := util.GetRequestBody(w,r)

	if err != nil {
		util.ResultError(w,err.Error(),-1)
		return
	}

	requ := &request{}

	err = json.Unmarshal(b,requ)

	if err !=nil {
		util.ResultError(w,err.Error(),-3)
		return
	}

	req := &netnease.PlayListDetailReq{
		N:100000,
		S:8,
	}

	if requ.Id != 0 {
		req.Id = requ.Id
	}
	if requ.S != 0 {
		req.S  = requ.S
	}
	if requ.N != 0 {
		req.N  = requ.N
	}

	playListDetail,err := req.GetPlayListDetail()

	if err != nil {
		util.ResultError(w,err.Error(),-2)
		return
	}

	response := &response{
		Ret:0,
		Msg:"ok",
		Data:playListDetail,
	}

	m,err := json.Marshal(response)

	if err !=nil {
		util.ResultError(w,err.Error(),-3)
		return
	}

	_,_ = w.Write(m)
}