package main

import (
	"WYMusicBackend/src/daemon"
	"WYMusicBackend/src/server"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Text struct {
	A int `json:"a"`
}

func main() {
	// 默认输出到日志
	file, err := os.OpenFile("WYMusic.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		os.Exit(-1)
	}
	defer func() {
		_ = file.Close()
	}()
	os.Stderr, os.Stdout = file, file

	daemon.AppPath = os.Args[0]
	switch isDaemon, err := daemon.Daemonize(); {
	case !isDaemon:
		fmt.Println("other process is running!")
		return
	case err != nil:
		panic(err)
	}

	srv := server.NewServer()
	srv.Start()
	fmt.Println("srv starting...")

	c := make(chan os.Signal)
	signal.Notify(c)
	signal.Ignore(syscall.SIGCHLD, syscall.SIGPIPE, syscall.SIGHUP)

	<-c

	srv.Stop()

	fmt.Println("srv stopped")

}
