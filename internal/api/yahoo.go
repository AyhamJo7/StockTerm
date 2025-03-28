package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"stockterm/internal/model"
)

// YahooFinanceClient is a client for the Yahoo Finance API
type YahooFinanceClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewYahooFinanceClient creates a new Yahoo Finance API client
func NewYahooFinanceClient() *YahooFinanceClient {
	return &YahooFinanceClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://query1.finance.yahoo.com/v8/finance/chart/%s?region=US&lang=en-US&includePrePost=false&interval=2m&useYfid=true&range=%s&corsDomain=finance.yahoo.com&.tsrc=finance",
	}
}

// FetchStockData fetches stock data for a given ticker and time range
func (c *YahooFinanceClient) FetchStockData(ctx context.Context, ticker, timeRange string) (model.ChartResponse, error) {
	var response model.ChartResponse

	// Default to 1d if no time range is specified
	if timeRange == "" {
		timeRange = "1d"
	}

	// Create the URL
	url := fmt.Sprintf(c.baseURL, ticker, timeRange)

	// Create a new request with the provided context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return response, fmt.Errorf("error creating request: %w", err)
	}

	// Execute the request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return response, fmt.Errorf("error fetching data: %w", err)
	}
	defer res.Body.Close()

	// Check for non-200 status codes
	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Decode the response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("error decoding response: %w", err)
	}

	return response, nil
}

// FetchMultipleStocks fetches data for multiple tickers
func (c *YahooFinanceClient) FetchMultipleStocks(ctx context.Context, tickers []string, timeRange string) ([]model.ChartResponse, error) {
	var responses []model.ChartResponse

	for _, ticker := range tickers {
		res, err := c.FetchStockData(ctx, ticker, timeRange)
		if err != nil {
			// Log the error but continue with other tickers
			fmt.Printf("Error fetching data for %s: %v\n", ticker, err)
			continue
		}
		responses = append(responses, res)
	}

	return responses, nil
}
