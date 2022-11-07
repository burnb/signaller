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

func (r *Mysql) Trader(uid string) (*entities.Trader, error) {
	var trader entities.Trader
	err := r.db.Get(&trader, `SELECT * FROM traders WHERE uid=?`, uid)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &trader, nil
}

func (r *Mysql) CreateTrader(trader *entities.Trader) error {
	trader.CreatedAt = time.Now()
	trader.UpdatedAt = trader.CreatedAt
	_, err :=
		r.db.NamedExec(
			`INSERT INTO traders(
						uid, pnl, roi, roi_weekly, pnl_weekly, roi_monthly, pnl_monthly, roi_yearly, pnl_yearly, 
						position_shared, publisher, published_at, created_at, updated_at
	    			) VALUES (
						:uid, :pnl, :roi, :roi_weekly, :pnl_weekly, :roi_monthly, :pnl_monthly, :roi_yearly, :pnl_yearly, 
						:position_shared, :publisher, :published_at, :created_at, :updated_at
					)`,
			trader,
		)

	return err
}

func (r *Mysql) UpdateTrader(trader *entities.Trader) error {
	trader.UpdatedAt = time.Now()
	_, err :=
		r.db.NamedExec(
			`UPDATE traders 
					SET 
						pnl=:pnl, roi=:roi, roi_weekly=:roi_weekly, pnl_weekly=:pnl_weekly, 
						roi_monthly=:roi_monthly, pnl_monthly=:pnl_monthly, roi_yearly=:roi_yearly, pnl_yearly=:pnl_yearly, 
						position_shared=:position_shared, publisher=:publisher, 
						published_at=:published_at, updated_at=:updated_at
					WHERE uid = :uid`,
			trader,
		)

	return err
}

func (r *Mysql) RefreshPublishTime(uid string) error {
	now := time.Now()
	_, err := r.db.NamedExec(
		`UPDATE traders SET updated_at=:updated_at, published_at=:published_at WHERE uid = :uid`,
		map[string]interface{}{
			"uid":          uid,
			"updated_at":   now,
			"published_at": now,
		},
	)

	return err
}

func (r *Mysql) Publishers() ([]*entities.Trader, error) {
	var traders []*entities.Trader
	err := r.db.Select(&traders, `SELECT * FROM traders WHERE publisher = 1`)
	if err != nil {
		return nil, err
	}

	return traders, nil
}

func (r *Mysql) OpenedPositions(trader *entities.Trader) ([]*entities.Position, error) {
	var positions []*entities.Position
	err := r.db.Select(&positions, `SELECT * FROM positions WHERE trader_uid = ? AND closed_at IS NULL`, trader.Uid)
	if err != nil {
		return nil, err
	}

	return positions, nil
}

func (r *Mysql) CreatePosition(position *entities.Position) error {
	if position.CreateTimestamp == 0 {
		position.CreatedAt = time.Now()
	}
	res, err :=
		r.db.NamedExec(
			`INSERT INTO positions(
					 trader_uid, symbol, positions.long, entry_price, pnl, roe, amount, leverage,
					 exchange, margin_mode, hedged, created_at, updated_at, closed_at
	    			) VALUES (
						:trader_uid, :symbol, :long, :entry_price, :pnl, :roe, :amount, :leverage,
					 	:exchange, :margin_mode, :hedged, :created_at, :updated_at, :closed_at
					)`,
			position,
		)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	position.Id = id

	return nil
}

func (r *Mysql) UpdatePosition(position *entities.Position) error {
	_, err := r.db.NamedExec(
		`UPDATE positions 
				SET 
				     trader_uid=:trader_uid, symbol=:symbol, positions.long=:long, entry_price=:entry_price, pnl=:pnl, roe=:roe, 
				     amount=:amount, leverage=:leverage, exchange=:exchange, 
				     margin_mode=:margin_mode, hedged=:hedged, updated_at=:updated_at, closed_at=:closed_at
				WHERE id = :id`,
		position,
	)

	return err
}
