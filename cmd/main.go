// cmd/main.go
package main

import (
	"ActualBudgetNormalizerV2/internal/app"
	"ActualBudgetNormalizerV2/pkg/gotawrapper"
	"log"
	"path/filepath"
)

func main() {
	// Define the CSV path
	csvPath := filepath.Join("DATA", "EXTRATO.csv")

	// Load the CSV data
	dataFrame, err := gotawrapper.LoadCSV(csvPath)
	if err != nil {
		log.Fatalf("Error loading CSV from path %s: %v", csvPath, err)
	}

	// Create a new application
	app := app.NewApp(dataFrame)

	// Run the application
	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}

	// Log successful execution
	log.Println("Data processed and saved successfully.")
}
