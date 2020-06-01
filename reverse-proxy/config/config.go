package config

import (
	"errors"

	"github.com/bygui86/go-reverse-proxy/reverse-proxy/logging"
	"github.com/bygui86/go-reverse-proxy/reverse-proxy/utils"
)

const (
	proxyPortEnvVar = "PROXY_PORT"
	targetUrlEnvVar = "TARGET_URL" // format: {protocol}://{host}:{port}

	proxyPortDefault = 8080
	targetUrlDefault = ""
)

type config struct {
	ProxyPort int
	TargetUrl string
}

func LoadConfig() *config {
	logging.Log.Debug("Load configurations")
	return &config{
		ProxyPort: utils.GetIntEnv(proxyPortEnvVar, proxyPortDefault),
		TargetUrl: utils.GetStringEnv(targetUrlEnvVar, targetUrlDefault),
	}
}

func (c *config) ValidateConfig() error {
	if c.TargetUrl == targetUrlDefault {
		return errors.New("target url not defined")
	}
	return nil
}
