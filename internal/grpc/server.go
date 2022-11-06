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
	"google.golang.org/grpc/codes"
	grpcKeepalive "google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/burnb/signaller/pkg/grpc/api/proto"
)

const (
	serverName = "ServerGRPC"

	defaultServerMaxReceiveMessageSize = 1024 * 1024
	defaultServerMaxSendMessageSize    = 1024 * 1024
)

type Server struct {
	address             string
	log                 *zap.Logger
	grpcServer          *grpc.Server
	listener            net.Listener
	subscriptions       sync.Map
	followTraderUidCh   chan string
	unFollowTraderUidCh chan string

	proto.UnimplementedSignallerApiServer
}

func NewServer(address string, log *zap.Logger) *Server {
	return &Server{
		address:             address,
		log:                 log.Named(serverName),
		followTraderUidCh:   make(chan string, 10),
		unFollowTraderUidCh: make(chan string, 10),
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

func (s *Server) Stream(req *proto.SubscribeRequest, stream proto.SignallerApi_StreamServer) error {
	if req.Uids == nil {
		return status.Errorf(codes.InvalidArgument, "empty uid list")
	}

	uids := make(map[string]struct{})
	for _, uid := range req.Uids {
		uids[uid] = struct{}{}
		s.followTraderUidCh <- uid
	}

	cUid := uuid.New().String()
	s.subscriptions.Store(cUid, NewSubscription(uids, stream))
	s.log.Info("client connected", zap.String("sub", cUid))

	<-stream.Context().Done()

	s.subscriptions.Delete(cUid)

	s.log.Info("client disconnected", zap.String("sub", cUid))

	return nil
}

func (s *Server) FollowTraderUidCh() <-chan string {
	return s.followTraderUidCh
}

func (s *Server) UnFollowTraderUidCh() <-chan string {
	return s.unFollowTraderUidCh
}

func (s *Server) Publish(event *proto.PositionEvent) error {
	var unsubscribe []string
	s.subscriptions.Range(
		func(k, v any) bool {
			uid, ok := k.(string)
			if !ok {
				s.log.Error("unable to cast subscription uid type", zap.String("type", fmt.Sprintf("%T", k)))

				return false
			}

			subscription, ok := v.(Subscription)
			if !ok {
				s.log.Error("unable to cast subscription type", zap.String("type", fmt.Sprintf("%T", v)))

				return false
			}

			if _, ok := subscription.uids[event.TraderUid]; ok {
				if err := subscription.stream.Send(event); err != nil {
					s.log.Error("unable to send data to client", zap.Error(err))

					unsubscribe = append(unsubscribe, uid)
				}
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
