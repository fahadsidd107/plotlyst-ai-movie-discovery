package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"plotlyst/server/apiresponse"
	"plotlyst/server/config"
	"plotlyst/server/controllers"
	"plotlyst/server/database"
	"plotlyst/server/middleware"
	"plotlyst/server/routes"
	"plotlyst/server/tmdb"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	connectCtx, cancelConnect := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := database.Connect(connectCtx, cfg.MongoDBURI)
	cancelConnect()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	indexCtx, cancelIndexes := context.WithTimeout(context.Background(), 15*time.Second)
	if err := database.EnsureIndexes(indexCtx, client); err != nil {
		cancelIndexes()
		log.Fatalf("database index setup failed: %v", err)
	}
	cancelIndexes()

	router := gin.New()
	metrics := middleware.NewMetrics()
	router.Use(gin.Recovery(), middleware.RequestID(), middleware.SecurityHeaders(), metrics.Middleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/hello", func(c *gin.Context) { c.String(http.StatusOK, "Hello, Plotlyst!") })
	router.GET("/health/live", controllers.Liveness)
	router.GET("/health/ready", controllers.Readiness(client))
	router.GET("/metrics", metrics.Handler)
	router.NoRoute(func(c *gin.Context) {
		apiresponse.Error(c, http.StatusNotFound, "route_not_found", "The requested API route does not exist")
	})
	routes.SetupUnProtectedRoutes(router, client)
	tmdbClient := tmdb.New(cfg.TMDBBaseURL, cfg.TMDBAPIKey, cfg.TMDBTimeout, cfg.TMDBCacheTTL, cfg.TMDBMaxRetries)
	routes.SetupCatalogRoutes(router, controllers.NewCatalogController(tmdbClient, cfg.TMDBImageBaseURL))
	routes.SetupSpecialistRoutes(router, client, tmdbClient, cfg)
	routes.SetupIntelligenceRoutes(router, client, tmdbClient, cfg)
	routes.SetupStoryDNARoutes(router, client, tmdbClient)
	routes.SetupProtectedRoutes(router, client)
	routes.SetupCommunityRoutes(router, client)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Plotlyst API listening on http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, os.Interrupt, syscall.SIGTERM)
	select {
	case sig := <-shutdownSignals:
		log.Printf("shutdown signal received: %s", sig)
	case err := <-serverErrors:
		log.Printf("HTTP server failed: %v", err)
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancelShutdown()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP shutdown incomplete: %v", err)
	}
	if err := client.Disconnect(shutdownCtx); err != nil {
		log.Printf("database disconnect failed: %v", err)
	}
	log.Println("Plotlyst API stopped")
}
