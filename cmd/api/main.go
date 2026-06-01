package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/magnifimind/sentinel/internal/config"
	"github.com/magnifimind/sentinel/internal/handler"
	"github.com/magnifimind/sentinel/internal/kafka"
	"github.com/magnifimind/sentinel/internal/policy"
	"github.com/magnifimind/sentinel/internal/store"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Database
	db, err := store.New(ctx, cfg.DB.DSN())
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.RunMigrations(ctx); err != nil {
		logger.Error("migration failed", "error", err)
		os.Exit(1)
	}
	logger.Info("database ready")

	// Policy engine
	engine := policy.NewEngine()

	// Kafka consumer (runs in background)
	consumer := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.Topics,
		cfg.Kafka.ConsumerGroup,
		db,
		logger,
	)
	go consumer.Run(ctx)
	logger.Info("kafka consumer started", "brokers", cfg.Kafka.Brokers, "topics", cfg.Kafka.Topics)

	// HTTP router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Handlers
	policyH := handler.NewPolicyHandler(engine, logger)
	auditH := handler.NewAuditHandler(db, logger)

	// Routes
	r.Get("/healthz", handler.Healthz)
	r.Get("/readyz", handler.Readyz)
	r.Handle("/metrics", promhttp.Handler())

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/policy/evaluate", policyH.Evaluate)
		r.Get("/policy/rules", policyH.ListRules)

		r.Get("/audit/decisions", auditH.ListDecisions)
		r.Get("/audit/decisions/{agentID}", auditH.GetAgentDecisions)
		r.Get("/audit/trail", auditH.AuditTrail)

		r.Get("/agents", auditH.ListAgents)
	})

	// Server
	srv := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("api server starting", "addr", cfg.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			cancel()
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
}
