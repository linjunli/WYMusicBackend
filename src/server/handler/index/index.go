package index

import (
	"WYMusicBackend/src/netnease"
	"WYMusicBackend/src/util"
	"encoding/json"
	"net/http"
)

type ApiGetInfo struct {

}

type Response struct {
	Ret int `json:"ret"`
	Msg string `json:"msg"`
	PlayList *netnease.PlayListRes `json:"play_list"`
}

func (i *ApiGetInfo)ServeHTTP(w http.ResponseWriter, r *http.Request)  {

	_,err := util.GetRequestBody(w,r)

	if err != nil {
		util.ResultError(w,err.Error(),-1)
		return
	}

	playList := &netnease.PlayList{
		Cat: "全部",
		Order: "hot",
		Limit: 6,
		Offset: 0,
	}

	playListRes, err := playList.GetPlayList()

	if err != nil {
		util.ResultError(w,err.Error(),-2)
		return
	}


	response := &Response{
		Ret: 0,
		Msg: "ok",
		PlayList: playListRes,
	}

	m,err := json.Marshal(response)

	if err !=nil {
		util.ResultError(w,err.Error(),-3)
		return
	}

	_,_ = w.Write(m)
}