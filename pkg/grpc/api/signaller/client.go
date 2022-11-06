package signaller

import (
	"context"
	"fmt"
	"io"
	"time"

	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/burnb/signaller/pkg/grpc/api/proto"
)

const (
	clientName = "SignallerClientGRPC"

	defaultClientMaxReceiveMessageSize = 1024 * 1024
	defaultClientMaxSendMessageSize    = 1024 * 1024
)

type Client struct {
	address string
	log     *zap.Logger
	client  proto.SignallerApiClient
}

func NewClient(address string, log *zap.Logger) *Client {
	return &Client{
		address: address,
		log:     log.Named(clientName),
	}
}

func (c *Client) Init() error {
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

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, c.address, options...)
	if err != nil {
		return err
	}

	c.client = proto.NewSignallerApiClient(conn)

	return nil
}

func (c *Client) Subscribe(uids []string) (<-chan *proto.PositionEvent, error) {
	stream, err := c.client.Stream(context.Background(), &proto.SubscribeRequest{Uids: uids})
	if err != nil {
		return nil, fmt.Errorf("unable to get order event stream %w", err)
	}

	ch := make(chan *proto.PositionEvent)
	go c.runWorker(stream, ch)

	return ch, nil
}

func (c *Client) runWorker(stream proto.SignallerApi_StreamClient, ch chan<- *proto.PositionEvent) {
	for {
		pbPositionEvent, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				c.log.Error("unable to receive stream", zap.Error(err))
			}
			close(ch)

			return
		}

		ch <- pbPositionEvent
	}
}
