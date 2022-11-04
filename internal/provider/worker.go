package provider

import (
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burnb/signaller/internal/repository"
	"github.com/burnb/signaller/internal/repository/entities"
	"github.com/burnb/signaller/pkg/grpc/api/proto"
)

type Service struct {
	mu             sync.RWMutex
	log            *zap.Logger
	exchangeClient ExchangeClient
	repo           *repository.Mysql
	pub            publisher
	traders        []*entities.Trader
	positions      map[string]map[string]*entities.Position
}

func NewService(log *zap.Logger, exClient ExchangeClient, repo *repository.Mysql, pub publisher) *Service {
	return &Service{
		log:            log.Named(LoggerNameServiceWorker),
		exchangeClient: exClient,
		repo:           repo,
		pub:            pub,
		positions:      make(map[string]map[string]*entities.Position),
	}
}

func (s *Service) InitAndServe() error {
	if err := s.restore(); err != nil {
		return err
	}
	s.runPositionRefreshWorker()

	return nil
}

func (s *Service) restore() error {
	traders, err := s.repo.TradersWithSub()
	if err != nil {
		return err
	}
	s.traders = traders

	for _, trader := range s.traders {
		positions, err := s.repo.OpenedPositions(trader)
		if err != nil {
			return err
		}
		for _, position := range positions {
			if _, ok := s.positions[trader.EncryptedUid]; !ok {
				s.positions[trader.EncryptedUid] = make(map[string]*entities.Position)
			}
			s.positions[trader.EncryptedUid][position.Key()] = position
		}
	}

	return nil
}

func (s *Service) runPositionRefreshWorker() {
	go func() {
		for {
			for _, trader := range s.traders {
				positions, err := s.exchangeClient.TraderPositions(trader.EncryptedUid)
				if err != nil {
					s.log.Error("unable to get trader positions", zap.String("uid", trader.EncryptedUid), zap.Error(err))
					continue
				}
				trader.Positions = positions

				s.handleTrader(trader)
			}
			time.Sleep(10 * time.Second)
		}
	}()
}

func (s *Service) handleTrader(trader *entities.Trader) {
	var createEvents, updateEvents, closeEvents []*proto.PositionEvent
	symbolsCounter := make(map[string]int)
	for _, position := range trader.Positions {
		if _, ok := s.positions[position.UserId]; !ok {
			s.positions[position.UserId] = make(map[string]*entities.Position)
		}

		if _, ok := symbolsCounter[position.Symbol]; !ok {
			symbolsCounter[position.Symbol] = 0
		}
		symbolsCounter[position.Symbol]++

		cmd := proto.Command_CREATE
		oldPosition, ok := s.positions[position.UserId][position.Key()]
		if ok {
			if oldPosition.UpdateTimestamp == position.UpdateTimestamp {
				continue
			}
			cmd = proto.Command_UPDATE
			position.Id = oldPosition.Id
		}

		if err := s.repo.UpdatePosition(position); err != nil {
			s.log.Fatal("unable to update position", zap.Int64("id", position.Id), zap.Error(err))
		}

		direction := proto.Direction_LONG
		if position.Amount < 0 {
			direction = proto.Direction_SHORT
		}

		event := &proto.PositionEvent{
			Symbol:     position.Symbol,
			Uid:        position.UserId,
			Direction:  direction,
			PositionId: position.Id,
			Command:    cmd,
			Exchange:   s.exchangeClient.Name(),
			Leverage:   uint32(position.Leverage),
			EntryPrice: position.EntryPrice,
			CreatedAt:  timestamppb.New(time.Unix(position.UpdateTimestamp, 0)),
		}
		if cmd == proto.Command_CREATE {
			createEvents = append(createEvents, event)
		} else {
			updateEvents = append(updateEvents, event)
		}

		s.positions[position.UserId][position.Key()] = position
	}

	// check closed order
	for key, oldPosition := range s.positions[trader.EncryptedUid] {
		exists := false
		for _, position := range trader.Positions {
			position.Symbol = strings.ReplaceAll(position.Symbol, "-SWAP", "")
			if position.Key() == oldPosition.Key() {
				exists = true
				break
			}
		}
		if !exists {
			if err := s.repo.ClosePosition(oldPosition); err != nil {
				s.log.Fatal("unable to close position", zap.Int64("id", oldPosition.Id), zap.Error(err))
			}

			delete(s.positions[trader.EncryptedUid], key)

			dir := proto.Direction_LONG
			if oldPosition.Amount < 0 {
				dir = proto.Direction_SHORT
			}
			event := &proto.PositionEvent{
				Symbol:     oldPosition.Symbol,
				Uid:        trader.EncryptedUid,
				Direction:  dir,
				PositionId: oldPosition.Id,
				Command:    proto.Command_CLOSE,
				Exchange:   s.exchangeClient.Name(),
				Leverage:   uint32(oldPosition.Leverage),
				EntryPrice: oldPosition.EntryPrice,
				CreatedAt:  timestamppb.New(time.Now()),
			}
			closeEvents = append(closeEvents, event)
		}
	}

	// Check Hedge mode
	var eventsList []*proto.PositionEvent
	eventsList = append(eventsList, createEvents...)
	eventsList = append(eventsList, updateEvents...)
	eventsList = append(eventsList, closeEvents...)
	for _, event := range eventsList {
		if count, ok := symbolsCounter[event.Symbol]; ok && count > 1 {
			event.Hedged = true
		}
	}

	for _, event := range closeEvents {
		if pushErr := s.pub.Publish(event); pushErr != nil {
			s.log.Error("unable to push position closed event", zap.Int64("id", event.PositionId), zap.Error(pushErr))
		}
	}
	for _, event := range updateEvents {
		if pushErr := s.pub.Publish(event); pushErr != nil {
			s.log.Error("unable to push position updated event", zap.Int64("id", event.PositionId), zap.Error(pushErr))
		}
	}
	for _, event := range createEvents {
		if pushErr := s.pub.Publish(event); pushErr != nil {
			s.log.Error("unable to push position created event", zap.Int64("id", event.PositionId), zap.Error(pushErr))
		}
	}
}
