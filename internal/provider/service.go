package provider

import (
	"database/sql"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burnb/signaller/internal/configs"
	"github.com/burnb/signaller/internal/repository"
	"github.com/burnb/signaller/internal/repository/entities"
	"github.com/burnb/signaller/pkg/grpc/api/proto"
)

type Service struct {
	cfg            configs.Provider
	logger         *zap.Logger
	exchangeClient exchangeClient
	repo           *repository.Mysql
	publisher      publisher
	traders        sync.Map
	lastEventAt    time.Time
	lastSyncAt     time.Time
}

func NewService(cfg configs.Provider, log *zap.Logger, exClient exchangeClient, repo *repository.Mysql, pub publisher) *Service {
	return &Service{
		cfg:            cfg,
		logger:         log.Named(loggerName),
		exchangeClient: exClient,
		repo:           repo,
		publisher:      pub,
	}
}

func (s *Service) Init() error {
	if err := s.restore(); err != nil {
		return err
	}

	s.runPublisherUnFollowWorker()
	s.runPublisherFollowWorker()
	s.runPositionRefreshWorker()
	s.runTradersRefreshWorker()

	return nil
}

func (s *Service) TradersCount() uint {
	var length uint
	s.traders.Range(func(_, _ interface{}) bool {
		length++

		return true
	})

	return length
}

func (s *Service) LastEventAt() time.Time {
	return s.lastEventAt
}

func (s *Service) LastSyncAt() time.Time {
	return s.lastSyncAt
}

func (s *Service) restore() error {
	traders, err := s.repo.Publishers()
	if err != nil {
		return err
	}

	for _, trader := range traders {
		if err := s.restoreTradersOpenedPositions(trader); err != nil {
			return err
		}

		s.traders.Store(trader.Uid, trader)
	}

	return nil
}

func (s *Service) restoreTradersOpenedPositions(trader *entities.Trader) error {
	positions, err := s.repo.OpenedPositions(trader)
	if err != nil {
		return err
	}

	if trader.Positions == nil {
		trader.Positions = make(map[string]*entities.Position)
	}

	for _, position := range positions {
		trader.Positions[position.Key()] = position
	}

	return nil
}

func (s *Service) runPublisherFollowWorker() {
	go func() {
		for traderUid := range s.publisher.FollowTraderUidCh() {
			if _, ok := s.traders.Load(traderUid); ok {
				continue
			}

			trader, err := s.repo.Trader(traderUid)
			if err != nil {
				s.logger.Error("unable to get trader", zap.Error(err))
				continue
			}
			if trader == nil {
				trader = &entities.Trader{
					Uid:       traderUid,
					Publisher: true,
					Positions: make(map[string]*entities.Position),
				}

				s.exchangeClient.RefreshTraders([]*entities.Trader{trader})

				if !trader.PositionShared {
					s.logger.Warn("trader doesn't show his positions", zap.String("uid", trader.Uid))
					continue
				}

				if err := s.repo.CreateTrader(trader); err != nil {
					s.logger.Error("unable to create trader", zap.Error(err))
					continue
				}
			} else {
				trader.Publisher = true
				if err := s.repo.UpdateTrader(trader); err != nil {
					s.logger.Error("unable to update trader", zap.Error(err))
					continue
				}
				if err := s.restoreTradersOpenedPositions(trader); err != nil {
					s.logger.Error("unable to restore traders opened positions", zap.Error(err))
				}
			}
			s.traders.Store(traderUid, trader)
		}
	}()
}

func (s *Service) runPublisherUnFollowWorker() {
	go func() {
		for traderUid := range s.publisher.UnFollowTraderUidCh() {
			v, ok := s.traders.Load(traderUid)
			if !ok {
				continue
			}

			trader, _ := v.(*entities.Trader)
			trader.Publisher = false
			if err := s.repo.UpdateTrader(trader); err != nil {
				s.logger.Error(
					"unable to update trader after unfollow",
					zap.String("uid", traderUid),
					zap.Error(err),
				)
			}

			s.traders.Delete(traderUid)
		}
	}()
}

func (s *Service) runPositionRefreshWorker() {
	go func() {
		for {
			s.traders.Range(
				func(k, v any) bool {
					trader, _ := v.(*entities.Trader)
					positions, err := s.exchangeClient.TraderPositions(trader.Uid)
					if err != nil {
						s.logger.Error("unable to get trader positions", zap.String("uid", trader.Uid), zap.Error(err))

						return true
					}

					s.handleNewTraderPositions(trader, positions)

					return true
				},
			)

			if s.TradersCount() > 0 {
				s.lastSyncAt = time.Now()
			}

			time.Sleep(s.cfg.PositionRefreshTimeDuration())
		}
	}()
}

