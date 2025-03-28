package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the application configuration
type Config struct {
	// ConfigPath is the path to the configuration file
	ConfigPath string
	// WatchlistPath is the path to the watchlist file
	WatchlistPath string
	// DefaultTimeRange is the default time range for stock data
	DefaultTimeRange string
	// DefaultCurrency is the default currency for stock data
	DefaultCurrency string
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fall back to current directory if home directory can't be determined
		homeDir = "."
	}

	// Create the config directory if it doesn't exist
	configDir := filepath.Join(homeDir, ".stockterm")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Printf("Warning: Failed to create config directory: %v\n", err)
		}
	}

	return &Config{
		ConfigPath:       filepath.Join(configDir, "config.yaml"),
		WatchlistPath:    filepath.Join(configDir, "watchlist.txt"),
		DefaultTimeRange: "1d",
		DefaultCurrency:  "USD",
	}
}

// LoadWatchlist loads the watchlist from the file
func (c *Config) LoadWatchlist() ([]string, error) {
	// Check if the watchlist file exists
	if _, err := os.Stat(c.WatchlistPath); os.IsNotExist(err) {
		// Create an empty watchlist file
		if err := os.WriteFile(c.WatchlistPath, []byte(""), 0644); err != nil {
			return nil, fmt.Errorf("failed to create watchlist file: %w", err)
		}
		return []string{}, nil
	}

	// Read the watchlist file
	content, err := os.ReadFile(c.WatchlistPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read watchlist file: %w", err)
	}

	// If the file is empty, return an empty slice
	if len(content) == 0 {
		return []string{}, nil
	}

	// Split the content into an array by splitting with ","
	watchlist := strings.Split(string(content), ",")

	// Remove leading and trailing whitespaces from each element
	for i, item := range watchlist {
		watchlist[i] = strings.TrimSpace(item)
	}

	return watchlist, nil
}

// SaveWatchlist saves the watchlist to the file
func (c *Config) SaveWatchlist(watchlist []string) error {
	// Join the watchlist into a comma-separated string
	content := strings.Join(watchlist, ",")

	// Create the directory if it doesn't exist
	dir := filepath.Dir(c.WatchlistPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Write the content to the file
	if err := os.WriteFile(c.WatchlistPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write watchlist file: %w", err)
	}

	return nil
}

// MigrateFromLegacy migrates the watchlist from the legacy location
func (c *Config) MigrateFromLegacy(legacyPath string) error {
	// Check if the legacy file exists
	if _, err := os.Stat(legacyPath); os.IsNotExist(err) {
		// Legacy file doesn't exist, nothing to migrate
		return nil
	}

	// Check if the new watchlist file already exists
	if _, err := os.Stat(c.WatchlistPath); err == nil {
		// New watchlist file already exists, don't overwrite it
		return nil
	}

	// Read the legacy file
	content, err := os.ReadFile(legacyPath)
	if err != nil {
		return fmt.Errorf("failed to read legacy watchlist file: %w", err)
	}

	// Create the directory if it doesn't exist
	dir := filepath.Dir(c.WatchlistPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Write the content to the new file
	if err := os.WriteFile(c.WatchlistPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write watchlist file: %w", err)
	}

	return nil
}
