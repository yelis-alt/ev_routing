package dto

import "time"

const (
	PlugType2   Plug = "TYPE_2"
	PlugCCS     Plug = "CCS"
	PlugCHAdeMO Plug = "CHADEMO"
)

type Plug string

type Edge struct {
	Distance       float64       `json:"distance"`
	TripDuration   time.Duration `json:"tripDuration"`
	ChargeDuration time.Duration `json:"chargeDuration"`
	Cost           float64       `json:"cost"`
}

type CoordsDTO struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type RouteRequestDTO struct {
	StartCoords      CoordsDTO    `json:"startCoords"`
	FinishCoords     CoordsDTO    `json:"finishCoords"`
	AccLevel         float64      `json:"accLevel"`
	AccMax           float64      `json:"accMax"`
	SpendOpt         float64      `json:"spendOpt"`
	Temperature      float64      `json:"temperature"`
	FilteredStations []StationDTO `json:"filteredStations"`
}

type StationDTO struct {
	Id       int       `json:"id"`
	Coords   CoordsDTO `json:"coords"`
	Price    int       `json:"price"`
	Plug     Plug      `json:"plug"`
	Power    int       `json:"power"`
	IsActive bool      `json:"isActive"`
}

type RouteNodeDTO struct {
	RouteNode      StationDTO    `json:"routeNode"`
	Distance       float64       `json:"distance"`
	Cost           float64       `json:"cost"`
	ChargeDuration time.Duration `json:"chargeDuration"`
	ReachDuration  time.Duration `json:"reachDuration"`
}
