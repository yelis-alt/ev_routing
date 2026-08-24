package dto

import "time"

const (
	PlugType2   Plug = "TYPE_2"
	PlugCCS     Plug = "CCS"
	PlugCHAdeMO Plug = "CHADEMO"
)

// Plug is a charging connector standard.
type Plug string

// Edge is a directed i->j edge of the adjacency matrix.
type Edge struct {
	Distance       float64       `json:"distance"`       // S_ij, km
	TripDuration   time.Duration `json:"tripDuration"`   // driving time i->j
	ChargeDuration time.Duration `json:"chargeDuration"` // t_jm^charge at j's chosen slot
	Cost           float64       `json:"cost"`           // Z_s + (if charging) Z_w + Z_ch
	SlotId         int           `json:"slotId"`         // m, unique within j only
}

// CoordsDTO is a WGS84 geographic point.
type CoordsDTO struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// RouteRequestDTO is the request body for every POST /route/* endpoint.
type RouteRequestDTO struct {
	StartCoords      CoordsDTO    `json:"startCoords"`      // S
	FinishCoords     CoordsDTO    `json:"finishCoords"`     // D
	AccLevel         float64      `json:"accLevel"`         // C_start, kWh
	AccMax           float64      `json:"accMax"`           // C_max, kWh
	SpendOpt         float64      `json:"spendOpt"`         // E0, passport consumption at 25°C, kWh/100km
	Temperature      float64      `json:"temperature"`      // T, °C
	FilteredStations []StationDTO `json:"filteredStations"` // pool for C_k
}

// StationDTO is vertex j.
type StationDTO struct {
	Id     int       `json:"id"`
	Coords CoordsDTO `json:"coords"`
	Plug   Plug      `json:"plug"`
	Slots  []SlotDTO `json:"slots"` // M_j slots, indexed by m
}

// SlotDTO is slot m at station j.
type SlotDTO struct {
	Id       int     `json:"id"`       // m, unique within the station only
	Price    float64 `json:"price"`    // r_jm, currency/kWh
	Power    float64 `json:"power"`    // P_jm, kW
	WaitTime float64 `json:"waitTime"` // t_jm^wait, hours
	IsActive bool    `json:"isActive"` // excludes the slot from x_jm selection when false
}

// RouteNodeDTO is one stop on the winning path.
type RouteNodeDTO struct {
	RouteNode      StationDTO    `json:"routeNode"`      // vertex j
	Distance       float64       `json:"distance"`       // cumulative, km
	Cost           float64       `json:"cost"`           // cumulative Z
	ChargeDuration time.Duration `json:"chargeDuration"` // t_jm^charge at this stop
	ReachDuration  time.Duration `json:"reachDuration"`  // cumulative drive+charge+wait time
	SlotId         int           `json:"slotId"`         // m, unique within RouteNode.Id only
}
