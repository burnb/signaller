package configs

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Db struct {
	Host     string `envconfig:"DB_HOST"`
	Port     uint16 `envconfig:"DB_PORT"`
	Name     string `envconfig:"DB_DATABASE"`
	Username string `envconfig:"DB_USERNAME"`
	Password string `envconfig:"DB_PASSWORD"`
}

func (c *Db) Prepare() error {
	return envconfig.Process("", c)
}

func (c *Db) GetDatabaseDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", c.Username, c.Password, c.Host, c.Port, c.Name)
}
