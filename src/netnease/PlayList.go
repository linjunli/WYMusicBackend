package netnease

import (
	"WYMusicBackend/src/util"
	"encoding/json"
)
/**
歌单列表获取
 */
type PlayList struct {
	Cat string `json:"cat"`
	Order string `json:"order"`
	Limit int `json:"limit"`
	Offset int `json:"offset"`
	Total bool `json:"total"`
	CsrfToken string `json:"csrf_token"`
}

type PlayListRes struct {
	Cat string `json:"cat"`
	Code int `json:"code"`
	More bool `json:"more"`
	PlayLists []interface{} `json:"playlists"`
	Total int `json:"total"`
}

func (p *PlayList)GetPlayList() (*PlayListRes, error) {
	str,err := json.Marshal(p)

	if err != nil {
		return nil, err
	}

	res,err := util.PostReq("https://music.163.com/weapi/playlist/list",[]byte(str),"weapi")
	if err != nil {
		return nil, err
	}
	r := &PlayListRes{}
	err = json.Unmarshal(res,r)

	if err != nil {
		return nil, err
	}
	return r, nil
}

type PlayListDetailReq struct {
	Id int `json:"id"`
	N int `json:"n"`
	S int `json:"s"`
}

type PlayListDetailRes struct {
	PlayList interface{} `json:"playlist"`
	Code int `json:"code"`
	Privileges []interface{} `json:"privileges"`
}

func (p *PlayListDetailReq)GetPlayListDetail() (*PlayListDetailRes, error) {
	str,err := json.Marshal(p)

	if err != nil {
		return nil, err
	}

	res,err := util.PostReq("https://music.163.com/weapi/v3/playlist/detail",[]byte(str),"linuxapi")
	if err != nil {
		return nil, err
	}
	r := &PlayListDetailRes{}
	err = json.Unmarshal(res,r)

	if err != nil {
		return nil, err
	}
	return r, nil
}