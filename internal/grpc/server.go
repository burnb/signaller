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
	followCh      chan string
	unFollowCh    chan string

	proto.UnimplementedSignallerApiServer
}

func NewServer(address string, log *zap.Logger) *Server {
	return &Server{
		address: address,
		log:     log.Named(serverName),
		grpcServer: grpc.NewServer(
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
		),
		followCh:   make(chan string, 10),
		unFollowCh: make(chan string, 10),
	}
}

func (s *Server) Init() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("unable to listen port %w", err)
	}
	s.listener = listener

	reflection.Register(s.grpcServer)
	proto.RegisterSignallerApiServer(s.grpcServer, s)

	go s.serve()

	return nil
}

func (s *Server) serve() {
	s.log.Info("grpc serve success")

	if err := s.grpcServer.Serve(s.listener); err != nil {
		s.log.Fatal("unable to serve", zap.Error(err))
	}
}

func (s *Server) Stream(stream proto.SignallerApi_StreamServer) error {
	cUid := uuid.New().String()
	subscription := NewSubscription(stream)
	s.subscriptions.Store(cUid, subscription)

	go s.runStreamWorker(subscription)

	s.log.Info("client connected", zap.String("sub", cUid))

	<-stream.Context().Done()

	s.subscriptions.Delete(cUid)

	s.log.Info("client disconnected", zap.String("sub", cUid))

	return nil
}

func (s *Server) runStreamWorker(subscription *Subscription) {
	for {
		req, err := subscription.stream.Recv()
		if err != nil {
			select {
			case <-subscription.stream.Context().Done():
				return
			default:
				s.log.Error("unable to receive stream", zap.Error(err))

				return
			}
		}

		if req.GetType() == proto.SubscriptionRequestType_ADD {
			for _, uid := range req.Uids {
				subscription.uids[uid] = struct{}{}
				s.followCh <- uid
			}
		} else {
			for _, uid := range req.Uids {
				delete(subscription.uids, uid)
				exist := false
				s.subscriptions.Range(func(k, v any) bool {
					subscription, _ := v.(*Subscription)
					if _, ok := subscription.uids[uid]; ok {
						exist = true

						return false
					}

					return true
				})
				if !exist {
					s.unFollowCh <- uid
				}
			}
		}
	}
}

func (s *Server) FollowTraderUidCh() <-chan string {
	return s.followCh
}

func (s *Server) UnFollowTraderUidCh() <-chan string {
	return s.unFollowCh
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

			subscription, ok := v.(*Subscription)
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
