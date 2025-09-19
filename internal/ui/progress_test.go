package ui

import (
	"testing"

	"github.com/armitageee/faceit-cli/internal/config"
	"github.com/armitageee/faceit-cli/internal/logger"
)

func TestProgressBar(t *testing.T) {
	// Create a test model
	config := &config.Config{
		FaceitAPIKey: "test-api-key",
		CacheEnabled: false,
	}
	loggerConfig := logger.Config{
		Level:          logger.LogLevelInfo,
		KafkaEnabled:   false,
		ServiceName:    "test-service",
		ProductionMode: false,
		LogToStdout:    false,
	}
	appLogger, _ := logger.New(loggerConfig)
	
	model := AppModel{
		state:           StateLoading,
		loading:         true,
		progress:        0.5,
		progressMessage: "Loading test data...",
		progressType:    "matches",
		width:           80,
		height:          24,
		config:          config,
		logger:          appLogger,
	}

	// Test progress bar rendering
	progressBar := model.renderProgressBar()
	if progressBar == "" {
		t.Error("Progress bar should not be empty")
	}

	// Test with different progress values
	testCases := []struct {
		progress     float64
		message      string
		progressType string
	}{
		{0.0, "Starting...", "matches"},
		{0.25, "Loading data...", "stats"},
		{0.5, "Processing...", "match_stats"},
		{0.75, "Almost done...", "matches"},
		{1.0, "Complete!", "matches"},
	}

	for _, tc := range testCases {
		model.progress = tc.progress
		model.progressMessage = tc.message
		model.progressType = tc.progressType

		progressBar := model.renderProgressBar()
		if progressBar == "" {
			t.Errorf("Progress bar should not be empty for progress %f", tc.progress)
		}

		// Check that the message is included
		if !containsProgress(progressBar, tc.message) {
			t.Errorf("Progress bar should contain message '%s'", tc.message)
		}
	}
}

func TestProgressBarEdgeCases(t *testing.T) {
	config := &config.Config{
		FaceitAPIKey: "test-api-key",
		CacheEnabled: false,
	}
	loggerConfig := logger.Config{
		Level:          logger.LogLevelInfo,
		KafkaEnabled:   false,
		ServiceName:    "test-service",
		ProductionMode: false,
		LogToStdout:    false,
	}
	appLogger, _ := logger.New(loggerConfig)
	
	model := AppModel{
		state:           StateLoading,
		loading:         true,
		progress:        0.0,
		progressMessage: "",
		progressType:    "",
		width:           80,
		height:          24,
		config:          config,
		logger:          appLogger,
	}

	// Test with empty message
	progressBar := model.renderProgressBar()
	if progressBar != "" {
		t.Error("Progress bar should be empty when message is empty")
	}

	// Test with negative progress
	model.progress = -0.1
	model.progressMessage = "Test"
	progressBar = model.renderProgressBar()
	if progressBar == "" {
		t.Error("Progress bar should not be empty with negative progress")
	}

	// Test with progress > 1.0
	model.progress = 1.5
	progressBar = model.renderProgressBar()
	if progressBar == "" {
		t.Error("Progress bar should not be empty with progress > 1.0")
	}
}

func TestProgressBarResponsive(t *testing.T) {
	config := &config.Config{
		FaceitAPIKey: "test-api-key",
		CacheEnabled: false,
	}
	loggerConfig := logger.Config{
		Level:          logger.LogLevelInfo,
		KafkaEnabled:   false,
		ServiceName:    "test-service",
		ProductionMode: false,
		LogToStdout:    false,
	}
	appLogger, _ := logger.New(loggerConfig)
	
	model := AppModel{
		state:           StateLoading,
		loading:         true,
		progress:        0.5,
		progressMessage: "Test message",
		progressType:    "matches",
		config:          config,
		logger:          appLogger,
	}

	// Test different widths
	widths := []int{40, 60, 80, 120}
	for _, width := range widths {
		model.width = width
		model.height = 24

		progressBar := model.renderProgressBar()
		if progressBar == "" {
			t.Errorf("Progress bar should not be empty for width %d", width)
		}
	}
}

func TestProgressUpdate(t *testing.T) {
	config := &config.Config{
		FaceitAPIKey: "test-api-key",
		CacheEnabled: false,
	}
	loggerConfig := logger.Config{
		Level:          logger.LogLevelInfo,
		KafkaEnabled:   false,
		ServiceName:    "test-service",
		ProductionMode: false,
		LogToStdout:    false,
	}
	appLogger, _ := logger.New(loggerConfig)
	
	model := AppModel{
		state:           StateLoading,
		loading:         true,
		progress:        0.0,
		progressMessage: "",
		progressType:    "",
		width:           80,
		height:          24,
		config:          config,
		logger:          appLogger,
	}

	// Test progress update
	updatedModel := model.updateProgress(0.5, "Test message", "matches")
	if updatedModel.progress != 0.5 {
		t.Errorf("Expected progress 0.5, got %f", updatedModel.progress)
	}
	if updatedModel.progressMessage != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", updatedModel.progressMessage)
	}
	if updatedModel.progressType != "matches" {
		t.Errorf("Expected type 'matches', got '%s'", updatedModel.progressType)
	}

	// Test progress reset
	resetModel := model.resetProgress()
	if resetModel.progress != 0.0 {
		t.Errorf("Expected progress 0.0, got %f", resetModel.progress)
	}
	if resetModel.progressMessage != "" {
		t.Errorf("Expected empty message, got '%s'", resetModel.progressMessage)
	}
	if resetModel.progressType != "" {
		t.Errorf("Expected empty type, got '%s'", resetModel.progressType)
	}
}

