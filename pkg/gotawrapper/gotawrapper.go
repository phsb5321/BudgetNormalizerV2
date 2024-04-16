// pkg/gotawrapper/gotawrapper.go
package gotawrapper

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

func LoadCSV(filePath string) (dataframe.DataFrame, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file %s: %v", filePath, err)
		return dataframe.DataFrame{}, err
	}
	defer file.Close()

	df := dataframe.ReadCSV(file)
	if df.Err != nil {
		log.Printf("Failed to read CSV data from file %s: %v", filePath, df.Err)
		return dataframe.DataFrame{}, df.Err
	}

	return df, nil
}

func RowToString(df dataframe.DataFrame, rowIndex int) string {
	var strBuilder strings.Builder
	if rowIndex < 0 || rowIndex >= df.Nrow() {
		return ""
	}

	for _, colName := range df.Names() {
		strBuilder.WriteString(fmt.Sprintf("%s: %v; ", colName, df.Col(colName).Elem(rowIndex)))
	}

	return strBuilder.String()
}

func SaveCSV(df dataframe.DataFrame, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file %s: %v", filePath, err)
		return err
	}
	defer file.Close()

	err = df.WriteCSV(file)
	if err != nil {
		log.Printf("Failed to write CSV data to file %s: %v", filePath, err)
		return err
	}

	return nil
}

func MutateDataFrameWithSeries(
	df dataframe.DataFrame,
	data []string, columnType string,
	columnName string,
) dataframe.DataFrame {
	var col series.Series
	switch columnType {
	case "string":
		col = series.New(data, series.String, columnName)
	case "float":
		col = series.New(data, series.Float, columnName)
	case "int":
		col = series.New(data, series.Int, columnName)
	default:
		return df
	}

	return df.Mutate(col)
}
