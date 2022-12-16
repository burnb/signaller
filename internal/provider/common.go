package provider

import (
	"github.com/burnb/signaller/internal/repository/entities"
	"github.com/burnb/signaller/pkg/grpc/api/proto"
)

type exchangeClient interface {
	Name() string
	RefreshTraders(traders []*entities.Trader)
	TraderPositions(uid string) ([]*entities.Position, error)
}

type publisher interface {
	FollowTraderUidCh() <-chan string
	UnFollowTraderUidCh() <-chan string
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
