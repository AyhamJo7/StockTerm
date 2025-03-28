package model

// ChartResponse represents the response from Yahoo Finance API
type ChartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency             string  `json:"currency"`
				Symbol               string  `json:"symbol"`
				ExchangeName         string  `json:"exchangeName"`
				InstrumentType       string  `json:"instrumentType"`
				FirstTradeDate       int     `json:"firstTradeDate"`
				RegularMarketTime    int     `json:"regularMarketTime"`
				HasPrePostMarketData bool    `json:"hasPrePostMarketData"`
				GMTOffset            int     `json:"gmtoffset"`
				Timezone             string  `json:"timezone"`
				ExchangeTimezoneName string  `json:"exchangeTimezoneName"`
				RegularMarketPrice   float64 `json:"regularMarketPrice"`
				ChartPreviousClose   float64 `json:"chartPreviousClose"`
				PreviousClose        float64 `json:"previousClose"`
				Scale                int     `json:"scale"`
				PriceHint            int     `json:"priceHint"`
				CurrentTradingPeriod struct {
					Pre     TradingPeriod `json:"pre"`
					Regular TradingPeriod `json:"regular"`
					Post    TradingPeriod `json:"post"`
				} `json:"currentTradingPeriod"`
				TradingPeriods  [][]TradingPeriod `json:"tradingPeriods"`
				DataGranularity string            `json:"dataGranularity"`
				Range           string            `json:"range"`
				ValidRanges     []string          `json:"validRanges"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Close []float64 `json:"close"`
					Low   []float64 `json:"low"`
					High  []float64 `json:"high"`
					Open  []float64 `json:"open"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}

// TradingPeriod represents a trading period in the market
type TradingPeriod struct {
	Timezone  string `json:"timezone"`
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	GMToffset int    `json:"gmtoffset"`
}

// StockData represents the essential stock data for display
type StockData struct {
	Ticker        string
	LastPrice     float64
	Change        float64
	ChangePercent float64
	PreviousClose float64
	Currency      string
}

// NewStockData creates a StockData instance from a ChartResponse
func NewStockData(response ChartResponse) StockData {
	if len(response.Chart.Result) == 0 {
		return StockData{}
	}

	meta := response.Chart.Result[0].Meta
	diff := meta.RegularMarketPrice - meta.PreviousClose
	changePercent := diff / meta.PreviousClose * 100

	return StockData{
		Ticker:        meta.Symbol,
		LastPrice:     meta.RegularMarketPrice,
		Change:        diff,
		ChangePercent: changePercent,
		PreviousClose: meta.PreviousClose,
		Currency:      meta.Currency,
	}
}
