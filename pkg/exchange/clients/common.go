package clients

import (
	"github.com/burnb/signaller/internal/repository/entities"
)

type Client interface {
	Name() string
	Traders(uids []string) (traders []*entities.Trader, err error)
	RefreshTraders(traders []*entities.Trader)
	TopTraders() (traders []*entities.Trader, err error)
	TraderPositions(uid string) ([]*entities.Position, error)
}
