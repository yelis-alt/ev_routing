package dto

// RouteRequest is the OpenRouteService directions request for one from/to
// pair.
type RouteRequest struct {
	Coordinates      [][]float64 `json:"coordinates"` // [[fromLon, fromLat], [toLon, toLat]]
	SuppressWarnings bool        `json:"suppress_warnings,omitempty"`
	Units            string      `json:"units,omitempty"`
	Instructions     bool        `json:"instructions,omitempty"`
	Preference       string      `json:"preference,omitempty"`
	Geometry         bool        `json:"geometry,omitempty"`
}

// RouteResponse is OpenRouteService's directions response; only the
// summary of the (single) requested route is decoded.
type RouteResponse struct {
	Routes []struct {
		Summary struct {
			Distance float64 `json:"distance"` // meters
			Duration float64 `json:"duration"` // seconds
		} `json:"summary"`
	} `json:"routes"`
}

// RouteResult is a RouteResponse normalized into RoutingService's units.
type RouteResult struct {
	Distance     float64 `json:"distance"`     // S_ij, km
	TripDuration float64 `json:"tripDuration"` // hours
}
