package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type Proxy struct {
	Path    string  `envconfig:"PROXY_LIST_PATH" default:"./proxy.txt"`
	Gateway *string `envconfig:"PROXY_GATEWAY"`
}

func (c *Proxy) Prepare() error {
	return envconfig.Process("", c)
}
