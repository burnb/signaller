package exchange

import (
	"github.com/burnb/signaller/internal/repository/entities"
)

type Client interface {
	Name() string
	RefreshTraders(traders []*entities.Trader)
	TopTraders() (traders []*entities.Trader, err error)
	TraderPositions(uid string) ([]*entities.Position, error)
}
