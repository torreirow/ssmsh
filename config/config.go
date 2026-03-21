package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	gcfg "gopkg.in/gcfg.v1"
)

const (
	LegacyConfigFileName = ".ssmshrc"
	ConfigDirName        = "ssmsh"
	ConfigFileName       = "config"
	HistoryFileName      = "history"
)

// Paths contains all ssmsh config paths
type Paths struct {
	ConfigDir   string
	ConfigFile  string
	HistoryFile string
	CacheFile   string
	LogFile     string
}

// GetConfigPaths returns the paths for ssmsh configuration files
func GetConfigPaths() (*Paths, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	return &Paths{
		ConfigDir:   configDir,
		ConfigFile:  filepath.Join(configDir, ConfigFileName),
		HistoryFile: filepath.Join(configDir, HistoryFileName),
		CacheFile:   filepath.Join(configDir, "cache.db"),
		LogFile:     filepath.Join(configDir, "session.log"),
	}, nil
}

// getConfigDir returns the configuration directory path
func getConfigDir() (string, error) {
	// Check $XDG_CONFIG_HOME first
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, ConfigDirName), nil
	}

	// Fall back to ~/.config/ssmsh
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".config", ConfigDirName), nil
}

// Config holds the default shell configuration
type Config struct {
	Default struct {
		// AWS settings
		Decrypt   bool
		Key       string
		Profile   string
		Region    string
		Overwrite bool
		Type      string
		Output    string

		// Tab-completion settings
		Completion         bool
		CompletionMaxItems int `gcfg:"completion-max-items"`
		CompletionCacheTTL int `gcfg:"completion-cache-ttl"`

		// Persistent cache settings
		CacheEnabled     bool   `gcfg:"cache-enabled"`
		CacheLocation    string `gcfg:"cache-location"`
		CacheMaxSize     int    `gcfg:"cache-max-size"`
		CacheCompression bool   `gcfg:"cache-compression"`

		// History settings
		HistorySize int  `gcfg:"history-size"`
		HistoryFile bool `gcfg:"history-file"`

		// Logging settings
		EnableLogging bool   `gcfg:"enable-logging"`
		LogLevel      string `gcfg:"log-level"`
	}
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	paths, err := GetConfigPaths()
	if err != nil {
		return err
	}

	// Create directory with restricted permissions
	if err := os.MkdirAll(paths.ConfigDir, 0700); err != nil {
		return err
	}

	// Verify permissions (important on multi-user systems)
	info, err := os.Stat(paths.ConfigDir)
	if err != nil {
		return err
	}

	// Warn if permissions are too open
	if info.Mode().Perm() != 0700 {
		fmt.Fprintf(os.Stderr, "Warning: Config directory has insecure permissions: %o\n",
			info.Mode().Perm())
		fmt.Fprintf(os.Stderr, "Fix with: chmod 700 %s\n", paths.ConfigDir)
	}

	return nil
}

// AutoMigrateConfig automatically migrates legacy config if needed
func AutoMigrateConfig() error {
	paths, err := GetConfigPaths()
	if err != nil {
		return err
	}

	// If new config exists, no migration needed
	if _, err := os.Stat(paths.ConfigFile); err == nil {
		return nil
	}

	// Check for legacy config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	legacyPath := filepath.Join(homeDir, LegacyConfigFileName)
	if _, err := os.Stat(legacyPath); os.IsNotExist(err) {
		// No legacy config - fresh install
		return nil
	}

	// AUTO-MIGRATE
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("  ssmsh Configuration Migration")
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Printf("\nMigrating configuration to new location...\n\n")

	// Ensure config directory exists
	if err := EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Read legacy config
	data, err := os.ReadFile(legacyPath)
	if err != nil {
		return fmt.Errorf("failed to read legacy config: %w", err)
	}

	// Write to new location with restricted permissions
	if err := os.WriteFile(paths.ConfigFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write new config: %w", err)
	}

	fmt.Printf("  ✓ Config migrated to: %s\n", paths.ConfigFile)

	// Rename (not delete) legacy file for safety
	backupPath := legacyPath + ".backup"
	if err := os.Rename(legacyPath, backupPath); err != nil {
		// Non-fatal - just warn
		fmt.Printf("  ⚠ Could not rename old config: %v\n", err)
		fmt.Printf("    Please manually remove: %s\n", legacyPath)
	} else {
		fmt.Printf("  ✓ Legacy config backed up to: %s\n", backupPath)
	}

	fmt.Println("\nMigration complete!")
	fmt.Printf("You can safely delete the backup after verifying:\n")
	fmt.Printf("  rm %s\n\n", backupPath)
	fmt.Println("═══════════════════════════════════════════════")

	// Brief pause so user can read the message
	time.Sleep(2 * time.Second)

	return nil
}

