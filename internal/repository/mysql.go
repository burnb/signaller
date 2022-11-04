package repository

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/burnb/signaller/internal/repository/entities"
)

type Mysql struct {
	db *sqlx.DB
}

func NewMysql(db *sqlx.DB) *Mysql {
	return &Mysql{db: db}
}

func (r *Mysql) DB() *sqlx.DB {
	return r.db
}

func (r *Mysql) User(uid string) (*entities.Trader, error) {
	user := entities.Trader{}
	err := r.db.Get(
		&user,
		`SELECT 
    				uid, nickName, userPhotoUrl, followerCount, pnlValue, roiValue,
    				weeklyPnl, weeklyRoe, monthlyPnl, monthlyRoe  
				FROM traders
				WHERE uid = ?`,
		uid,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Mysql) PositionById(id int64) (*entities.Position, error) {
	position := entities.Position{Id: id}
	err := r.db.Get(&position, "SELECT id, uid, symbol, entryPrice, markPrice, pnl, roe, amount, leverage, invested, opened, `long`, unix_timestamp(updatedAt) AS updatedAt, exchange, margin_mode FROM trader_positions WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &position, nil
}

func (r *Mysql) FirstTraderLogAt() (*time.Time, error) {
	var at time.Time
	if err := r.db.Get(&at, "SELECT createdAt FROM traders_log LIMIT 1"); err != nil {
		return nil, err
	}

	return &at, nil
}

func (r *Mysql) OpenedPositions(trader *entities.Trader) ([]*entities.Position, error) {
	return nil, nil
}

func (r *Mysql) TradersWithSub() ([]*entities.Trader, error) {
	return nil, nil
}

func (r *Mysql) UpdatePosition(position *entities.Position) error {
	return nil
}

func (r *Mysql) ClosePosition(position *entities.Position) error {
	return nil
}
