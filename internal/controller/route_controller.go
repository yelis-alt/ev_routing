package controller

import (
	"encoding/json"
	"net/http"

	"ev_routing/internal/dto"
	"ev_routing/internal/service"
)

// RouteController exposes an HTTP endpoint per route-search strategy.
type RouteController struct {
	routing        *service.RoutingService
	genetic        *service.GeneticRouteService
	dijkstra       *service.DijkstraRouteService
	vns            *service.VNSRouteService
	branchAndBound *service.BranchAndBoundRouteService
	aco            *service.ACORouteService
}

// NewRouteController builds a RouteController backed by the given services.
func NewRouteController(
	routing *service.RoutingService,
	genetic *service.GeneticRouteService,
	dijkstra *service.DijkstraRouteService,
	vns *service.VNSRouteService,
	branchAndBound *service.BranchAndBoundRouteService,
	aco *service.ACORouteService,
) *RouteController {
	return &RouteController{
		routing:        routing,
		genetic:        genetic,
		dijkstra:       dijkstra,
		vns:            vns,
		branchAndBound: branchAndBound,
		aco:            aco,
	}
}

// RegisterRoutes wires the controller's endpoints onto mux.
func (rc *RouteController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /route/genetic", rc.handleGenetic)
	mux.HandleFunc("POST /route/dijkstra", rc.handleDijkstra)
	mux.HandleFunc("POST /route/vns", rc.handleVNS)
	mux.HandleFunc("POST /route/branch-and-bound", rc.handleBranchAndBound)
	mux.HandleFunc("POST /route/aco", rc.handleACO)
}

// handleGenetic finds the cheapest route using GeneticRouteService.
func (rc *RouteController) handleGenetic(w http.ResponseWriter, r *http.Request) {
	routeRequest, adjacencyMatrix, ok := rc.decodeAndBuildMatrix(w, r)
	if !ok {
		return
	}

	routeNodes := rc.genetic.GetRouteWithEvolution(adjacencyMatrix, routeRequest)
	writeJSON(w, http.StatusOK, routeNodes)
}

// handleDijkstra finds the cheapest route using DijkstraRouteService.
func (rc *RouteController) handleDijkstra(w http.ResponseWriter, r *http.Request) {
	routeRequest, adjacencyMatrix, ok := rc.decodeAndBuildMatrix(w, r)
	if !ok {
		return
	}

	routeNodes := rc.dijkstra.GetRouteWithDijkstra(adjacencyMatrix, routeRequest)
	writeJSON(w, http.StatusOK, routeNodes)
}

// handleVNS finds the cheapest route using VNSRouteService.
func (rc *RouteController) handleVNS(w http.ResponseWriter, r *http.Request) {
	routeRequest, adjacencyMatrix, ok := rc.decodeAndBuildMatrix(w, r)
	if !ok {
		return
	}

	routeNodes := rc.vns.GetRouteWithVNS(adjacencyMatrix, routeRequest)
	writeJSON(w, http.StatusOK, routeNodes)
}

// handleBranchAndBound finds the cheapest route using BranchAndBoundRouteService.
func (rc *RouteController) handleBranchAndBound(w http.ResponseWriter, r *http.Request) {
	routeRequest, adjacencyMatrix, ok := rc.decodeAndBuildMatrix(w, r)
	if !ok {
		return
	}

	routeNodes := rc.branchAndBound.GetRouteWithBranchAndBound(adjacencyMatrix, routeRequest)
	writeJSON(w, http.StatusOK, routeNodes)
}

// handleACO finds the cheapest route using ACORouteService.
func (rc *RouteController) handleACO(w http.ResponseWriter, r *http.Request) {
	routeRequest, adjacencyMatrix, ok := rc.decodeAndBuildMatrix(w, r)
	if !ok {
		return
	}

	routeNodes := rc.aco.GetRouteWithACO(adjacencyMatrix, routeRequest)
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

	adjacencyMatrix, err := rc.routing.GetAdjacencyMatrix(&routeRequest)
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
