package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"stockterm/internal/model"
)

// TableRenderer renders stock data in a table
type TableRenderer struct {
	writer io.Writer
	style  table.Style
}

// NewTableRenderer creates a new table renderer
func NewTableRenderer() *TableRenderer {
	return &TableRenderer{
		writer: os.Stdout,
		style:  table.StyleColoredCyanWhiteOnBlack,
	}
}

// WithWriter sets the writer for the table renderer
func (r *TableRenderer) WithWriter(writer io.Writer) *TableRenderer {
	r.writer = writer
	return r
}

// WithStyle sets the style for the table renderer
func (r *TableRenderer) WithStyle(style table.Style) *TableRenderer {
	r.style = style
	return r
}

// RenderStocks renders a table of stock data
func (r *TableRenderer) RenderStocks(stocks []model.StockData) {
	t := table.NewWriter()
	t.SetOutputMirror(r.writer)
	t.AppendHeader(table.Row{"Ticker", "Last Price", "Change", "Change %", "Prev. Close", "Currency"})

	for _, stock := range stocks {
		t.AppendRow(table.Row{
			stock.Ticker,
			fmt.Sprintf("%.2f", stock.LastPrice),
			appendPlus(stock.Change),
			appendPlus(stock.ChangePercent),
			fmt.Sprintf("%.2f", stock.PreviousClose),
			stock.Currency,
		})
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name: "Change",
			Transformer: text.Transformer(func(val interface{}) string {
				return getColoredChangeCell(val, "")
			}),
		},
		{
			Name: "Change %",
			Transformer: text.Transformer(func(val interface{}) string {
				return getColoredChangeCell(val, "%")
			}),
		},
	})

	t.SetStyle(r.style)
	t.Render()
}

// RenderChartResponses renders a table of chart responses
func (r *TableRenderer) RenderChartResponses(responses []model.ChartResponse) {
	var stocks []model.StockData
	for _, response := range responses {
		stocks = append(stocks, model.NewStockData(response))
	}
	r.RenderStocks(stocks)
}

// getColoredChangeCell returns a colored cell for a change value
func getColoredChangeCell(val interface{}, postfix string) string {
	strVal, ok := val.(string)
	if !ok {
		return "0.00" + postfix
	}

	var color text.Color
	if strings.Contains(strVal, "-") {
		color = text.FgRed
	} else if strings.Contains(strVal, "+") {
		color = text.FgGreen
	}

	return text.Colors{color}.Sprint(strVal + postfix)
}

// appendPlus adds a plus sign to positive numbers
func appendPlus(num float64) string {
	if num >= 0 {
		return fmt.Sprintf("+%.2f", num)
	}
	return fmt.Sprintf("%.2f", num)
}
