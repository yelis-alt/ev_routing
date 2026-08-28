package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"ev_routing/config"
	"ev_routing/internal/controller"
	"ev_routing/internal/repository"
	"ev_routing/internal/service/algo"
	"ev_routing/internal/service/geo"
	"ev_routing/internal/service/schedule"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	apiKey := os.Getenv("OPENROUTESERVICE_API_KEY")

	orsService := geo.NewService(cfg.OpenRouteService.Request.URL, apiKey)
	routingService := geo.NewRoutingService(orsService)
	geneticRouteService := algo.NewGeneticRouteService()
	dijkstraRouteService := algo.NewDijkstraRouteService()
	vnsRouteService := algo.NewVNSRouteService()
	branchAndBoundRouteService := algo.NewBranchAndBoundRouteService()
	acoRouteService := algo.NewACORouteService()

	routeController := controller.NewRouteController(
		routingService,
		geneticRouteService,
		dijkstraRouteService,
		vnsRouteService,
		branchAndBoundRouteService,
		acoRouteService,
	)

	db, err := sql.Open("postgres", databaseDSN(cfg))
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	scheduleRepository := repository.NewScheduleRepository(db)
	scheduleService := schedule.NewScheduleService(scheduleRepository)
	scheduleController := controller.NewScheduleController(scheduleService)

	mux := http.NewServeMux()
	routeController.RegisterRoutes(mux)
	scheduleController.RegisterRoutes(mux)

	log.Println("Application started, listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// databaseDSN builds a Postgres connection string from cfg, with the
// password taken from DB_PASSWORD.
func databaseDSN(cfg *config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		os.Getenv("DB_PASSWORD"),
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
}