func (s *Service) runTradersRefreshWorker() {
	go func() {
		refreshDuration := s.cfg.TradersRefreshTimeDuration()
		for {
			var traders []*entities.Trader
			s.traders.Range(
				func(k, v any) bool {
					trader, _ := v.(*entities.Trader)
					traders = append(traders, trader)

					return true
				},
			)

			s.exchangeClient.RefreshTraders(traders)
			for _, trader := range traders {
				if err := s.repo.UpdateTrader(trader); err != nil {
					s.logger.Error(
						"unable to update trader after refresh",
						zap.String("uid", trader.Uid),
						zap.Error(err),
					)
				}
				if !trader.PositionShared {
					s.logger.Warn("trader hide his positions", zap.String("uid", trader.Uid))
				}
			}

			time.Sleep(refreshDuration)
		}
	}()
}

func (s *Service) handleNewTraderPositions(trader *entities.Trader, newPositions []*entities.Position) {
	var createEvents, updateEvents, closeEvents []*proto.PositionEvent
	symbols := make(map[string]struct{})
	for _, newPosition := range newPositions {
		if trader.Positions == nil {
			trader.Positions = make(map[string]*entities.Position)
		}

		var amountChange float64
		eventType := proto.Type_CREATE
		oldPosition, ok := trader.Positions[newPosition.Key()]
		if ok {
			newPosition.Id = oldPosition.Id
			if err := s.repo.UpdatePosition(newPosition); err != nil {
				s.logger.Panic("unable to update position", zap.Int64("id", newPosition.Id), zap.Error(err))
			}

			if oldPosition.UpdatedAt == newPosition.UpdatedAt {
				continue
			}

			eventType = proto.Type_UPDATE
			amountChange =
				float64(newPosition.Leverage)/float64(oldPosition.Leverage)*(newPosition.Amount/oldPosition.Amount) - 1
		} else if err := s.repo.CreatePosition(newPosition); err != nil {
			s.logger.Panic("unable to create position", zap.Int64("id", newPosition.Id), zap.Error(err))
		}

		direction := proto.Direction_LONG
		if newPosition.Amount < 0 {
			direction = proto.Direction_SHORT
		}

		var hedged bool
		if _, ok := symbols[newPosition.Symbol]; !ok {
			symbols[newPosition.Symbol] = struct{}{}
		} else {
			hedged = true
		}

		event := &proto.PositionEvent{
			Symbol:       newPosition.Symbol,
			TraderUid:    newPosition.TraderUID,
			Direction:    direction,
			PositionId:   newPosition.Id,
			Type:         eventType,
			Exchange:     s.exchangeClient.Name(),
			Leverage:     uint32(newPosition.Leverage),
			AmountChange: amountChange,
			EntryPrice:   newPosition.EntryPrice,
			CreatedAt:    timestamppb.New(time.UnixMilli(newPosition.UpdatedAt)),
			Hedged:       hedged,
		}
		if eventType == proto.Type_CREATE {
			createEvents = append(createEvents, event)
		} else {
			updateEvents = append(updateEvents, event)
		}

		trader.Positions[newPosition.Key()] = newPosition
	}

	// check closed order
	for key, oldPosition := range trader.Positions {
		exists := false
		for _, newPosition := range newPositions {
			if newPosition.Key() == oldPosition.Key() {
				exists = true
				break
			}
		}
		if !exists {
			oldPosition.ClosedAt = sql.NullTime{Time: time.Now(), Valid: true}
			if err := s.repo.UpdatePosition(oldPosition); err != nil {
				s.logger.Panic("unable to close position", zap.Int64("id", oldPosition.Id), zap.Error(err))
			}

			delete(trader.Positions, key)

			dir := proto.Direction_LONG
			if oldPosition.Amount < 0 {
				dir = proto.Direction_SHORT
			}
			event := &proto.PositionEvent{
				Symbol:     oldPosition.Symbol,
				TraderUid:  trader.Uid,
				Direction:  dir,
				PositionId: oldPosition.Id,
				Type:       proto.Type_CLOSE,
				Exchange:   s.exchangeClient.Name(),
				Leverage:   uint32(oldPosition.Leverage),
				EntryPrice: oldPosition.EntryPrice,
				CreatedAt:  timestamppb.New(time.Now()),
			}
			closeEvents = append(closeEvents, event)
		}
	}

	s.push(closeEvents)
	s.push(updateEvents)
	s.push(createEvents)
}

func (s *Service) push(events []*proto.PositionEvent) {
	for _, event := range events {
		if err := s.publisher.Publish(event); err != nil {
			s.logger.Error(
				"unable to push position event",
				zap.Int64("id", event.PositionId),
				zap.String("type", proto.Type_name[int32(event.Type)]),
				zap.Error(err),
			)
		}
		if err := s.repo.RefreshPublishTime(event.TraderUid); err != nil {
			s.logger.Error(
				"unable to refresh trader publish time",
				zap.String("trader_uid", event.TraderUid),
				zap.Error(err),
			)
		}

		s.lastEventAt = time.Now()
	}
}
