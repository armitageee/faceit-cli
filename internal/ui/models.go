package ui

import (
	"faceit-cli/internal/config"
	"faceit-cli/internal/entity"
	"faceit-cli/internal/repository"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Custom message types for async operations
type profileLoadedMsg struct {
	profile entity.PlayerProfile
}

type matchesLoadedMsg struct {
	matches []entity.PlayerMatchSummary
}

type errorMsg struct {
	err string
}

// InitialModel creates the initial application model
func InitialModel(repo repository.FaceitRepository, config *config.Config) AppModel {
	model := AppModel{
		state:         StateSearch,
		repo:          repo,
		config:        config,
		searchInput:   "",
		recentPlayers: make([]string, 0),
		loading:       false,
	}

	// If default player is configured, load it automatically
	if config.DefaultPlayer != "" {
		model.searchInput = config.DefaultPlayer
		model.loading = true
		model.state = StateLoading
	}

	return model
}

// Init initializes the model
func (m AppModel) Init() tea.Cmd {
	// If we have a default player, load it
	if m.config.DefaultPlayer != "" {
		return m.loadPlayerProfile(m.config.DefaultPlayer)
	}
	return nil
}

// Update handles messages and updates the model
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case StateSearch:
			return m.updateSearch(msg)
		case StateProfile:
			return m.updateProfile(msg)
		case StateMatches:
			return m.updateMatches(msg)
		case StateStats:
			return m.updateStats(msg)
		case StateMatchDetail:
			return m.updateMatchDetail(msg)
		case StatePlayerSwitch:
			return m.updatePlayerSwitch(msg)
		case StateComparisonInput:
			return m.updateComparisonInput(msg)
		case StateComparison:
			return m.updateComparison(msg)
		case StateError:
			return m.updateError(msg)
		}

	case profileLoadedMsg:
		m.loading = false
		m.player = &msg.profile
		m.state = StateProfile
		// Add to recent players
		m.addToRecentPlayers(msg.profile.Nickname)
		// Load lifetime stats
		return m, m.loadLifetimeStats()

	case matchesLoadedMsg:
		m.loading = false
		m.matches = msg.matches
		m.selectedMatchIndex = 0
		m.state = StateMatches
		return m, nil

	case statsLoadedMsg:
		m.loading = false
		m.stats = &msg.stats
		m.state = StateStats
		return m, nil

	case matchDetailLoadedMsg:
		m.loading = false
		m.matchDetail = &msg.matchDetail
		m.state = StateMatchDetail
		return m, nil

	case comparisonLoadedMsg:
		m.loading = false
		m.comparison = &msg.comparison
		m.state = StateComparison
		return m, nil

	case lifetimeStatsLoadedMsg:
		m.loading = false
		m.lifetimeStats = msg.stats
		return m, nil

	case errorMsg:
		m.loading = false
		m.error = msg.err
		m.state = StateError
		return m, nil
	}

	return m, nil
}

// View renders the current state
func (m AppModel) View() string {
	if m.loading {
		return m.viewLoading()
	}

	switch m.state {
	case StateSearch:
		return m.viewSearch()
	case StateProfile:
		return m.viewProfile()
	case StateMatches:
		return m.viewMatches()
	case StateStats:
		return m.viewStats()
	case StateMatchDetail:
		return m.viewMatchDetail()
	case StatePlayerSwitch:
		return m.viewPlayerSwitch()
	case StateComparisonInput:
		return m.viewComparisonInput()
	case StateComparison:
		return m.viewComparison()
	case StateError:
		return m.viewError()
	default:
		return "Unknown state"
	}
}

// Additional styling constants that are used in views
var (
	searchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(1, 2)

	profileStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2)

	matchesStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2)

	statsStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)

	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	winStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	lossStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)
)
