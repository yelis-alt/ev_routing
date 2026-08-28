package additional

import (
	"math"

	"ev_routing/internal/dto"
)

const (
	earthRadiusKm = 6371.0
	kmPerDegree   = 111.32
)

// haversineDistanceKm is the great-circle distance between a and b, km.
func haversineDistanceKm(a, b dto.CoordsDTO) float64 {
	lat1 := a.Latitude * math.Pi / 180
	lat2 := b.Latitude * math.Pi / 180
	dLat := (b.Latitude - a.Latitude) * math.Pi / 180
	dLon := (b.Longitude - a.Longitude) * math.Pi / 180

	h := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)

	return 2 * earthRadiusKm * math.Asin(math.Sqrt(h))
}

// distancePointToSegmentKm approximates the distance from p to segment
// a-b, km, via flat-plane projection; not a true geodesic distance.
func distancePointToSegmentKm(p, a, b dto.CoordsDTO) float64 {
	latScale := math.Cos(a.Latitude * math.Pi / 180)

	toXY := func(c dto.CoordsDTO) (float64, float64) {
		x := (c.Longitude - a.Longitude) * latScale * kmPerDegree
		y := (c.Latitude - a.Latitude) * kmPerDegree
		return x, y
	}

	px, py := toXY(p)
	bx, by := toXY(b)

	abLen2 := bx*bx + by*by
	if abLen2 == 0 {
		return haversineDistanceKm(p, a)
	}

	t := (px*bx + py*by) / abLen2
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	dx := px - t*bx
	dy := py - t*by

	return math.Sqrt(dx*dx + dy*dy)
}
