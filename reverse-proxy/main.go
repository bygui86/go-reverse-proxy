package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bygui86/go-reverse-proxy/reverse-proxy/config"
)

var reverseProxy *ReverseProxy

func main() {
	Log.Info("Start reverse proxy")

	Log.Info("Load configurations")
	cfg := loadConfig()

	Log.Debug("Validate configurations")
	cfgErr := cfg.validateConfig()
	if cfgErr != nil {
		SugaredLog.Errorf("Config validation failed: %s", cfgErr.Error())
		os.Exit(501)
	}

	reverseProxy = startReverseProxy(cfg.targetUrl, cfg.reverseProxyPort)

	Log.Info("Reverse proxy up and running")

	startSysCallChannel()

	shutdownAndWait(3)
}

func startReverseProxy(targetUrl string, port int) *ReverseProxy {
	Log.Debug("Start reverse proxy")
	proxy, err := newReverseProxy(targetUrl, port)
	if err != nil {
		SugaredLog.Errorf("Reverse proxy creation failed: %s", err.Error())
		os.Exit(501)
	}
	Log.Debug("Reverse proxy server successfully created")

	proxy.start()
	Log.Debug("Reverse proxy server successfully started")

	return proxy
}

func startSysCallChannel() {
	syscallCh := make(chan os.Signal)
	signal.Notify(syscallCh, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-syscallCh
}

func shutdownAndWait(timeout int) {
	SugaredLog.Warnf("Termination signal received! Timeout %d", timeout)

	if reverseProxy != nil {
		reverseProxy.shutdown(timeout)
	}

	time.Sleep(time.Duration(timeout) * time.Second)
}
