package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"stockterm/internal/api"
	"stockterm/internal/config"
	"stockterm/internal/ui"
	"stockterm/internal/watchlist"
)

func main() {
	// Create a context that can be cancelled on SIGINT or SIGTERM
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Initialize configuration
	cfg := config.DefaultConfig()

	// Migrate from legacy config if needed
	if err := cfg.MigrateFromLegacy("./ggs.config"); err != nil {
		fmt.Printf("Warning: Failed to migrate from legacy config: %v\n", err)
	}

	// Initialize services
	yahooClient := api.NewYahooFinanceClient()
	watchlistService := watchlist.NewService(cfg)
	tableRenderer := ui.NewTableRenderer()

	// Parse command-line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	// Execute the command
	if err := executeCommand(ctx, command, args, yahooClient, watchlistService, tableRenderer); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func executeCommand(
	ctx context.Context,
	command string,
	args []string,
	yahooClient *api.YahooFinanceClient,
	watchlistService *watchlist.Service,
	tableRenderer *ui.TableRenderer,
) error {
	switch command {
	case "get":
		if len(args) < 1 {
			return fmt.Errorf("missing ticker argument")
		}
		return getTickersPrice(ctx, args[0], yahooClient, tableRenderer)

	case "get-all":
		return getWatchlistPrice(ctx, yahooClient, watchlistService, tableRenderer)

	case "list":
		return displayWatchlist(watchlistService)

	case "add":
		if len(args) < 1 {
			return fmt.Errorf("missing ticker argument")
		}
		return addTickersToWatchlist(args[0], watchlistService)

	case "remove":
		if len(args) < 1 {
			return fmt.Errorf("missing ticker argument")
		}
		return removeTickersFromWatchlist(args[0], watchlistService)

	case "help":
		printUsage()
		return nil

	case "version":
		printVersion()
		return nil

	default:
		return fmt.Errorf("invalid command: '%s'\n\n%s", command, getUsageText())
	}
}

func getTickersPrice(ctx context.Context, tickersArg string, yahooClient *api.YahooFinanceClient, tableRenderer *ui.TableRenderer) error {
	// Split the tickers by comma
	tickers := strings.Split(tickersArg, ",")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Fetch the stock data
	responses, err := yahooClient.FetchMultipleStocks(ctx, tickers, "1d")
	if err != nil {
		return fmt.Errorf("error fetching stock data: %w", err)
	}

	// Render the table
	tableRenderer.RenderChartResponses(responses)

	return nil
}

func getWatchlistPrice(ctx context.Context, yahooClient *api.YahooFinanceClient, watchlistService *watchlist.Service, tableRenderer *ui.TableRenderer) error {
	// Get the watchlist
	watchlist, err := watchlistService.GetWatchlist()
	if err != nil {
		return fmt.Errorf("error getting watchlist: %w", err)
	}

	if len(watchlist) == 0 {
		fmt.Println("Watchlist is empty. Add tickers with 'stockterm add <ticker>'")
		return nil
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Fetch the stock data
	responses, err := yahooClient.FetchMultipleStocks(ctx, watchlist, "1d")
	if err != nil {
		return fmt.Errorf("error fetching stock data: %w", err)
	}

	// Render the table
	tableRenderer.RenderChartResponses(responses)

	return nil
}

func displayWatchlist(watchlistService *watchlist.Service) error {
	// Get the watchlist
	watchlist, err := watchlistService.GetWatchlist()
	if err != nil {
		return fmt.Errorf("error getting watchlist: %w", err)
	}

	if len(watchlist) == 0 {
		fmt.Println("Watchlist is empty. Add tickers with 'stockterm add <ticker>'")
		return nil
	}

	// Run the watchlist editor
	updatedWatchlist, saved := ui.RunWatchlistEditor(watchlist)
	if !saved {
		fmt.Println("Watchlist not updated")
		return nil
	}

	// Update the watchlist
	if err := watchlistService.UpdateWatchlist(updatedWatchlist); err != nil {
		return fmt.Errorf("error updating watchlist: %w", err)
	}

	fmt.Println("Watchlist has been updated!")
	return nil
}

func addTickersToWatchlist(tickersArg string, watchlistService *watchlist.Service) error {
	// Split the tickers by comma
	tickers := strings.Split(tickersArg, ",")

	// Add each ticker to the watchlist
	for _, ticker := range tickers {
		ticker = strings.TrimSpace(ticker)
		if ticker == "" {
			continue
		}

		if err := watchlistService.AddTicker(ticker); err != nil {
			fmt.Printf("Warning: %v\n", err)
			continue
		}

		fmt.Printf("%s has been added to the watchlist\n", ticker)
	}

	return nil
}

func removeTickersFromWatchlist(tickersArg string, watchlistService *watchlist.Service) error {
	// Split the tickers by comma
	tickers := strings.Split(tickersArg, ",")

	// Remove each ticker from the watchlist
	for _, ticker := range tickers {
		ticker = strings.TrimSpace(ticker)
		if ticker == "" {
			continue
		}

		if err := watchlistService.RemoveTicker(ticker); err != nil {
			fmt.Printf("Warning: %v\n", err)
			continue
		}

		fmt.Printf("%s has been removed from the watchlist\n", ticker)
	}

	return nil
}

func printUsage() {
	fmt.Println(getUsageText())
}

func getUsageText() string {
	return `StockTerm - A terminal-based stock viewer

Usage:
  stockterm <command> [arguments]

Commands:
  get <ticker>       Display stock price in a table. Multiple tickers can be separated by commas.
  get-all            Display watchlist in a table.
  list               Display an editable list of all tickers in the watchlist.
  add <ticker>       Add ticker to watchlist. Multiple tickers can be separated by commas.
  remove <ticker>    Remove ticker from watchlist. Multiple tickers can be separated by commas.
  help               Display this help message.
  version            Display version information.

Examples:
  stockterm get MSFT
  stockterm get AAPL,GOOGL,MSFT
  stockterm add TSLA
  stockterm add AAPL,META,TSLA
  stockterm remove TSLA
  stockterm get-all
  stockterm list`
}

func printVersion() {
	fmt.Println("StockTerm v1.0.0")
	fmt.Println("A terminal-based stock viewer")
	fmt.Println("https://github.com/AyhamJo7/StockTerm")
}
