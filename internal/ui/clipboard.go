package ui

import (
	"os/exec"
	"runtime"
	"strings"
)

// GetClipboardContent retrieves content from the system clipboard
func GetClipboardContent() (string, error) {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("pbpaste")
	case "linux":
		// Try xclip first, then xsel
		cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
		if _, err := exec.LookPath("xclip"); err != nil {
			cmd = exec.Command("xsel", "--clipboard", "--output")
		}
	case "windows":
		cmd = exec.Command("powershell", "-command", "Get-Clipboard")
	default:
		return "", nil // Unsupported platform
	}
	
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(string(output)), nil
}

// isPasteKey checks if the key combination is a paste command
func isPasteKey(key string) bool {
	return key == "ctrl+v" || key == "cmd+v"
}
