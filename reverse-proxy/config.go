package main

import (
	"errors"
)

const (
	reverseProxyPortEnvVar = "REVERSE_PROXY_PORT"
	targetUrlEnvVar        = "TARGET_URL"

	reverseProxyPortDefault = 8080
	targetUrlDefault        = ""
)

type config struct {
	reverseProxyPort int
	targetUrl        string
}

func loadConfig() *config {
	Log.Debug("Load configurations")
	return &config{
		reverseProxyPort: GetIntEnv(reverseProxyPortEnvVar, reverseProxyPortDefault),
		targetUrl:        GetStringEnv(targetUrlEnvVar, targetUrlDefault),
	}
}

func (c *config) validateConfig() error {
	if c.targetUrl == targetUrlDefault {
		return errors.New("target url not defined")
	}
	return nil
}
