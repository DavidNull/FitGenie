// @title FitGenie API
// @version 1.0
// @description AI-powered outfit recommendation API with Flutter mobile app
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.fitgenie.local/support
// @contact.email support@fitgenie.local
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fitgenie/internal/api"
	"fitgenie/internal/config"
	"fitgenie/pkg/database"
	"fitgenie/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "fitgenie/docs"
)

func main() {
	log := logger.NewLogger()

	cfg := config.Load()

	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		log.Fatal("failed to run migrations", "error", err)
	}

	router := api.NewRouter(db, log, cfg)

	mux := http.NewServeMux()
	mux.Handle("/", router)
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        mux,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Info("starting server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed to start", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", "error", err)
	}

	log.Info("server exited")
}
