package controller

import (
	"encoding/json"
	"net/http"

	"ev_routing/internal/dto"
	"ev_routing/internal/service"
)

// RouteController exposes HTTP endpoints for the genetic and Dijkstra
// cheapest-path search approaches.
type RouteController struct {
	Routing  *service.RoutingService
	Genetic  *service.GeneticRouteService
	Dijkstra *service.DijkstraRouteService
}

// NewRouteController builds a RouteController backed by the given services.
func NewRouteController(
	routing *service.RoutingService,
	genetic *service.GeneticRouteService,
	dijkstra *service.DijkstraRouteService,
) *RouteController {
	return &RouteController{
		Routing:  routing,
		Genetic:  genetic,
		Dijkstra: dijkstra,
	}
}

// RegisterRoutes wires the controller's endpoints onto mux.
func (rc *RouteController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /route/genetic", rc.handleGenetic)
	mux.HandleFunc("POST /route/dijkstra", rc.handleDijkstra)
}

// handleGenetic finds the cheapest route using GeneticRouteService.
func (rc *RouteController) handleGenetic(w http.ResponseWriter, r *http.Request) {
	routeRequest, adjacencyMatrix, ok := rc.decodeAndBuildMatrix(w, r)
	if !ok {
		return
	}

	routeNodes := rc.Genetic.GetRouteWithEvolution(adjacencyMatrix, routeRequest)
	writeJSON(w, http.StatusOK, routeNodes)
}

// handleDijkstra finds the cheapest route using DijkstraRouteService.
func (rc *RouteController) handleDijkstra(w http.ResponseWriter, r *http.Request) {
	routeRequest, adjacencyMatrix, ok := rc.decodeAndBuildMatrix(w, r)
	if !ok {
		return
	}

	routeNodes := rc.Dijkstra.GetRouteWithDijkstra(adjacencyMatrix, routeRequest)
	writeJSON(w, http.StatusOK, routeNodes)
}

// decodeAndBuildMatrix decodes the request body and builds its adjacency
// matrix; on failure it writes the error response itself and returns ok=false.
func (rc *RouteController) decodeAndBuildMatrix(
	w http.ResponseWriter,
	r *http.Request,
) (*dto.RouteRequestDTO, map[int]map[int]dto.Edge, bool) {
	var routeRequest dto.RouteRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&routeRequest); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return nil, nil, false
	}

	adjacencyMatrix, err := rc.Routing.GetAdjacencyMatrix(&routeRequest)
	if err != nil {
		http.Error(w, "failed to build adjacency matrix: "+err.Error(), http.StatusBadGateway)
		return nil, nil, false
	}

	return &routeRequest, adjacencyMatrix, true
}

// writeJSON writes payload as a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
