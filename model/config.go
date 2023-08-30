package model

import (
	"strings"
	"sync"
)

type ServerConfig struct {
	Claude    *Claude `mapstructure:"claude" json:"claude" yaml:"claude"`
	HttpProxy string `mapstructure:"http-proxy" json:"http-proxy" yaml:"http-proxy"`
	TlsCert   string `mapstructure:"tls-cert" json:"tls-cert" yaml:"tls-cert"`
	TlsKey    string `mapstructure:"tls-key" json:"tls-key" yaml:"tls-key"`
	ListenHost string `mapstructure:"listen-host" json:"listen-host" yaml:"listen-host"`
}

type Claude struct {
	SessionKeys []string   `mapstructure:"session-keys" json:"session-keys" yaml:"session-keys"`
	Lock        sync.Mutex `json:"-"`
}

func (c *Claude) GetSessionKey() (sessionKey string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	if len(c.SessionKeys) == 0 {
		return
	}
	sessionKey = c.SessionKeys[0]
	if !strings.HasPrefix(sessionKey, "sessionKey=") {
		sessionKey = "sessionKey=" + sessionKey
	}
	c.SessionKeys = append(c.SessionKeys[1:], sessionKey)
	return
}
