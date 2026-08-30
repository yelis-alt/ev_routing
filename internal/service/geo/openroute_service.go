package geo

import (
	"bytes"
	"encoding/json"
	"ev_routing/internal/dto"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	unitsKilometers   = "km"
	preferenceFastest = "fastest"

	rateLimitMaxRetries  = 8
	rateLimitBaseBackoff = 2 * time.Second
	rateLimitMaxBackoff  = 30 * time.Second
)

// Service is a client for the OpenRouteService directions API.
type Service struct {
	requestURL string
	apiKey     string
	client     *http.Client
}

// NewService builds a Service that sends requests to requestURL,
// authenticated with apiKey. Its client's connection pool is sized above
// the default so RoutingService can fetch many routes concurrently without
// serializing on connection reuse.
func NewService(requestURL, apiKey string) *Service {
	return &Service{
		requestURL: requestURL,
		apiKey:     apiKey,
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        adjacencyMatrixConcurrency,
				MaxIdleConnsPerHost: adjacencyMatrixConcurrency,
			},
		},
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

	var routeResponse dto.RouteResponse

	backoff := rateLimitBaseBackoff
	for attempt := 0; ; attempt++ {
		req, err := http.NewRequest(
			http.MethodPost,
			s.requestURL,
			bytes.NewBuffer(body),
		)
		if err != nil {
			return nil, fmt.Errorf("create route request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", s.apiKey)

		resp, err := s.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("send route request: %w", err)
		}

		if resp.StatusCode == http.StatusTooManyRequests && attempt < rateLimitMaxRetries {
			wait := retryAfterDelay(resp.Header.Get("Retry-After"), backoff)
			_ = resp.Body.Close()
			time.Sleep(wait)
			if backoff < rateLimitMaxBackoff {
				backoff *= 2
			}
			continue
		}

		if resp.StatusCode != http.StatusOK {
			_ = resp.Body.Close()
			return nil, fmt.Errorf(
				"openrouteservice returned status %d",
				resp.StatusCode,
			)
		}

		decodeErr := json.NewDecoder(resp.Body).Decode(&routeResponse)
		_ = resp.Body.Close()
		if decodeErr != nil {
			return nil, fmt.Errorf("decode route response: %w", decodeErr)
		}

		break
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

// retryAfterDelay honors a Retry-After header (seconds) if present and
// parseable, otherwise falls back to the given backoff duration.
func retryAfterDelay(retryAfterHeader string, backoff time.Duration) time.Duration {
	if seconds, err := strconv.Atoi(retryAfterHeader); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}

	return backoff
}

// roundToTwoDecimals rounds value to two decimal places.
func roundToTwoDecimals(value float64) float64 {
	return float64(int(value*100+0.5)) / 100
}
