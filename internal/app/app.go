package app

import (
	"context"
	"fmt"
	"time"

	"faceit-cli/internal/cache"
	"faceit-cli/internal/config"
	"faceit-cli/internal/logger"
	"faceit-cli/internal/repository"
	"faceit-cli/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// App represents the main application
type App struct {
	config *config.Config
	repo   repository.FaceitRepository
	logger *logger.Logger
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config, appLogger *logger.Logger) *App {
	// Initialize repository with optional caching
	var repo repository.FaceitRepository = repository.NewFaceitRepository(cfg.FaceitAPIKey)
	
	if cfg.CacheEnabled {
		appLogger.Info("Cache enabled", map[string]interface{}{
			"ttl_minutes": cfg.CacheTTL,
		})
		repo = cache.NewCachedFaceitRepository(repo, time.Duration(cfg.CacheTTL)*time.Minute)
	}
	
	return &App{
		config: cfg,
		repo:   repo,
		logger: appLogger,
	}
}

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	a.logger.Info("Initializing UI model")
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
