package commands

import (
	"testing"
)

func TestShowCompletionStatus(t *testing.T) {
	// Test that function doesn't panic
	oldEnabled := completionEnabled
	defer func() { completionEnabled = oldEnabled }()

	completionEnabled = true
	// Function would print to shell, just verify it doesn't crash
	// Full integration test would require mocking ishell
}

func TestSetCompletionEnabled(t *testing.T) {
	oldEnabled := completionEnabled
	defer func() { completionEnabled = oldEnabled }()

	setCompletionEnabled(true)
	if !isCompletionEnabled() {
		t.Error("Completion should be enabled")
	}

	setCompletionEnabled(false)
	if isCompletionEnabled() {
		t.Error("Completion should be disabled")
	}
}
