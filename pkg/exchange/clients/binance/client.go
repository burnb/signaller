package binance

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

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
	if err != nil {
		return nil, err
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
			continue
		}

		break
	}
	if err != nil {
		return nil, err
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
		pos.Symbol = strings.ReplaceAll(pos.Symbol, "-SWAP", "")
		pos.TraderUID = uid
		pos.Long = pos.Amount > 0
		pos.Exchange = entities.ExchangeBinance
		pos.MarginMode = entities.MarginModeCross
		pos.UpdatedAt = time.UnixMilli(pos.UpdateTimestamp)
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
	tradersCh := make(chan *entities.Trader)
	wg := &sync.WaitGroup{}
	for w := 0; w < int(c.workersCnt); w++ {
		wg.Add(1)
		go func() {
			for trader := range tradersCh {
				if err := c.refreshTraderBaseInfo(trader); err != nil {
					c.log.Error("unable to get trader base info", zap.String("uid", trader.Uid), zap.Error(err))
					continue
				}
				if err := c.refreshTraderStats(trader); err != nil {
					c.log.Error("unable to get trader stats", zap.String("uid", trader.Uid), zap.Error(err))
					continue
				}
			}
			wg.Done()
		}()
	}

	for _, trader := range traders {
		tradersCh <- trader
	}

	close(tradersCh)

	wg.Wait()
}

func (c *Client) refreshTraderStats(trader *entities.Trader) error {
	var err error
	var body []byte
	for i := 0; i < DefaultRetry; i++ {
		body, err = c.do(
			true,
			http.MethodPost,
			`https://www.binance.com/bapi/futures/v1/public/future/leaderboard/getOtherPerformance`,
			fmt.Sprintf(`{"encryptedUid":"%s","tradeType":"PERPETUAL"}`, trader.Uid),
		)
		if err != nil {
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
			trader.RoiWeekly = period.Value
		}
		if period.PeriodType == "EXACT_WEEKLY" && period.StatisticsType == "PNL" {
			trader.PnlWeekly = period.Value
		}
		if period.PeriodType == "EXACT_MONTHLY" && period.StatisticsType == "ROI" {
			trader.RoiMonthly = period.Value
		}
		if period.PeriodType == "EXACT_MONTHLY" && period.StatisticsType == "PNL" {
			trader.PnlMonthly = period.Value
		}
		if period.PeriodType == "EXACT_YEARLY" && period.StatisticsType == "ROI" {
			trader.RoiYearly = period.Value
		}
		if period.PeriodType == "EXACT_YEARLY" && period.StatisticsType == "PNL" {
			trader.PnlYearly = period.Value
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
			fmt.Sprintf(`{"encryptedUid":"%s"}`, trader.Uid),
		)
		if err != nil {
			continue
		}

		break
	}
	if err != nil {
		return err
	}

	data := BaseInfoResponse{}
	err = data.UnmarshalJSON(body)
	if err != nil {
		return err
	}
	if !data.Success {
		return errors.New(data.Message)
	}

	trader.PositionShared = data.Data.PositionShared

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
