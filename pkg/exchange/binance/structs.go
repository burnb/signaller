package binance

import (
	"github.com/burnb/signaller/internal/repository/entities"
)

//go:generate easyjson -all structs.go

type SymbolInfo struct {
	Symbol            string
	Quote             string
	Base              string
	PricePrecision    int
	MinQuantity       float64
	MinSize           float64
	QuantityPrecision int
}

type TraderResponse struct {
	Code          string             `json:"code"`
	Message       string             `json:"message"`
	MessageDetail string             `json:"messageDetail"`
	Data          []*entities.Trader `json:"data"`
	Success       bool               `json:"success"`
}

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

type BaseInfoResponse struct {
	Code          string    `json:"code"`
	Message       string    `json:"message"`
	MessageDetail string    `json:"messageDetail"`
	Data          *BaseInfo `json:"data"`
	Success       bool      `json:"success"`
}

type Stats struct {
	PeriodType     string  `json:"periodType"`
	StatisticsType string  `json:"statisticsType"`
	Value          float64 `json:"value"`
	Rank           int     `json:"rank"`
}

type StatsResponse struct {
	Code          string   `json:"code"`
	Message       string   `json:"message"`
	MessageDetail string   `json:"messageDetail"`
	Data          []*Stats `json:"data"`
	Success       bool     `json:"success"`
}

type PositionsResponse struct {
	Code          string `json:"code"`
	Message       string `json:"message"`
	MessageDetail string `json:"messageDetail"`
	Data          struct {
		OtherPositionRetList []*entities.Position `json:"otherPositionRetList"`
	} `json:"data"`
	Success bool `json:"success"`
}
