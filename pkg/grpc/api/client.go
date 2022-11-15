package api

import (
	"context"
	"fmt"
	"sync"
	"time"

	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/burnb/signaller/pkg/grpc/api/proto"
)

type Client struct {
	address    string
	publishers sync.Map
	followCh   chan []string
	unfollowCh chan []string
	log        *zap.Logger
	client     proto.SignallerApiClient
	stream     proto.SignallerApi_StreamClient
}

func NewClient(address string, publishers []string, log *zap.Logger) *Client {
	c := &Client{
		address:    address,
		followCh:   make(chan []string),
		unfollowCh: make(chan []string),
		log:        log.Named(loggerName),
	}
	for _, publisher := range publishers {
		c.publishers.Store(publisher, struct{}{})
	}

	return c
}

func (c *Client) Init(ctx context.Context) error {
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithChainUnaryInterceptor(
			retry.UnaryClientInterceptor(
				retry.WithMax(10),
				retry.WithBackoff(retry.BackoffLinear(time.Second)),
				retry.WithCodes(codes.Aborted, codes.Unavailable),
			),
		),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(defaultClientMaxReceiveMessageSize),
			grpc.MaxCallSendMsgSize(defaultClientMaxSendMessageSize),
		),
	}

	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	conn, err := grpc.DialContext(ctx, c.address, options...)
	if err != nil {
		return err
	}

	c.client = proto.NewSignallerApiClient(conn)

	return nil
}

func (c *Client) AddPublishers(publishers []string) {
	c.followCh <- publishers
}

func (c *Client) RemovePublishers(publishers []string) {
	c.unfollowCh <- publishers
}

func (c *Client) Subscribe(ctx context.Context) (<-chan *proto.PositionEvent, error) {
	err := c.subscribeToTheStream(ctx)
	if err != nil {
		return nil, err
	}

	go c.runPublishersWorker(ctx)

	ch := make(chan *proto.PositionEvent)
	go c.runStreamWorker(ctx, ch)

	return ch, nil
}

func (c *Client) subscribeToTheStream(ctx context.Context) error {
	stream, err := c.client.Stream(ctx)
	if err != nil {
		return fmt.Errorf("unable to get position event stream %w", err)
	}

	var publishers []string
	c.publishers.Range(
		func(k, _ any) bool {
			publisher, _ := k.(string)
			publishers = append(publishers, publisher)

			return true
		},
	)
	err = stream.Send(&proto.SubscriptionRequest{Type: proto.SubscriptionRequestType_ADD, Uids: publishers})
	if err != nil {
		return fmt.Errorf("unable to send initial subscription request to the stream %w", err)
	}
	c.stream = stream

	return nil
}

func (c *Client) runPublishersWorker(ctx context.Context) {
	for {
		select {
		case publishers := <-c.followCh:
			if c.stream != nil {
				err := c.stream.Send(&proto.SubscriptionRequest{Type: proto.SubscriptionRequestType_ADD, Uids: publishers})
				if err != nil {
					c.log.Error("unable to send new additional publishers to the stream", zap.Error(err))
				}
			}
			for _, publisher := range publishers {
				c.publishers.Store(publisher, struct{}{})
			}
		case rmPublishers := <-c.unfollowCh:
			if c.stream != nil {
				err := c.stream.Send(&proto.SubscriptionRequest{Type: proto.SubscriptionRequestType_REMOVE, Uids: rmPublishers})
				if err != nil {
					c.log.Error("unable to send publishers for unfollow to the stream", zap.Error(err))
				}
			}
			for _, rmPublisher := range rmPublishers {
				c.publishers.LoadAndDelete(rmPublisher)
			}
		case <-ctx.Done():
			break
		}
	}
}

func (c *Client) runStreamWorker(ctx context.Context, ch chan<- *proto.PositionEvent) {
	for {
		pbPositionEvent, err := c.stream.Recv()
		if err != nil {
			select {
			case <-ctx.Done():
				close(ch)

				return
			default:
				c.log.Error("unable to receive stream", zap.Error(err))
				for {
					subErr := c.subscribeToTheStream(ctx)
					if subErr != nil {
						c.log.Error("unable to resubscribe to the position event stream", zap.Error(subErr))
						time.Sleep(2 * time.Second)
						continue
					}

					c.log.Info("position event stream connection established")
					break
				}

				continue
			}
		}

		ch <- pbPositionEvent
	}
}
