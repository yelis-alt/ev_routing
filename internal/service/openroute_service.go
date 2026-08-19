package service

import (
	"bytes"
	"encoding/json"
	"ev_routing/internal/dto"
	"fmt"
	"io"
	"net/http"
)

const (
	unitsKilometers   = "km"
	preferenceFastest = "fastest"
)

// Service is a client for the OpenRouteService directions API.
type Service struct {
	RequestURL string
	APIKey     string
	Client     *http.Client
}

// NewService builds a Service that sends requests to requestURL,
// authenticated with apiKey.
func NewService(requestURL, apiKey string) *Service {
	return &Service{
		RequestURL: requestURL,
		APIKey:     apiKey,
		Client:     &http.Client{},
	}
}

// GetRoute fetches the fastest driving route between two coordinates from
// OpenRouteService and returns its distance and trip duration.
func (s *Service) GetRoute(
	startLat float64,
	startLong float64,
	finishLat float64,
	finishLong float64,
) (*dto.RouteResult, error) {
	request := dto.RouteRequest{
		Coordinates: [][]float64{
			{startLong, startLat},
			{finishLong, finishLat},
		},
		SuppressWarnings: true,
		Units:            unitsKilometers,
		Instructions:     false,
		Preference:       preferenceFastest,
		Geometry:         false,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal route request: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		s.RequestURL,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("create route request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.APIKey)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send route request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"openrouteservice returned status %d",
			resp.StatusCode,
		)
	}

	var routeResponse dto.RouteResponse

	if err := json.NewDecoder(resp.Body).Decode(&routeResponse); err != nil {
		return nil, fmt.Errorf("decode route response: %w", err)
	}

	if len(routeResponse.Routes) == 0 {
		return nil, fmt.Errorf("no routes returned by openrouteservice")
	}

	summary := routeResponse.Routes[0].Summary

	return &dto.RouteResult{
		Distance:     roundToTwoDecimals(summary.Distance),
		TripDuration: roundToTwoDecimals(summary.Duration / 3600),
	}, nil
}

// roundToTwoDecimals rounds value to two decimal places.
func roundToTwoDecimals(value float64) float64 {
	return float64(int(value*100+0.5)) / 100
}
