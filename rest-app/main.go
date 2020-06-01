package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bygui86/go-reverse-proxy/rest-app/blogpost"
	"github.com/bygui86/go-reverse-proxy/rest-app/logging"
)

var restServer *blogpost.Server

func main() {
	logging.Log.Info("Start rest-app")

	restServer = startBlogPostServer()

	logging.Log.Info("rest-app up and running")

	startSysCallChannel()

	shutdownAndWait(3)
}

func startBlogPostServer() *blogpost.Server {
	logging.Log.Debug("Start BlogPost server")

	server := blogpost.NewRestServer()
	logging.Log.Debug("BlogPost server successfully created")

	server.Start()
	logging.Log.Debug("BlogPost server successfully started")

	return server
}

func startSysCallChannel() {
	syscallCh := make(chan os.Signal)
	signal.Notify(syscallCh, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-syscallCh
}

func shutdownAndWait(timeout int) {
	logging.SugaredLog.Warnf("Termination signal received! Timeout %d", timeout)

	if restServer != nil {
		restServer.Shutdown(timeout)
	}

	time.Sleep(time.Duration(timeout+1) * time.Second)
}
