package main

import (
	"log"
	"net/http"
	"os"

	"ev_routing/config"
	"ev_routing/internal/controller"
	"ev_routing/internal/service"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	apiKey := os.Getenv("OPENROUTESERVICE_API_KEY")

	orsService := service.NewService(cfg.OpenRouteService.Request.URL, apiKey)
	routingService := service.NewRoutingService(orsService)
	geneticRouteService := service.NewGeneticRouteService()
	dijkstraRouteService := service.NewDijkstraRouteService()

	routeController := controller.NewRouteController(routingService, geneticRouteService, dijkstraRouteService)

	mux := http.NewServeMux()
	routeController.RegisterRoutes(mux)

	log.Println("Application started, listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
