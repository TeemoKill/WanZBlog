package main

import (
	"github.com/TeemoKill/WanZBlog/framework"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	engine := framework.New()

	if err := engine.Init(); err != nil {
		return
	}

	engine.StartService()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
	if err := engine.Stop(); err != nil {
		return
	}

}
