package commands

import (
	"testing"
)

func TestFilterAndLimit(t *testing.T) {
	items := []string{
		"/dev/app/",
		"/dev/app/url",
		"/dev/app/domain",
		"/dev/admin/",
		"/dev/database/",
	}

	tests := []struct {
		name        string
		prefix      string
		dirsOnly    bool
		maxItems    int
		expectCount int
	}{
		{
			name:        "filter by prefix",
			prefix:      "/dev/app",
			dirsOnly:    false,
			maxItems:    50,
			expectCount: 3, // /dev/app/, /dev/app/url, /dev/app/domain
		},
		{
			name:        "directories only",
			prefix:      "/dev/",
			dirsOnly:    true,
			maxItems:    50,
			expectCount: 3, // /dev/app/, /dev/admin/, /dev/database/
		},
		{
			name:        "apply max limit",
			prefix:      "/dev/",
			dirsOnly:    false,
			maxItems:    2,
			expectCount: 2, // First 2 matches
		},
		{
			name:        "no matches",
			prefix:      "/prod/",
			dirsOnly:    false,
			maxItems:    50,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Temporarily set max items
			oldMaxItems := completionMaxItems
			completionMaxItems = tt.maxItems
			defer func() { completionMaxItems = oldMaxItems }()

			result := filterAndLimit(items, tt.prefix, tt.dirsOnly)
			if len(result) != tt.expectCount {
				t.Errorf("Expected %d results, got %d", tt.expectCount, len(result))
			}
		})
	}
}

func TestContainsControlChars(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect bool
	}{
		{
			name:   "normal string",
			input:  "/dev/app/url",
			expect: false,
		},
		{
			name:   "with newline",
			input:  "/dev/app\n/url",
			expect: true,
		},
		{
			name:   "with tab",
			input:  "/dev/app\t/url",
			expect: true,
		},
		{
			name:   "with null",
			input:  "/dev/app\x00/url",
			expect: true,
		},
		{
			name:   "with delete",
			input:  "/dev/app\x7f/url",
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsControlChars(tt.input)
			if result != tt.expect {
				t.Errorf("Expected %v, got %v for input %q", tt.expect, result, tt.input)
			}
		})
	}
}

func TestWrapCompleter(t *testing.T) {
	// Mock completer that returns fixed values
	mockCompleter := func(args []string) []string {
		return []string{"val1", "val2"}
	}

	wrapped := wrapCompleter(mockCompleter)

	// Enable completion
	oldEnabled := completionEnabled
	completionEnabled = true
	defer func() { completionEnabled = oldEnabled }()

	// Should call real completer
	result := wrapped([]string{})
	if len(result) != 2 {
		t.Errorf("Expected 2 results, got %d", len(result))
	}

	// Disable completion
	completionEnabled = false

	// Should return empty
	result = wrapped([]string{})
	if len(result) != 0 {
		t.Errorf("Expected 0 results when disabled, got %d", len(result))
	}
}

func TestFilterLongNames(t *testing.T) {
	longName := "/dev/" + string(make([]byte, 250)) // >200 chars total
	items := []string{
		"/dev/short",
		longName,
	}

	result := filterAndLimit(items, "/dev/", false)

	// Should only include short name
	if len(result) != 1 {
		t.Errorf("Expected 1 result (long name filtered), got %d", len(result))
	}

	if result[0] != "/dev/short" {
		t.Errorf("Expected /dev/short, got %s", result[0])
	}
}
