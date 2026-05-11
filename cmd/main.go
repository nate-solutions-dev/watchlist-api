package main

import (
	"context"
	"log"

	"github.com/Friel909/watchlist-api/config"
	"github.com/Friel909/watchlist-api/internal/controller"
	"github.com/Friel909/watchlist-api/internal/database"
	"github.com/Friel909/watchlist-api/internal/logger"
	"github.com/Friel909/watchlist-api/internal/repository"
	"github.com/Friel909/watchlist-api/internal/router"
	"github.com/Friel909/watchlist-api/internal/service"
)

// @title           Watchlist API
// @version         1.0
// @description     Movie and show watchlist API
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// main boots dependencies and starts the HTTP server.
func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	logger.Init(cfg.Environment)

	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("init postgres pool: %v", err)
	}
	defer pool.Close()

	authRepo := repository.NewAuthRepository(pool)
	watchListRepo := repository.NewWatchListRepository(pool)

	tmdbService := service.NewTMDBService(cfg)
	authService := service.NewAuthService(authRepo, cfg, tmdbService)
	watchListService := service.NewWatchListService(watchListRepo, tmdbService)

	healthController := controller.NewHealthController()
	authController := controller.NewAuthController(authService)
	watchListController := controller.NewWatchListController(watchListService)
	discoverController := controller.NewDiscoverController(tmdbService)

	r := router.NewRouter(cfg, healthController, authController, watchListController, discoverController)

	logger.Info(context.Background(), "main", "server starting", "port", cfg.Port, "env", cfg.Environment)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
