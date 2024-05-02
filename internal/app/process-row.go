package app

import (
	"ActualBudgetNormalizerV2/internal/utils"
	"ActualBudgetNormalizerV2/pkg/gotawrapper"
	"fmt"
)

// ProcessRow processes a single row of the original DataFrame. It constructs a prompt for the LLM,
// queries the LLM for processed data, updates internal state based on the LLM's response, and
// manages progress reporting.
//
// Parameters:
// - rowIndex: The current index of the row being processed.
// - totalRows: The total number of rows in the DataFrame.
//
// Returns:
// - A pointer to the processed data or nil if an error occurs during processing.
func (a *App) ProcessRow(rowIndex, totalRows int) *utils.ResponseData {
	// Convert the row data to a string format suitable for LLM processing.
	rowStr := gotawrapper.RowToString(a.originalDataFrame, rowIndex)

	// Construct a prompt with current row data and existing categorizations.
	prompt := utils.PromptLLM(rowStr, map[string][]string{
		"categories": a.setsOfRows.Categories.Items(),
		"payee":      a.setsOfRows.Payees.Items(),
		"notes":      a.setsOfRows.Notes.Items(),
	})

	// Query the LLM with the constructed prompt to get the processed data.
	result, err := utils.QueryLLM(a.llmName, prompt)
	if err != nil {
		// Log the error and return nil to indicate the processing failure for this row.
		fmt.Printf("Error processing row %d: %v\n", rowIndex, err)
		return nil
	}

	// Update categorization sets based on the new data received from the LLM.
	updates := map[string]string{
		"categories": result.Category,
		"payee":      result.Payee,
		"notes":      result.Notes,
	}
	a.setsOfRows.Categories.Add(updates["categories"])
	a.setsOfRows.Payees.Add(updates["payee"])
	a.setsOfRows.Notes.Add(updates["notes"])

	// Return the result from the LLM query to be further processed or stored.
	return result
}
