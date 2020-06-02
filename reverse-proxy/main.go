package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bygui86/go-reverse-proxy/reverse-proxy/config"
	"github.com/bygui86/go-reverse-proxy/reverse-proxy/logging"
	"github.com/bygui86/go-reverse-proxy/reverse-proxy/reverse_proxy"
)

var reverseProxy *reverse_proxy.ReverseProxy

func main() {
	logging.Log.Info("Start reverse proxy")

	logging.Log.Info("Load configurations")
	cfg := config.LoadConfig()

	logging.Log.Debug("Validate configurations")
	cfgErr := cfg.ValidateConfig()
	if cfgErr != nil {
		logging.SugaredLog.Errorf("Config validation failed: %s", cfgErr.Error())
		os.Exit(501)
	}

	reverseProxy = startReverseProxy(cfg.TargetUrl, cfg.ProxyHost, cfg.ProxyPort, cfg.NestedLevelNum)

	logging.Log.Info("Reverse proxy up and running")

	startSysCallChannel()

	shutdownAndWait(3)
}

func startReverseProxy(targetUrl string, host string, port int, nestedLevelNum int) *reverse_proxy.ReverseProxy {
	logging.Log.Debug("Start reverse proxy")
	proxy, err := reverse_proxy.NewReverseProxy(targetUrl, host, port, nestedLevelNum)
	if err != nil {
		logging.SugaredLog.Errorf("Reverse proxy creation failed: %s", err.Error())
		os.Exit(501)
	}
	logging.Log.Debug("Reverse proxy server successfully created")

	proxy.Start()
	logging.Log.Debug("Reverse proxy server successfully started")

	return proxy
}

func startSysCallChannel() {
	syscallCh := make(chan os.Signal)
	signal.Notify(syscallCh, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-syscallCh
}

func shutdownAndWait(timeout int) {
	logging.SugaredLog.Warnf("Termination signal received! Timeout %d", timeout)

	if reverseProxy != nil {
		reverseProxy.Shutdown(timeout)
	}

	time.Sleep(time.Duration(timeout) * time.Second)
}
