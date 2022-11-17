package entities

import (
	"database/sql"
	"fmt"
)

//go:generate easyjson -all position.go

type Position struct {
	Id         int64        `json:"-"`
	TraderUID  string       `json:"-" db:"trader_uid"`
	Symbol     string       `json:"symbol"`
	EntryPrice float64      `json:"entryPrice" db:"entry_price"`
	Pnl        float64      `json:"pnl"`
	Roe        float64      `json:"roe"`
	Amount     float64      `json:"amount"`
	Leverage   uint8        `json:"leverage"`
	MarginMode MarginMode   `json:"-" db:"margin_mode"`
	Long       bool         `json:"-"`
	Exchange   Exchange     `json:"-"`
	Hedged     bool         `json:"-"`
	CreatedAt  int64        `json:"createTimeStamp" db:"created_at"`
	UpdatedAt  int64        `json:"updateTimeStamp" db:"updated_at"`
	ClosedAt   sql.NullTime `json:"-" db:"closed_at"`
}

func (p *Position) Key() string {
	return fmt.Sprintf("%s:%v:%d", p.Symbol, p.Long, p.MarginMode)
}
