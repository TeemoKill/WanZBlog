package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/TeemoKill/WanZBlog/framework"
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