func TestProgressUpdateMsg(t *testing.T) {
	// Test progress update message
	msg := progressUpdateMsg{
		progress:     0.75,
		message:      "Test message",
		progressType: "stats",
	}

	if msg.progress != 0.75 {
		t.Errorf("Expected progress 0.75, got %f", msg.progress)
	}
	if msg.message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", msg.message)
	}
	if msg.progressType != "stats" {
		t.Errorf("Expected type 'stats', got '%s'", msg.progressType)
	}
}

func TestSimulateProgress(t *testing.T) {
	config := &config.Config{
		FaceitAPIKey: "test-api-key",
		CacheEnabled: false,
	}
	loggerConfig := logger.Config{
		Level:          logger.LogLevelInfo,
		KafkaEnabled:   false,
		ServiceName:    "test-service",
		ProductionMode: false,
		LogToStdout:    false,
	}
	appLogger, _ := logger.New(loggerConfig)
	
	model := AppModel{
		state:           StateLoading,
		loading:         true,
		progress:        0.0,
		progressMessage: "Test message",
		progressType:    "matches",
		config:          config,
		logger:          appLogger,
	}

	// Test simulate progress command
	cmd := model.simulateProgress()
	if cmd == nil {
		t.Error("Simulate progress command should not be nil")
	}

	// Test the command execution
	msg := cmd()
	progressMsg, ok := msg.(progressUpdateMsg)
	if !ok {
		t.Error("Expected progressUpdateMsg")
	}

	if progressMsg.progress <= 0.0 || progressMsg.progress > 0.9 {
		t.Errorf("Expected progress between 0 and 0.9, got %f", progressMsg.progress)
	}
	if progressMsg.message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", progressMsg.message)
	}
	if progressMsg.progressType != "matches" {
		t.Errorf("Expected type 'matches', got '%s'", progressMsg.progressType)
	}
}

func TestUpdateProgressCmd(t *testing.T) {
	// Test update progress command
	cmd := updateProgressCmd(0.5, "Test message", "matches")
	if cmd == nil {
		t.Error("Update progress command should not be nil")
	}

	// Test the command execution
	msg := cmd()
	progressMsg, ok := msg.(progressUpdateMsg)
	if !ok {
		t.Error("Expected progressUpdateMsg")
	}

	if progressMsg.progress != 0.5 {
		t.Errorf("Expected progress 0.5, got %f", progressMsg.progress)
	}
	if progressMsg.message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", progressMsg.message)
	}
	if progressMsg.progressType != "matches" {
		t.Errorf("Expected type 'matches', got '%s'", progressMsg.progressType)
	}
}

func TestLoadingScreen(t *testing.T) {
	config := &config.Config{
		FaceitAPIKey: "test-api-key",
		CacheEnabled: false,
	}
	loggerConfig := logger.Config{
		Level:          logger.LogLevelInfo,
		KafkaEnabled:   false,
		ServiceName:    "test-service",
		ProductionMode: false,
		LogToStdout:    false,
	}
	appLogger, _ := logger.New(loggerConfig)
	
	model := AppModel{
		state:           StateLoading,
		loading:         true,
		progress:        0.5,
		progressMessage: "Loading test data...",
		progressType:    "matches",
		width:           80,
		height:          24,
		config:          config,
		logger:          appLogger,
	}

	// Test loading screen rendering
	loadingScreen := model.renderLoadingScreen()
	if loadingScreen == "" {
		t.Error("Loading screen should not be empty")
	}

	// Test with different progress states
	model.progress = 0.0
	loadingScreen = model.renderLoadingScreen()
	if loadingScreen == "" {
		t.Error("Loading screen should not be empty with 0 progress")
	}

	model.progress = 1.0
	loadingScreen = model.renderLoadingScreen()
	if loadingScreen == "" {
		t.Error("Loading screen should not be empty with 100% progress")
	}
}

// Helper function to check if string contains substring
func containsProgress(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		   len(s) > len(substr) && containsProgress(s[1:], substr)
}

// Benchmark tests
func BenchmarkRenderProgressBar(b *testing.B) {
	config := &config.Config{
		FaceitAPIKey: "test-api-key",
		CacheEnabled: false,
	}
	loggerConfig := logger.Config{
		Level:          logger.LogLevelInfo,
		KafkaEnabled:   false,
		ServiceName:    "test-service",
		ProductionMode: false,
		LogToStdout:    false,
	}
	appLogger, _ := logger.New(loggerConfig)
	
	model := AppModel{
		state:           StateLoading,
		loading:         true,
		progress:        0.5,
		progressMessage: "Loading test data...",
		progressType:    "matches",
		width:           80,
		height:          24,
		config:          config,
		logger:          appLogger,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.renderProgressBar()
	}
}

func BenchmarkSimulateProgress(b *testing.B) {
	config := &config.Config{
		FaceitAPIKey: "test-api-key",
		CacheEnabled: false,
	}
	loggerConfig := logger.Config{
		Level:          logger.LogLevelInfo,
		KafkaEnabled:   false,
		ServiceName:    "test-service",
		ProductionMode: false,
		LogToStdout:    false,
	}
	appLogger, _ := logger.New(loggerConfig)
	
	model := AppModel{
		state:           StateLoading,
		loading:         true,
		progress:        0.0,
		progressMessage: "Test message",
		progressType:    "matches",
		config:          config,
		logger:          appLogger,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.simulateProgress()
	}
}
