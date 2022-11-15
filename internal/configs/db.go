package configs

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Db struct {
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     uint16 `envconfig:"DB_PORT" required:"true"`
	Name     string `envconfig:"DB_DATABASE" required:"true"`
	Username string `envconfig:"DB_USERNAME" required:"true"`
	Password string `envconfig:"DB_PASSWORD"`
}

func (c *Db) Prepare() error {
	return envconfig.Process("", c)
}

func (c *Db) GetDatabaseDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", c.Username, c.Password, c.Host, c.Port, c.Name)
}
