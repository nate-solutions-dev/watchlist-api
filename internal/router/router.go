package router

import (
	"github.com/Friel909/watchlist-api/config"
	"github.com/Friel909/watchlist-api/internal/controller"
	"github.com/Friel909/watchlist-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

// NewRouter registers middlewares and all public/private API routes.
func NewRouter(cfg *config.Config, healthController *controller.HealthController, authController *controller.AuthController, watchListController *controller.WatchListController, discoverController *controller.DiscoverController) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID(), gin.Logger(), gin.Recovery())

	r.GET("/health", healthController.Health)

	public := r.Group("/public")
	{
		public.POST("/auth/register", authController.Register)
		public.POST("/auth/login", authController.Login)
	}

	private := r.Group("/private")
	private.Use(middleware.JWTMiddleware(cfg))
	{
		private.GET("/auth/me", authController.Me)
		private.GET("/watchlist", watchListController.GetAll)
		private.POST("/watchlist", watchListController.Create)
		private.PATCH("/watchlist/:id", watchListController.Update)
		private.DELETE("/watchlist/:id", watchListController.Delete)
		private.GET("/search", discoverController.SearchTitles)
		private.GET("/discover/trending", discoverController.GetTrending)
	}

	return r
}
