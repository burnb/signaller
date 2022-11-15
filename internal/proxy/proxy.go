package proxy

import (
	"bufio"
	"context"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	netProxy "golang.org/x/net/proxy"

	"github.com/burnb/signaller/internal/configs"
)

type Service struct {
	mu      sync.RWMutex
	cfg     configs.Proxy
	logger  *zap.Logger
	current uint16
	proxies []string
}

func New(cfg configs.Proxy, log *zap.Logger) *Service {
	return &Service{cfg: cfg, logger: log.Named(loggerName)}
}

func (s *Service) Init() error {
	inFile, err := os.Open(s.cfg.Path)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return nil
		}

		return err
	}
	defer func() {
		if err := inFile.Close(); err != nil {
			s.logger.Error("unable to close proxy list file", zap.Error(err))
		}
	}()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		s.proxies = append(s.proxies, strings.TrimSpace(strings.ReplaceAll(scanner.Text(), "socks5://", "")))
	}

	s.logger.Info("loaded proxies", zap.Int("cnt", len(s.proxies)))

	return nil
}

func (s *Service) IsEnabled() bool {
	s.mu.RLock()
	cnt := len(s.proxies)
	s.mu.RUnlock()

	return cnt > 0
}

func (s *Service) Gateway() *string {
	return s.cfg.Gateway
}

func (s *Service) Address(rnd bool) (address string) {
	if rnd {
		i := uint16(rand.Intn(len(s.proxies) - 1))

		s.mu.RLock()
		address = s.proxies[i]
		s.mu.RUnlock()

		return address
	}

	s.mu.Lock()
	address = s.proxies[s.current]
	s.current++
	if s.current == uint16(len(s.proxies)) {
		s.current = 0
	}
	s.mu.Unlock()

	return address
}

func (s *Service) Dialer(rnd bool) (netProxy.Dialer, error) {
	if !s.IsEnabled() {
		return netProxy.FromEnvironment(), nil
	}

	forwardDialer := netProxy.Dialer(netProxy.Direct)
	if s.Gateway() != nil {
		dialerGate, err := netProxy.SOCKS5("tcp", *s.Gateway(), nil, netProxy.Direct)
		if err != nil {
			return nil, err
		}
		forwardDialer = dialerGate
	}

	return netProxy.SOCKS5("tcp", s.Address(rnd), nil, forwardDialer)
}

func (s *Service) HttpClient(rndProxy bool) (*http.Client, error) {
	dialer, err := s.Dialer(rndProxy)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
				return dialer.Dial(network, address)
			},
		},
		Timeout: time.Second * 5,
	}, nil
}
