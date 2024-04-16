// internal/app/app.go
package app

import (
	"ActualBudgetNormalizerV2/internal/models"
	"ActualBudgetNormalizerV2/internal/ui"
	"ActualBudgetNormalizerV2/internal/utils"
	"ActualBudgetNormalizerV2/pkg/gotawrapper"
	"fmt"
	"log"
	"sync"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

type SetsOfRows struct {
	Categories *models.Set
	Payees     *models.Set
	Notes      *models.Set
}

func NewSetsOfRows() *SetsOfRows {
	return &SetsOfRows{
		Categories: models.NewSet(),
		Payees:     models.NewSet(),
		Notes:      models.NewSet(),
	}
}

type App struct {
	originalDataFrame  dataframe.DataFrame
	processedDataFrame dataframe.DataFrame
	rowSets            *SetsOfRows
	chanUIIsLoading    chan bool
	chanUIProgressInfo chan string
	dataProcessingWg   sync.WaitGroup
}

func NewApp(
	originalDataFrame dataframe.DataFrame,
) *App {
	return &App{
		originalDataFrame:  originalDataFrame,
		rowSets:            NewSetsOfRows(),
		chanUIIsLoading:    make(chan bool),
		chanUIProgressInfo: make(chan string),
	}
}

func (a *App) Run() error {
	go ui.StartUI(a.chanUIIsLoading, a.chanUIProgressInfo)

	// Get the total number of rows in the DataFrame
	totalRows := a.originalDataFrame.Nrow()
	resultsChan := make(chan *utils.ResponseData)
	for i := 0; i < totalRows; i++ {
		a.dataProcessingWg.Add(1)
		go a.ProcessRow(i, totalRows, resultsChan)
	}

	// Wait for all goroutines to finish and close the result channel
	go func() {
		a.dataProcessingWg.Wait()
		close(resultsChan)
		a.chanUIIsLoading <- true
	}()

	// Collect the processed row results
	var processedData []*utils.ResponseData
	for result := range resultsChan {
		if result != nil {
			processedData = append(processedData, result)
		}
	}

	// Convert Processed Data back to a DataFrame
	a.processedDataFrame = a.CreateDataFrameFromResponseData(processedData)

	// Save the processed DataFrame to a CSV file
	if err := gotawrapper.SaveCSV(a.processedDataFrame, "DATA/EXTRATO_processed.csv"); err != nil {
		return fmt.Errorf("error saving processed DataFrame to CSV: %v", err)
	}

	return nil
}

func (a *App) ProcessRow(
	rowIndex, // The index of the row to process
	totalRows int, // The total number of rows in the DataFrame
	resultsChan chan<- *utils.ResponseData, // The channel to send the processed data to
) {
	defer a.dataProcessingWg.Done()

	rowStr := gotawrapper.RowToString(a.originalDataFrame, rowIndex)
	prompt := utils.PromptLLM(rowStr, map[string][]string{
		"categories": a.rowSets.Categories.Items(),
		"payee":      a.rowSets.Payees.Items(),
		"notes":      a.rowSets.Notes.Items(),
	})

	result, err := utils.QueryLLM("mistral", prompt)
	if err != nil {
		fmt.Errorf("Error querying LLM:", err)
		resultsChan <- nil
		return
	}
	resultsChan <- result // Send the processed data to the results channel

	const percentageFactor = 100 // Used to convert progress to a percentage
	progress := float64(rowIndex+1) / float64(totalRows)
	progressPercent := fmt.Sprintf("%.2f%%", progress*percentageFactor)
	a.chanUIProgressInfo <- fmt.Sprintf(
		"Processing row %d of %d (%s)",
		rowIndex+1,
		totalRows,
		progressPercent,
	)

	// Record the new data in the row sets
	updates := map[string]string{
		"categories": result.Category,
		"payee":      result.Payee,
		"notes":      result.Notes,
	}
	a.rowSets.Categories.Add(updates["categories"])
	a.rowSets.Payees.Add(updates["payee"])
	a.rowSets.Notes.Add(updates["notes"])
}

func (a *App) CreateDataFrameFromResponseData(processedData []*utils.ResponseData) dataframe.DataFrame {
	// Check for empty data to prevent unnecessary processing
	if len(processedData) == 0 {
		log.Println("No processed data to create DataFrame.")
		return dataframe.DataFrame{}
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
				log.Printf("Error converting amount for item at index %d: %v\n", i, err)
				amounts[i] = "0.00" // Default to "0.00" on error
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
