package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"InteractiveScraper/internal/config"
	"InteractiveScraper/internal/scraper"
	"InteractiveScraper/internal/storage"
)

func main() {
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic("Log klasörü oluşturulamadı: " + err.Error())
	}

	logFile, err := os.OpenFile("logs/scraper.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Log dosyası açılamadı: " + err.Error())
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	logger := slog.New(slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
		Level: getLogLevel(os.Getenv("LOG_LEVEL")),
	}))

	slog.SetDefault(logger)

	logger.Info("starting CTI scraper")

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatal("config validation failed:", err)
	}

	logger.Info("configuration loaded",
		slog.String("tor_proxy", cfg.TorProxy),
		slog.String("db_host", cfg.DBHost),
		slog.Duration("scraper_interval", cfg.ScraperInterval),
	)

	store, err := storage.NewPostgresStorage(cfg.GetDSN(), logger)
	if err != nil {
		log.Fatal("failed to initialize storage:", err)
	}
	defer store.Close()

	scraperInstance, err := scraper.NewScraper(cfg, store, logger)
	if err != nil {
		log.Fatal("failed to initialize scraper:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	healthChecker := scraper.NewHealthChecker(scraperInstance, logger)
	go startHealthServer(cfg.HealthPort, healthChecker, logger)

	go handleShutdown(cancel, logger)

	logger.Info("scraper running")
	if err := scraperInstance.Run(ctx); err != nil && err != context.Canceled {
		logger.Error("scraper error", slog.String("error", err.Error()))
	}

	logger.Info("scraper stopped")
}

func startHealthServer(port string, healthChecker *scraper.HealthChecker, logger *slog.Logger) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthChecker.Handler())

	addr := ":" + port
	logger.Info("health check server starting", slog.String("addr", addr))

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("health server error", slog.String("error", err.Error()))
	}
}

func handleShutdown(cancel context.CancelFunc, logger *slog.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	logger.Info("shutdown signal received", slog.String("signal", sig.String()))

	cancel()

	time.Sleep(30 * time.Second)
}

func getLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
