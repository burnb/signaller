package grpc

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	grpcKeepalive "google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/burnb/signaller/pkg/grpc/api/proto"
)

const (
	serverName = "ServerGRPC"

	defaultServerMaxReceiveMessageSize = 1024 * 1024
	defaultServerMaxSendMessageSize    = 1024 * 1024
)

type Server struct {
	address       string
	log           *zap.Logger
	grpcServer    *grpc.Server
	listener      net.Listener
	subscriptions sync.Map

	proto.UnimplementedSignallerApiServer
}

func NewServer(address string, log *zap.Logger) *Server {
	return &Server{
		address: address,
		log:     log.Named(serverName),
	}
}

func (s *Server) Init() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("unable to listen port %w", err)
	}
	s.listener = listener

	go func() {
		options := []grpc.ServerOption{
			grpc.MaxRecvMsgSize(defaultServerMaxReceiveMessageSize),
			grpc.MaxSendMsgSize(defaultServerMaxSendMessageSize),
			grpc.KeepaliveEnforcementPolicy(grpcKeepalive.EnforcementPolicy{
				MinTime:             5 * time.Second,
				PermitWithoutStream: true,
			}),
			grpc.KeepaliveParams(grpcKeepalive.ServerParameters{
				Time:    5 * time.Second,
				Timeout: 1 * time.Second,
			}),
			grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())),
		}
		s.grpcServer = grpc.NewServer(options...)
		reflection.Register(s.grpcServer)
		proto.RegisterSignallerApiServer(s.grpcServer, s)

		s.log.Info("grpc serve success")

		if err := s.grpcServer.Serve(s.listener); err != nil {
			s.log.Fatal("unable to serve", zap.Error(err))
		}
	}()

	return nil
}

func (s *Server) Stream(_ *proto.SubscribeRequest, stream proto.SignallerApi_StreamServer) error {
	uid := uuid.New().String()
	s.subscriptions.Store(uid, stream)
	s.log.Info("client connected", zap.String("sub", uid))

	ctx := stream.Context()
	<-ctx.Done()

	s.subscriptions.Delete(uid)

	s.log.Info("client disconnected", zap.String("sub", uid))

	return nil
}

func (s *Server) Publish(event *proto.PositionEvent) error {
	var unsubscribe []string
	s.subscriptions.Range(
		func(k, v interface{}) bool {
			uid, ok := k.(string)
			if !ok {
				s.log.Error("unable to cast subscription uid type", zap.String("type", fmt.Sprintf("%T", k)))

				return false
			}
			stream, ok := v.(proto.SignallerApi_StreamServer)
			if !ok {
				s.log.Error("unable to cast subscription stream type", zap.String("type", fmt.Sprintf("%T", v)))

				return false
			}

			if err := stream.Send(event); err != nil {
				s.log.Error("unable to send data to client", zap.Error(err))

				unsubscribe = append(unsubscribe, uid)
			}

			return true
		},
	)

	for _, uid := range unsubscribe {
		s.subscriptions.Delete(uid)
	}

	return nil
}

func (s *Server) Shutdown() {
	s.grpcServer.Stop()
}
