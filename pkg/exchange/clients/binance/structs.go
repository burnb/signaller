package binance

import (
	"github.com/burnb/signaller/internal/repository/entities"
)

//go:generate easyjson structs.go

type SymbolInfo struct {
	Symbol            string
	Quote             string
	Base              string
	PricePrecision    int
	MinQuantity       float64 // Количество токенов
	MinSize           float64 // Размер позиции в USDT
	QuantityPrecision int
}

// easyjson:json
type TraderResponse struct {
	Code          string             `json:"code"`
	Message       string             `json:"message"`
	MessageDetail string             `json:"messageDetail"`
	Data          []*entities.Trader `json:"data"`
	Success       bool               `json:"success"`
}

// easyjson:json
type BaseInfo struct {
	DeliveryPositionShared bool   `json:"deliveryPositionShared"`
	FollowerCount          int    `json:"followerCount"`
	FollowingCount         int    `json:"followingCount"`
	Introduction           string `json:"introduction"`
	NickName               string `json:"nickName"`
	PositionShared         bool   `json:"positionShared"`
	TwitterUrl             string `json:"twitterUrl"`
	UserPhotoUrl           string `json:"userPhotoUrl"`
}

// easyjson:json
type BaseInfoResponse struct {
	Code          string    `json:"code"`
	Message       string    `json:"message"`
	MessageDetail string    `json:"messageDetail"`
	Data          *BaseInfo `json:"data"`
	Success       bool      `json:"success"`
}

// easyjson:json
type Stats struct {
	PeriodType     string  `json:"periodType"`
	StatisticsType string  `json:"statisticsType"`
	Value          float64 `json:"value"`
	Rank           int     `json:"rank"`
}

// easyjson:json
type StatsResponse struct {
	Code          string   `json:"code"`
	Message       string   `json:"message"`
	MessageDetail string   `json:"messageDetail"`
	Data          []*Stats `json:"data"`
	Success       bool     `json:"success"`
}

// easyjson:json
type PositionsResponse struct {
	Code          string `json:"code"`
	Message       string `json:"message"`
	MessageDetail string `json:"messageDetail"`
	Data          struct {
		OtherPositionRetList []*entities.Position `json:"otherPositionRetList"`
	} `json:"data"`
	Success bool `json:"success"`
}
