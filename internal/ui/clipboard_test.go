package ui

import (
	"os"
	"runtime"
	"testing"
)

func TestGetClipboardContent(t *testing.T) {
	// This test is limited because clipboard functionality depends on system tools
	// We'll test the function doesn't panic and handles errors gracefully
	
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "clipboard access",
			wantErr: false, // May or may not have content, but shouldn't error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := GetClipboardContent()
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetClipboardContent() expected error but got none")
				}
				return
			}
			
			// On some systems, clipboard might be empty or inaccessible
			// We just want to ensure the function doesn't panic
			if err != nil {
				t.Logf("GetClipboardContent() returned error (expected on some systems): %v", err)
			}
			
			// Content might be empty, that's okay
			t.Logf("Clipboard content: %q", content)
		})
	}
}

func TestGetClipboardContent_PlatformSpecific(t *testing.T) {
	// Test platform-specific behavior
	switch runtime.GOOS {
	case "darwin":
		t.Run("macOS", func(t *testing.T) {
			content, err := GetClipboardContent()
			if err != nil {
				t.Logf("macOS clipboard access failed (might be expected in CI): %v", err)
			} else {
				t.Logf("macOS clipboard content: %q", content)
			}
		})
	case "linux":
		t.Run("Linux", func(t *testing.T) {
			content, err := GetClipboardContent()
			if err != nil {
				t.Logf("Linux clipboard access failed (might be expected in CI): %v", err)
			} else {
				t.Logf("Linux clipboard content: %q", content)
			}
		})
	case "windows":
		t.Run("Windows", func(t *testing.T) {
			content, err := GetClipboardContent()
			if err != nil {
				t.Logf("Windows clipboard access failed (might be expected in CI): %v", err)
			} else {
				t.Logf("Windows clipboard content: %q", content)
			}
		})
	default:
		t.Run("Unknown OS", func(t *testing.T) {
			content, err := GetClipboardContent()
			if err == nil {
				t.Errorf("GetClipboardContent() should return error on unknown OS, got content: %q", content)
			}
		})
	}
}

func TestGetClipboardContent_Environment(t *testing.T) {
	// Test in different environments
	originalDisplay := os.Getenv("DISPLAY")
	
	t.Run("No DISPLAY", func(t *testing.T) {
		os.Unsetenv("DISPLAY")
		content, err := GetClipboardContent()
		if err != nil {
			t.Logf("Clipboard access without DISPLAY failed (expected on Linux): %v", err)
		} else {
			t.Logf("Clipboard content without DISPLAY: %q", content)
		}
	})
	
	// Restore original DISPLAY
	if originalDisplay != "" {
		os.Setenv("DISPLAY", originalDisplay)
	}
}

// Benchmark tests
func BenchmarkGetClipboardContent(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetClipboardContent()
	}
}

func BenchmarkGetClipboardContent_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = GetClipboardContent()
		}
	})
}
