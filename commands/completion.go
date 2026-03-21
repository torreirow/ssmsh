package commands

import (
	"strconv"
	"strings"

	"github.com/abiosoft/ishell"
)

const completionUsage string = `
completion usage: completion [true|false|stats|clear-cache|save-cache|reload-cache]

Controls tab-completion for parameter paths and names.

  completion              Show current state
  completion true         Enable tab-completion
  completion false        Disable tab-completion
  completion stats        Show cache statistics
  completion clear-cache  Clear all cached entries
  completion save-cache   Force save cache to disk
  completion reload-cache Reload cache from disk

Settings (configured in ~/.config/ssmsh/config):
  Max items:  Shows up to N suggestions
  Cache TTL:  Results cached for N seconds
`

const completionError = "value for completion must be boolean (true or false)"

func completionCmd(c *ishell.Context) {
	if len(c.Args) == 0 {
		showCompletionStatus()
		return
	}

	switch strings.ToLower(c.Args[0]) {
	case "true":
		setCompletionEnabled(true)
		showCompletionStatus()
	case "false":
		setCompletionEnabled(false)
		showCompletionStatus()
	case "stats":
		showCacheStats()
	case "clear-cache":
		clearAllCaches()
	case "save-cache":
		saveCache()
	case "reload-cache":
		reloadCache()
	default:
		// Try to parse as boolean for backwards compat
		if v, err := strconv.ParseBool(c.Args[0]); err == nil {
			setCompletionEnabled(v)
			showCompletionStatus()
		} else {
			shell.Println(completionError)
		}
	}
}

func showCompletionStatus() {
	state := "disabled"
	if isCompletionEnabled() {
		state = "enabled"
	}
	shell.Printf("Tab completion is %s\n", state)

	if isCompletionEnabled() {
		shell.Printf("  Max items: %d\n", getCompletionMaxItems())
		shell.Printf("  Cache TTL: %d seconds\n", getCompletionCacheTTL())
	}
}

func showCacheStats() {
	shell.Println("Cache Statistics:")
	shell.Println()

	// Memory cache
	memStats := memoryCache.stats()
	shell.Println("  In-Memory Cache:")
	shell.Printf("    Entries:      %d / %d (max)\n", memStats.Entries, maxCacheEntries)
	shell.Printf("    Total items:  %d\n", memStats.TotalItems)
	shell.Printf("    TTL:          %d seconds\n", getCompletionCacheTTL())
	shell.Println()

	// Persistent cache
	if persistentCache != nil {
		pStats := persistentCache.Stats()
		shell.Println("  Persistent Cache:")
		shell.Printf("    Enabled:      %v\n", pStats.Enabled)
		shell.Printf("    Location:     %s\n", pStats.Location)
		shell.Printf("    Entries:      %d\n", pStats.Entries)
		shell.Printf("    Total items:  %d\n", pStats.TotalItems)
		shell.Printf("    File size:    %.2f MB / %.0f MB (max)\n",
			pStats.FileSizeMB, pStats.MaxSizeMB)
		shell.Printf("    Compression:  %v\n", pStats.Compression)
	} else {
		shell.Println("  Persistent Cache: Disabled")
	}
}

func clearAllCaches() {
	// Clear memory cache
	clearCache()

	// Clear persistent cache
	if persistentCache != nil {
		if err := persistentCache.clear(); err != nil {
			shell.Printf("Error clearing persistent cache: %v\n", err)
		} else {
			shell.Println("Persistent cache cleared")
		}
	}

	shell.Println("All caches cleared")
}

func saveCache() {
	if persistentCache == nil {
		shell.Println("Persistent cache is disabled")
		return
	}

	if err := persistentCache.save(); err != nil {
		shell.Printf("Error saving cache: %v\n", err)
	} else {
		shell.Println("Cache saved to disk")
	}
}

func reloadCache() {
	if persistentCache == nil {
		shell.Println("Persistent cache is disabled")
		return
	}

	if err := persistentCache.load(); err != nil {
		shell.Printf("Error reloading cache: %v\n", err)
	} else {
		shell.Println("Cache reloaded from disk")
	}
}
