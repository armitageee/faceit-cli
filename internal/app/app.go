package app

import (
	"context"
	"fmt"
	"time"

	"github.com/armitageee/faceit-cli/internal/cache"
	"github.com/armitageee/faceit-cli/internal/config"
	"github.com/armitageee/faceit-cli/internal/logger"
	"github.com/armitageee/faceit-cli/internal/repository"
	"github.com/armitageee/faceit-cli/internal/telemetry"
	"github.com/armitageee/faceit-cli/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// App represents the main application
type App struct {
	config     *config.Config
	repo       repository.FaceitRepository
	logger     *logger.Logger
	telemetry  *telemetry.Telemetry
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config, appLogger *logger.Logger, telemetryInstance *telemetry.Telemetry) *App {
	// Initialize repository with telemetry support
	var repo repository.FaceitRepository = repository.NewFaceitRepository(cfg.FaceitAPIKey, telemetryInstance)
	
	if cfg.CacheEnabled {
		appLogger.Info("Cache enabled", map[string]interface{}{
			"ttl_minutes": cfg.CacheTTL,
		})
		repo = cache.NewCachedFaceitRepository(repo, time.Duration(cfg.CacheTTL)*time.Minute)
	}
	
	return &App{
		config:    cfg,
		repo:      repo,
		logger:    appLogger,
		telemetry: telemetryInstance,
	}
}

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	if a.telemetry != nil {
		return a.telemetry.WithSpan(ctx, "app.run", func(ctx context.Context) error {
			return a.runInternal(ctx)
		})
	}
	return a.runInternal(ctx)
}

// runInternal contains the actual run logic
func (a *App) runInternal(ctx context.Context) error {
	a.logger.Info("Initializing UI model")
	
	// Create a span for UI initialization if telemetry is enabled
	if a.telemetry != nil {
		ctx, span := a.telemetry.StartSpan(ctx, "app.init_ui")
		span.SetAttributes(
			attribute.Bool("cache_enabled", a.config.CacheEnabled),
			attribute.Int("cache_ttl_minutes", a.config.CacheTTL),
			attribute.Int("matches_per_page", a.config.MatchesPerPage),
		)
		
		model := ui.InitialModel(a.repo, a.config, a.logger)
		span.End()
		
		a.logger.Info("Starting TUI program")
		
		// Create a span for TUI execution
		_, tuiSpan := a.telemetry.StartSpan(ctx, "app.tui_execution")
		defer tuiSpan.End()
		
		p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
		
		if _, err := p.Run(); err != nil {
			tuiSpan.RecordError(err)
			tuiSpan.SetStatus(codes.Error, err.Error())
			a.logger.Error("TUI program failed", map[string]interface{}{
				"error": err.Error(),
			})
			return fmt.Errorf("failed to run application: %w", err)
		}
		
		tuiSpan.SetStatus(codes.Ok, "TUI completed successfully")
		a.logger.Info("TUI program completed successfully")
		return nil
	}
	
	// No telemetry - run without tracing
	model := ui.InitialModel(a.repo, a.config, a.logger)
	a.logger.Info("Starting TUI program")
	
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	
	if _, err := p.Run(); err != nil {
		a.logger.Error("TUI program failed", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to run application: %w", err)
	}
	
	a.logger.Info("TUI program completed successfully")
	return nil
}
