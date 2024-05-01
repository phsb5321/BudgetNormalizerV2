package app

import (
	"ActualBudgetNormalizerV2/internal/utils"
	"fmt"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

func (a *App) CreateAppDF(processedData []*utils.ResponseData) dataframe.DataFrame {
	// Check for empty data to prevent unnecessary processing
	if len(processedData) == 0 {
		return dataframe.New()
	}

	// Initialize slices for DataFrame columns
	dates := make([]string, len(processedData))
	amounts := make([]string, len(processedData)) // Use string to store formatted amounts
	payees := make([]string, len(processedData))
	notes := make([]string, len(processedData))
	categories := make([]string, len(processedData))

	// Populate columns from processed data
	for i, item := range processedData {
		if item != nil {
			dates[i] = item.Date
			payees[i] = item.Payee
			notes[i] = item.Notes
			categories[i] = item.Category

			// Format amount to two decimal places as it is money related
			amt, err := item.Amount.Float64()
			if err != nil {
				amounts[i] = fmt.Sprintf("%.2f", 0.0) // Default to 0.0 if error
			} else {
				amounts[i] = fmt.Sprintf("%.2f", amt) // Use formatted string
			}
		}
	}

	// Create a new DataFrame with the augmented data
	df := dataframe.New(
		series.New(dates, series.String, "Date"),
		series.New(amounts, series.String, "Amount"), // Store amounts as strings
		series.New(payees, series.String, "Payee"),
		series.New(notes, series.String, "Notes"),
		series.New(categories, series.String, "Categories"),
	)

	return df
}
