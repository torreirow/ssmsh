package commands

import (
	"os"
	"os/exec"

	"github.com/abiosoft/ishell"
	"github.com/bwhaley/ssmsh/config"
)

const configUsage string = `
config usage: config <show|edit>

Manage ssmsh configuration.

  config show    Display config directory and file paths
  config edit    Open config file in $EDITOR (default: vi)
`

func configCmd(c *ishell.Context) {
	if len(c.Args) == 0 {
		shell.Println("Usage: config <show|edit>")
		return
	}

	switch c.Args[0] {
	case "show":
		configShow(c)
	case "edit":
		configEdit(c)
	default:
		shell.Printf("Unknown subcommand: %s\n", c.Args[0])
		shell.Println("Usage: config <show|edit>")
	}
}

func configShow(c *ishell.Context) {
	paths, err := config.GetConfigPaths()
	if err != nil {
		shell.Printf("Error: %v\n", err)
		return
	}

	shell.Println("Configuration:")
	shell.Printf("  Directory: %s\n", paths.ConfigDir)
	shell.Printf("  Config:    %s\n", paths.ConfigFile)
	shell.Printf("  History:   %s\n", paths.HistoryFile)
	shell.Printf("  Cache:     %s\n", paths.CacheFile)

	// Show if files exist
	checkFile := func(path string) string {
		if _, err := os.Stat(path); err == nil {
			return "✓"
		}
		return "✗"
	}

	shell.Printf("\nFiles:\n")
	shell.Printf("  [%s] config\n", checkFile(paths.ConfigFile))
	shell.Printf("  [%s] history\n", checkFile(paths.HistoryFile))
	shell.Printf("  [%s] cache\n", checkFile(paths.CacheFile))
}

func configEdit(c *ishell.Context) {
	paths, err := config.GetConfigPaths()
	if err != nil {
		shell.Printf("Error: %v\n", err)
		return
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi" // Default fallback
	}

	// Create default config if doesn't exist
	if _, err := os.Stat(paths.ConfigFile); os.IsNotExist(err) {
		if err := config.GenerateDefaultConfig(); err != nil {
			shell.Printf("Error creating config: %v\n", err)
			return
		}
	}

	// Open in editor
	cmd := exec.Command(editor, paths.ConfigFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		shell.Printf("Error opening editor: %v\n", err)
		return
	}

	shell.Println("Config updated. Restart ssmsh to apply changes.")
}
