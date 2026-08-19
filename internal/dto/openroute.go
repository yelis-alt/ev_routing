package dto

type RouteRequest struct {
	Coordinates      [][]float64 `json:"coordinates"`
	SuppressWarnings bool        `json:"suppress_warnings,omitempty"`
	Units            string      `json:"units,omitempty"`
	Instructions     bool        `json:"instructions,omitempty"`
	Preference       string      `json:"preference,omitempty"`
	Geometry         bool        `json:"geometry,omitempty"`
}

type RouteResponse struct {
	Routes []struct {
		Summary struct {
			Distance float64 `json:"distance"`
			Duration float64 `json:"duration"`
		} `json:"summary"`
	} `json:"routes"`
}

type RouteResult struct {
	Distance     float64 `json:"distance"`
	TripDuration float64 `json:"tripDuration"`
}
