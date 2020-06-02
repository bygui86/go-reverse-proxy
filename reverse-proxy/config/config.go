package config

import (
	"errors"

	"github.com/bygui86/go-reverse-proxy/reverse-proxy/logging"
	"github.com/bygui86/go-reverse-proxy/reverse-proxy/utils"
)

const (
	proxyHostEnvVar      = "PROXY_HOST"
	proxyPortEnvVar      = "PROXY_PORT"
	targetUrlEnvVar      = "TARGET_URL" // format: {protocol}://{host}:{port}
	nestedLevelNumEnvVar = "NESTED_LEVEL_NUMBER"

	proxyHostDefault      = "localhost"
	proxyPortDefault      = 8080
	targetUrlDefault      = ""
	nestedLevelNumDefault = 5
)

type config struct {
	ProxyHost      string
	ProxyPort      int
	TargetUrl      string
	NestedLevelNum int
}

func LoadConfig() *config {
	logging.Log.Debug("Load configurations")
	return &config{
		ProxyHost:      utils.GetStringEnv(proxyHostEnvVar, proxyHostDefault),
		ProxyPort:      utils.GetIntEnv(proxyPortEnvVar, proxyPortDefault),
		TargetUrl:      utils.GetStringEnv(targetUrlEnvVar, targetUrlDefault),
		NestedLevelNum: utils.GetIntEnv(nestedLevelNumEnvVar, nestedLevelNumDefault),
	}
}

func (c *config) ValidateConfig() error {
	if c.TargetUrl == targetUrlDefault {
		return errors.New("target url not defined")
	}
	return nil
}
