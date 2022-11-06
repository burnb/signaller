package entities

import (
	"database/sql"
	"fmt"
	"time"
)

//go:generate easyjson position.go

// easyjson:json
type Position struct {
	Id         int64
	TraderUID  string  `db:"trader_uid"`
	Symbol     string  `json:"symbol"`
	EntryPrice float64 `json:"entryPrice" db:"entry_price"`
	Pnl        float64
	Roe        float64
	Amount     float64
	Invested   float64
	Leverage   uint8
	MarginMode MarginMode   `json:"-" db:"margin_mode"`
	Opened     bool         `json:"-"`
	Long       bool         `json:"-"`
	Exchange   Exchange     `json:"-"`
	Hedged     bool         `json:"-"`
	CreatedAt  time.Time    `json:"createTimestamp" db:"created_at"`
	UpdatedAt  time.Time    `json:"updateTimestamp" db:"updated_at"`
	ClosedAt   sql.NullTime `json:"closedTimestamp" db:"closed_at"`
}

func (p *Position) Key() string {
	return fmt.Sprintf("%s:%v:%d", p.Symbol, p.Long, p.MarginMode)
}
