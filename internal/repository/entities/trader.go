package entities

//go:generate easyjson trader.go

// easyjson:json
type Trader struct {
	EncryptedUid           string      `db:"uid" json:"encryptedUid"`
	NickName               string      `db:"nickName" json:"nickName"`
	UserPhotoUrl           string      `db:"userPhotoUrl" json:"userPhotoUrl"`
	FollowerCount          int64       `db:"followerCount" json:"followerCount"`
	PnlValue               float64     `db:"pnlValue" json:"pnlValue"`
	RoiValue               float64     `db:"roiValue" json:"roiValue"`
	WeeklyRoi              float64     `db:"weeklyRoe" json:"weeklyRoi"`
	WeeklyPnl              float64     `db:"weeklyPnl" json:"weeklyPnl"`
	MonthlyRoi             float64     `db:"monthlyRoe" json:"monthlyRoi"`
	MonthlyPnl             float64     `db:"monthlyPnl" json:"monthlyPnl"`
	YearlyRoi              float64     `db:"yearlyRoe" json:"yearlyRoi"`
	YearlyPnl              float64     `db:"yearlyPnl" json:"yearlyPnl"`
	Rank                   int         `db:"rank"`
	PositionShared         bool        `json:"positionShared"`
	DeliveryPositionShared bool        `json:"deliveryPositionShared"`
	LastUpdate             int64       `db:"lastUpdate"`
	Positions              []*Position `json:"-"`
}
