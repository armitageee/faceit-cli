package app

import (
	"context"
	"fmt"

	"faceit-cli/internal/config"
	"faceit-cli/internal/repository"
	"faceit-cli/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// App represents the main application
type App struct {
	config *config.Config
	repo   repository.FaceitRepository
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config) *App {
	repo := repository.NewFaceitRepository(cfg.FaceitAPIKey)
	return &App{
		config: cfg,
		repo:   repo,
	}
}

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	model := ui.InitialModel(a.repo, a.config)
	
	p := tea.NewProgram(model, tea.WithAltScreen())
	
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run application: %w", err)
	}
	
	return nil
}
