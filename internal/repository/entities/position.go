package entities

//go:generate easyjson position.go

const (
	MarginModeCross    = "cross"
	MarginModeIsolated = "isolated"
	MarginModeHedge    = "hedge"
)

// easyjson:json
type Position struct {
	Id              int64
	UserId          string  `json:"uid" db:"uid"`
	Symbol          string  `json:"symbol"`
	EntryPrice      float64 `json:"entryPrice" db:"entryPrice"`
	MarkPrice       float64 `json:"markPrice" db:"markPrice"`
	Pnl             float64 `json:"pnl"`
	Roe             float64 `json:"roe"`
	Amount          float64 `json:"amount"`
	Leverage        uint8   `json:"leverage"`
	Invested        float64
	Opened          bool
	Long            bool
	CreateTimestamp int64  `json:"createTimestamp" db:"createdAt"`
	UpdateTimestamp int64  `json:"updateTimestamp" db:"updatedAt"`
	ClosedTimestamp *int64 `json:"closedTimestamp" db:"closedAt"`
	Exchange        string
	MarginMode      string `db:"margin_mode"`
	Hedged          bool
}

func (p *Position) Direction() string {
	if p.Long {
		return "LONG"
	}
	return "SHORT"
}

func (p *Position) Key() string {
	return p.Symbol + p.Direction() + p.MarginMode
}
