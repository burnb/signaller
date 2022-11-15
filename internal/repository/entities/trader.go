package entities

import (
	"database/sql"
	"time"
)

//go:generate easyjson -all trader.go

type Trader struct {
	Uid            string               `db:"uid" json:"encryptedUid"`
	Pnl            float64              `db:"pnl" json:"pnlValue"`
	PnlWeekly      float64              `db:"pnl_weekly" json:"weeklyPnl"`
	PnlMonthly     float64              `db:"pnl_monthly" json:"monthlyPnl"`
	PnlYearly      float64              `db:"pnl_yearly" json:"yearlyPnl"`
	Roi            float64              `db:"roi" json:"roiValue"`
	RoiWeekly      float64              `db:"roi_weekly" json:"weeklyRoi"`
	RoiMonthly     float64              `db:"roi_monthly" json:"monthlyRoi"`
	RoiYearly      float64              `db:"roi_yearly" json:"yearlyRoi"`
	PositionShared bool                 `db:"position_shared" json:"positionShared"`
	Publisher      bool                 `db:"publisher" json:"-"`
	PublishedAt    sql.NullTime         `db:"published_at" json:"-"`
	CreatedAt      time.Time            `db:"created_at" json:"-"`
	UpdatedAt      time.Time            `db:"updated_at" json:"-"`
	Positions      map[string]*Position `json:"-"`
}
