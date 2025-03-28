package watchlist

import (
	"fmt"
	"sort"
	"strings"

	"stockterm/internal/config"
)

// Service provides operations for managing the watchlist
type Service struct {
	config *config.Config
}

// NewService creates a new watchlist service
func NewService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

// GetWatchlist returns the current watchlist
func (s *Service) GetWatchlist() ([]string, error) {
	return s.config.LoadWatchlist()
}

// AddTicker adds a ticker to the watchlist
func (s *Service) AddTicker(ticker string) error {
	// Normalize the ticker
	ticker = strings.TrimSpace(strings.ToUpper(ticker))
	if ticker == "" {
		return fmt.Errorf("ticker cannot be empty")
	}

	// Get the current watchlist
	watchlist, err := s.config.LoadWatchlist()
	if err != nil {
		return fmt.Errorf("failed to load watchlist: %w", err)
	}

	// Check if the ticker is already in the watchlist
	for _, t := range watchlist {
		if t == ticker {
			return fmt.Errorf("ticker %s is already in the watchlist", ticker)
		}
	}

	// Add the ticker to the watchlist
	watchlist = append(watchlist, ticker)

	// Sort the watchlist
	sort.Strings(watchlist)

	// Save the updated watchlist
	if err := s.config.SaveWatchlist(watchlist); err != nil {
		return fmt.Errorf("failed to save watchlist: %w", err)
	}

	return nil
}

// RemoveTicker removes a ticker from the watchlist
func (s *Service) RemoveTicker(ticker string) error {
	// Normalize the ticker
	ticker = strings.TrimSpace(strings.ToUpper(ticker))
	if ticker == "" {
		return fmt.Errorf("ticker cannot be empty")
	}

	// Get the current watchlist
	watchlist, err := s.config.LoadWatchlist()
	if err != nil {
		return fmt.Errorf("failed to load watchlist: %w", err)
	}

	// Check if the ticker is in the watchlist
	found := false
	var newWatchlist []string
	for _, t := range watchlist {
		if t != ticker {
			newWatchlist = append(newWatchlist, t)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("ticker %s is not in the watchlist", ticker)
	}

	// Save the updated watchlist
	if err := s.config.SaveWatchlist(newWatchlist); err != nil {
		return fmt.Errorf("failed to save watchlist: %w", err)
	}

	return nil
}

// AddMultipleTickers adds multiple tickers to the watchlist
func (s *Service) AddMultipleTickers(tickers []string) error {
	for _, ticker := range tickers {
		// Ignore errors for individual tickers
		_ = s.AddTicker(ticker)
	}
	return nil
}

// RemoveMultipleTickers removes multiple tickers from the watchlist
func (s *Service) RemoveMultipleTickers(tickers []string) error {
	for _, ticker := range tickers {
		// Ignore errors for individual tickers
		_ = s.RemoveTicker(ticker)
	}
	return nil
}

// UpdateWatchlist replaces the entire watchlist with a new one
func (s *Service) UpdateWatchlist(tickers []string) error {
	// Normalize and sort the tickers
	var normalizedTickers []string
	for _, ticker := range tickers {
		normalized := strings.TrimSpace(strings.ToUpper(ticker))
		if normalized != "" {
			normalizedTickers = append(normalizedTickers, normalized)
		}
	}
	sort.Strings(normalizedTickers)

	// Save the updated watchlist
	if err := s.config.SaveWatchlist(normalizedTickers); err != nil {
		return fmt.Errorf("failed to save watchlist: %w", err)
	}

	return nil
}
