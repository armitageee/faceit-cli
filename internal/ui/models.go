package ui

import (
	"github.com/armitageee/faceit-cli/internal/config"
	"github.com/armitageee/faceit-cli/internal/entity"
	"github.com/armitageee/faceit-cli/internal/logger"
	"github.com/armitageee/faceit-cli/internal/repository"

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
func InitialModel(repo repository.FaceitRepository, config *config.Config, appLogger *logger.Logger) AppModel {
	model := AppModel{
		state:          StateSearch,
		repo:           repo,
		config:         config,
		logger:         appLogger,
		searchInput:    "",
		recentPlayers:  make([]string, 0),
		loading:        false,
		currentPage:    1,
		totalMatches:   0,
		matchesPerPage: config.MatchesPerPage,
		hasMoreMatches: false,
	}

	// If default player is configured, load it automatically
	if config.DefaultPlayer != "" {
		appLogger.Info("Loading default player", map[string]interface{}{
			"player": config.DefaultPlayer,
		})
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
		case StateMatchSearch:
			return m.updateMatchSearch(msg)
		case StateMatchStats:
			return m.updateMatchStats(msg)
		case StatePlayerMatchDetail:
			return m.updatePlayerMatchDetail(msg)
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
		m.currentPage = 1
		m.totalMatches = len(msg.matches)
		// Calculate if there are more pages
		totalPages := (len(msg.matches) + m.matchesPerPage - 1) / m.matchesPerPage
		m.hasMoreMatches = m.currentPage < totalPages
		m.state = StateMatches
		
		// Start background loading if we loaded less than the maximum
		if len(msg.matches) < m.config.MaxMatchesToLoad {
			m.backgroundLoading = true
			// Use parallel loading for better performance
			return m, m.loadBackgroundMatchesParallel()
		}
		return m, nil

	case matchesPageLoadedMsg:
		m.loading = false
		// Replace matches with the new page data
		m.matches = msg.matches
		m.currentPage = msg.page
		m.hasMoreMatches = msg.hasMore
		m.totalMatches = msg.totalMatches
		// Reset selected index to first match of the page
		m.selectedMatchIndex = 0
		return m, nil

	case backgroundMatchesLoadedMsg:
		// Update matches if we have more data from background loading
		if len(msg.matches) > len(m.matches) {
			m.matches = msg.matches
			// Recalculate pagination info
			m.totalMatches = len(m.matches)
			totalPages := (len(m.matches) + m.matchesPerPage - 1) / m.matchesPerPage
			m.hasMoreMatches = m.currentPage < totalPages
		}
		
		// Check if we need to continue loading
		if len(m.matches) < m.config.MaxMatchesToLoad && len(msg.matches) > 0 {
			// Continue loading with smaller batches
			return m, m.loadBackgroundMatches()
		} else {
			// Background loading is complete
			m.backgroundLoading = false
		}
		return m, nil

	case matchStatsLoadedMsg:
		m.loading = false
		m.matchStats = msg.matchStats
		m.state = StateMatchStats
		return m, nil

	case playerMatchStatsLoadedMsg:
		m.loading = false
		m.playerMatchStats = msg.matchStats
		m.state = StatePlayerMatchDetail
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

	case progressUpdateMsg:
		m.progress = msg.progress
		m.progressMessage = msg.message
		m.progressType = msg.progressType
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
		return m.renderLoadingScreen()
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
	case StateMatchSearch:
		return m.viewMatchSearch()
	case StateMatchStats:
		return m.viewMatchStats()
	case StatePlayerMatchDetail:
		return m.viewPlayerMatchDetail()
	case StatePlayerSwitch:
		return m.viewPlayerSwitch()
	case StateComparisonInput:
		return m.viewComparisonInput()
	case StateComparison:
		return m.viewComparison()
	case StateLoading:
		return m.renderLoadingScreen()
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
			Foreground(lipgloss.Color("#FAFAFA")).
			Padding(1, 0)

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


	winStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	lossStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)
)
