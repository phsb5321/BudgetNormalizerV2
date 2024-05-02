package main

import (
	"log"
	"os"
	"path/filepath"

	"ActualBudgetNormalizerV2/internal/app"
	"ActualBudgetNormalizerV2/pkg/gotawrapper"
)

func main() {
	// Define the path to the log file
	logPath := filepath.Join("logs", "application.log")

	// Open or create the log file
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Set the log output to the file
	log.SetOutput(logFile)

	// Load the CSV file
	csvPath := filepath.Join("DATA", "EXTRATO.csv")
	dataFrame, err := gotawrapper.LoadCSV(csvPath)
	if err != nil {
		log.Fatalf("Error loading CSV from path %s: %v", csvPath, err)
	}

	// Create and run the application instance
	appInstance := app.NewApp(dataFrame, "phi3")
	if err := appInstance.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}

	log.Println("Data processed and saved successfully.")
}
