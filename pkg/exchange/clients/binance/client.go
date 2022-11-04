package binance

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/adshao/go-binance/v2/futures"
	"go.uber.org/zap"

	"github.com/burnb/signaller/internal/proxy"
	"github.com/burnb/signaller/internal/repository/entities"
)

type Client struct {
	name               string
	log                *zap.Logger
	workersCnt         uint8
	proxySrv           *proxy.Service
	client             *futures.Client
	symbols            map[string]*SymbolInfo
	crossWalletBalance float64
}

func NewClient(log *zap.Logger, proxySrv *proxy.Service) *Client {
	return &Client{log: log.Named(Name), name: Name, workersCnt: DefaultWorkersCnt, proxySrv: proxySrv}
}

func (c *Client) Name() string {
	return c.name
}

func (c *Client) TopTraders() (traders []*entities.Trader, err error) {
	var body []byte
	for i := 0; i < DefaultRetry; i++ {
		body, err =
			c.do(
				true,
				http.MethodPost,
				`https://www.binance.com/bapi/futures/v1/public/future/leaderboard/searchLeaderboard`,
				`{"limit":150,"sortType":"ROI","isShared":true,"periodType":"EXACT_WEEKLY","pnlGainType":null,"roiGainType":null,"symbol":"","tradeType":"PERPETUAL"}`,
			)
		if err != nil {
			c.log.Error("unable to do top traders request", zap.Error(err))
			continue
		}

		break
	}

	data := &TraderResponse{}
	err = data.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}
	if !data.Success {
		return nil, errors.New(data.Message)
	}

	return data.Data, nil
}

func (c *Client) Traders(uids []string) (traders []*entities.Trader, err error) {
	data := &TraderResponse{}
	if !data.Success {
		return nil, errors.New(data.Message)
	}

	return data.Data, nil
}

func (c *Client) TraderPositions(uid string) ([]*entities.Position, error) {
	var err error
	var body []byte
	for i := 0; i < DefaultRetry; i++ {
		body, err = c.do(
			false,
			http.MethodPost,
			`https://www.binance.com/bapi/futures/v1/public/future/leaderboard/getOtherPosition`,
			fmt.Sprintf(`{"encryptedUid":"%s","tradeType":"PERPETUAL"}`, uid),
		)
		if err != nil {
			c.log.Error("unable to do trader position request", zap.Error(err))
			continue
		}

		break
	}

	data := &PositionsResponse{}
	err = data.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}
	if !data.Success {
		return nil, errors.New(data.Message)
	}

	for _, pos := range data.Data.OtherPositionRetList {
		pos.UserId = uid
		pos.Opened = true
		pos.Long = pos.Amount > 0
		pos.UpdateTimestamp = pos.UpdateTimestamp / 1000
		pos.Exchange = c.name
		pos.MarginMode = entities.MarginModeCross
		for _, pos2 := range data.Data.OtherPositionRetList {
			if pos.Symbol == pos2.Symbol && pos.Amount == -pos2.Amount {
				pos.MarginMode = entities.MarginModeHedge
				break
			}
		}
	}

	return data.Data.OtherPositionRetList, nil
}

func (c *Client) RefreshTraders(traders []*entities.Trader) {
	tradersCh := make(chan *entities.Trader, c.workersCnt)
	defer close(tradersCh)
	wg := &sync.WaitGroup{}
	for w := 0; w < int(c.workersCnt); w++ {
		wg.Add(1)
		go func() {
			for trader := range tradersCh {
				if err := c.refreshTraderBaseInfo(trader); err != nil {
					c.log.Error("unable to get trader base info", zap.String("uid", trader.EncryptedUid), zap.Error(err))
					continue
				}
				if err := c.refreshTraderStats(trader); err != nil {
					c.log.Error("unable to get trader stats", zap.String("uid", trader.EncryptedUid), zap.Error(err))
					continue
				}
			}
			wg.Done()
		}()
	}

	for _, trader := range traders {
		tradersCh <- trader
	}

	wg.Wait()
}

func (c *Client) refreshTraderStats(user *entities.Trader) error {
	var err error
	var body []byte
	for i := 0; i < DefaultRetry; i++ {
		body, err = c.do(
			true,
			http.MethodPost,
			`https://www.binance.com/bapi/futures/v1/public/future/leaderboard/getOtherPerformance`,
			fmt.Sprintf(`{"encryptedUid":"%s","tradeType":"PERPETUAL"}`, user.EncryptedUid),
		)
		if err != nil {
			c.log.Error("unable to do trader stats request", zap.Error(err))
			continue
		}

		break
	}
	if err != nil {
		return err
	}

	data := StatsResponse{}
	err = data.UnmarshalJSON(body)
	if err != nil {
		return err
	}
	if !data.Success {
		return errors.New(data.Message)
	}

	for _, period := range data.Data {
		if period.PeriodType == "EXACT_WEEKLY" && period.StatisticsType == "ROI" {
			user.WeeklyRoi = period.Value
		}
		if period.PeriodType == "EXACT_WEEKLY" && period.StatisticsType == "PNL" {
			user.WeeklyPnl = period.Value
		}
		if period.PeriodType == "EXACT_MONTHLY" && period.StatisticsType == "ROI" {
			user.MonthlyRoi = period.Value
		}
		if period.PeriodType == "EXACT_MONTHLY" && period.StatisticsType == "PNL" {
			user.MonthlyPnl = period.Value
		}
		if period.PeriodType == "EXACT_YEARLY" && period.StatisticsType == "ROI" {
			user.YearlyRoi = period.Value
		}
		if period.PeriodType == "EXACT_YEARLY" && period.StatisticsType == "PNL" {
			user.YearlyPnl = period.Value
		}
	}
	return nil
}

func (c *Client) refreshTraderBaseInfo(trader *entities.Trader) error {
	var err error
	var body []byte
	for i := 0; i < DefaultRetry; i++ {
		body, err = c.do(
			true,
			http.MethodPost,
			`https://www.binance.com/bapi/futures/v2/public/future/leaderboard/getOtherLeaderboardBaseInfo`,
			fmt.Sprintf(`{"encryptedUid":"%s"}`, trader.EncryptedUid),
		)
		if err != nil {
			c.log.Error("unable to do trader base info request", zap.Error(err))
			continue
		}

		break
	}
	if err != nil {
		return err
	}

	c.log.Debug("got user base info", zap.String("uid", trader.EncryptedUid))

	data := BaseInfoResponse{}
	err = data.UnmarshalJSON(body)
	if err != nil {
		return err
	}
	if !data.Success {
		return errors.New(data.Message)
	}

	trader.NickName = data.Data.NickName
	trader.PositionShared = data.Data.PositionShared
	trader.UserPhotoUrl = data.Data.UserPhotoUrl

	return nil
}

func (c *Client) do(rndProxy bool, method string, url string, body string) ([]byte, error) {
	client, err := c.proxySrv.HttpClient(rndProxy)
	if err != nil {
		return nil, fmt.Errorf("unable to get http client %w", err)
	}

	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("unable to get new http request %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", ApiClientUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to do request %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read request body %w", err)
	}

	return respBody, nil
}
