package provider

import (
	"github.com/burnb/signaller/internal/repository/entities"
	"github.com/burnb/signaller/pkg/grpc/api/proto"
)

type ExchangeClient interface {
	Name() string
	TopTraders() (traders []*entities.Trader, err error)
	RefreshTraders(traders []*entities.Trader)
	TraderPositions(uid string) ([]*entities.Position, error)
}

type publisher interface {
	Publish(event *proto.PositionEvent) error
}

type PositionResult struct {
	Error     error
	Positions []*entities.Position
}

type PositionJob struct {
	User     *entities.Trader
	ResultCh chan *PositionResult
}

type PositionsCh struct {
	User      *entities.Trader
	Positions []*entities.Position
}
