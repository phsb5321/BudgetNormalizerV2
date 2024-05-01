// internal/app/app.go
package app

import (
	"ActualBudgetNormalizerV2/internal/models"
	"ActualBudgetNormalizerV2/internal/utils"
	"ActualBudgetNormalizerV2/pkg/gotawrapper"
	"fmt"

	"github.com/go-gota/gota/dataframe"
)

// SetsOfRows holds sets of categories, payees, and notes used across the application.
type SetsOfRows struct {
	Categories *models.Set
	Payees     *models.Set
	Notes      *models.Set
}

// NewSetsOfRows initializes new sets for categories, payees, and notes.
func NewSetsOfRows() *SetsOfRows {
	return &SetsOfRows{
		Categories: models.NewSet(),
		Payees:     models.NewSet(),
		Notes:      models.NewSet(),
	}
}

// App encapsulates all the necessary components of the application.
type App struct {
	llmName            string
	originalDataFrame  dataframe.DataFrame
	processedDataFrame dataframe.DataFrame
	setsOfRows         *SetsOfRows
	chanUIProgressInfo chan float64
}

// NewApp creates a new application instance with necessary initializations.
func NewApp(originalDataFrame dataframe.DataFrame, llmName string) *App {
	return &App{
		originalDataFrame:  originalDataFrame,
		setsOfRows:         NewSetsOfRows(),
		chanUIProgressInfo: make(chan float64),
		llmName:            llmName,
	}
}

// Run starts the application's main logic, processing each row and handling the UI.
func (a *App) Run() error {
	// go ui.StartUI(a.chanUIProgressInfo)

	totalRows := a.originalDataFrame.Nrow()
	var processedData []*utils.ResponseData
	for i := 0; i < totalRows; i++ {
		result := a.ProcessRow(i, totalRows)
		if result != nil {
			processedData = append(processedData, result)
		}
	}

	a.processedDataFrame = a.CreateAppDF(processedData)
	if err := gotawrapper.SaveCSV(a.processedDataFrame, "DATA/EXTRATO_processed.csv"); err != nil {
		return fmt.Errorf("error saving processed DataFrame to CSV: %v", err)
	}

	return nil
}
