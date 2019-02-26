package server

import (
	"WYMusicBackend/src/server/handler/index"
	"WYMusicBackend/src/server/handler/playlist"
	"WYMusicBackend/src/server/handler/top"
	"net/http"
	"os"
	"time"
)

type Server struct {
	localIP string
	server *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server)Start()  {
	s.localIP = "127.0.0.1"

	mux := http.NewServeMux()

	mux.Handle("/top/playlist",&top.PlayList{})
	mux.Handle("/index/apiGetInfo",&index.ApiGetInfo{})
	mux.Handle("/playlist/apiGetPlayListDetail", &playlist.ApiGetPlayListDetail{})


	s.server = &http.Server{
		Addr:  ":8000",
		ReadTimeout: 60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:mux,
	}

	go func() {
		err := s.server.ListenAndServe(); if err.Error() != http.ErrServerClosed.Error() {
			os.Exit(-1)
		}
	}()
}

func (s *Server)Stop()  {
	_ = s.server.Close()
}

func InitLog()  {
	
}