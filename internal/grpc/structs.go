package grpc

import (
	"github.com/burnb/signaller/pkg/grpc/api/proto"
)

type Subscription struct {
	uids   map[string]struct{}
	stream proto.SignallerApi_StreamServer
}

func NewSubscription(uids map[string]struct{}, stream proto.SignallerApi_StreamServer) *Subscription {
	return &Subscription{uids: uids, stream: stream}
}