// ReadConfig reads ssmsh configuration with priority chain
func ReadConfig(cfgFile string) (Config, error) {
	// If explicit config file provided via flag
	if cfgFile != "" {
		return readConfigFromFile(cfgFile)
	}

	// Check $SSMSH_CONFIG environment variable
	if envConfig := os.Getenv("SSMSH_CONFIG"); envConfig != "" {
		return readConfigFromFile(envConfig)
	}

	// Try XDG config path
	paths, err := GetConfigPaths()
	if err != nil {
		return Config{}, err
	}

	if _, err := os.Stat(paths.ConfigFile); err == nil {
		return readConfigFromFile(paths.ConfigFile)
	}

	// Try legacy ~/.ssmshrc (backwards compatibility)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	legacyPath := filepath.Join(homeDir, LegacyConfigFileName)
	if _, err := os.Stat(legacyPath); err == nil {
		// Found legacy config - warn user
		fmt.Fprintf(os.Stderr, "Warning: Using legacy config at %s\n", legacyPath)
		fmt.Fprintf(os.Stderr, "Consider migrating to: %s\n", paths.ConfigFile)
		fmt.Fprintf(os.Stderr, "Migration will happen automatically on next restart.\n\n")

		return readConfigFromFile(legacyPath)
	}

	// No config found - use defaults
	return Config{}, nil
}

// readConfigFromFile reads and parses a config file
func readConfigFromFile(path string) (Config, error) {
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, path)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file %s: %w", path, err)
	}
	return cfg, nil
}

// ExpandPath expands ~ to user home directory in paths
func ExpandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~/") {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, path[2:]), nil
}

// ValidateAndApplyDefaults validates config values and applies defaults
func ValidateAndApplyDefaults(cfg *Config) {
	// Enable completion by default if not explicitly configured
	// Note: gcfg doesn't distinguish between "not set" and "set to false"
	// So we check if ANY completion settings are configured. If none are, enable by default.
	hasCompletionConfig := cfg.Default.Completion || cfg.Default.CompletionMaxItems > 0 || cfg.Default.CompletionCacheTTL > 0
	if !hasCompletionConfig {
		cfg.Default.Completion = true
	}

	// Validate and set completion settings
	if cfg.Default.CompletionMaxItems <= 0 {
		cfg.Default.CompletionMaxItems = 50
	} else if cfg.Default.CompletionMaxItems > 500 {
		fmt.Fprintf(os.Stderr, "Warning: completion-max-items capped at 500 (was %d)\n",
			cfg.Default.CompletionMaxItems)
		cfg.Default.CompletionMaxItems = 500
	}

	// Validate TTL
	if cfg.Default.CompletionCacheTTL < 0 {
		fmt.Fprintf(os.Stderr, "Warning: completion-cache-ttl must be >= 0, using default (30)\n")
		cfg.Default.CompletionCacheTTL = 30
	} else if cfg.Default.CompletionCacheTTL > 3600 {
		fmt.Fprintf(os.Stderr, "Warning: completion-cache-ttl capped at 3600 seconds (1 hour)\n")
		cfg.Default.CompletionCacheTTL = 3600
	} else if cfg.Default.CompletionCacheTTL == 0 {
		cfg.Default.CompletionCacheTTL = 30 // Default
	}

	// Validate cache max size
	if cfg.Default.CacheMaxSize < 0 {
		cfg.Default.CacheMaxSize = 50
	} else if cfg.Default.CacheMaxSize > 500 {
		fmt.Fprintf(os.Stderr, "Warning: cache-max-size capped at 500 MB\n")
		cfg.Default.CacheMaxSize = 500
	} else if cfg.Default.CacheMaxSize == 0 {
		cfg.Default.CacheMaxSize = 50 // Default
	}

	// Apply other defaults if not set
	if cfg.Default.HistorySize == 0 {
		cfg.Default.HistorySize = 1000
	}

	// Expand tilde in cache location
	if cfg.Default.CacheLocation != "" {
		expanded, err := ExpandPath(cfg.Default.CacheLocation)
		if err == nil {
			cfg.Default.CacheLocation = expanded
		}
	}
}

// GenerateDefaultConfig creates a default config file with comments
func GenerateDefaultConfig() error {
	paths, err := GetConfigPaths()
	if err != nil {
		return err
	}

	// Check if config already exists
	if _, err := os.Stat(paths.ConfigFile); err == nil {
		return fmt.Errorf("config already exists at %s", paths.ConfigFile)
	}

	// Ensure directory exists
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	// Generate config with helpful comments
	configContent := `# ssmsh configuration file
# Location: ~/.config/ssmsh/config
#
# This file uses INI-like syntax. For more information:
# https://github.com/bwhaley/ssmsh

[default]
# AWS Configuration
# region=us-east-1
# profile=default
# key=your-kms-key-id

# Parameter Store Defaults
type=SecureString
overwrite=false
decrypt=true
output=json

# Tab Completion
completion=true
completion-max-items=50
completion-cache-ttl=30

# Persistent Cache
cache-enabled=true
cache-max-size=50
cache-compression=true
# cache-location=~/.config/ssmsh/cache.db

# Command History
history-size=1000
history-file=true

# Logging (for debugging)
# enable-logging=false
# log-level=info
`

	if err := os.WriteFile(paths.ConfigFile, []byte(configContent), 0600); err != nil {
		return err
	}

	fmt.Printf("Default config created at: %s\n", paths.ConfigFile)
	return nil
}
