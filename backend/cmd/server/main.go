package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"remember_my_pill/backend/internal/access"
	"remember_my_pill/backend/internal/config"
	"remember_my_pill/backend/internal/httpapi"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		logger.Error("startup configuration failed")
		os.Exit(1)
	}
	emailSender, err := access.NewEmailSender(access.ProviderConfig{Provider: cfg.EmailProvider, FromAddress: cfg.EmailFromAddress, FromName: cfg.EmailFromName, MessageStream: cfg.PostmarkMessageStream, PostmarkServerToken: cfg.PostmarkServerToken, PostmarkTemplateAlias: cfg.PostmarkStatusAccessTemplate})
	if err != nil {
		logger.Error("email provider configuration failed")
		os.Exit(1)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Error("database configuration failed")
		os.Exit(1)
	}
	defer pool.Close()
	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("database unavailable at startup")
		os.Exit(1)
	}

	server := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: httpapi.NewRouter(pool, pool, httpapi.Options{
			AllowedOrigin: cfg.AllowedOrigin, ConsentVersion: cfg.ConsentVersion, PilotLegalContentApproved: cfg.PilotLegalContentApproved, TrustedProxies: cfg.TrustedProxies, Logger: logger,
			StatusAccessEnabled: cfg.StatusAccessEnabled, StatusAccessRateLimitKey: cfg.StatusAccessRateLimitKey, StatusStore: pool, StatusEmailSender: emailSender,
			StatusAccessBaseURL: cfg.StatusAccessBaseURL, StatusSessionCookieSecure: cfg.StatusSessionCookieSecure,
		}),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	go func() {
		logger.Info("server started", "port", cfg.Port, "consent_configured", cfg.ConsentVersion != "", "pilot_legal_content_approved", cfg.PilotLegalContentApproved, "trusted_proxy_count", len(cfg.TrustedProxies), "email_provider", cfg.EmailProvider)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped unexpectedly")
			os.Exit(1)
		}
	}()

	<-shutdown
	logger.Info("server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed")
	}
}
