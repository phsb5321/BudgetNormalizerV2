// pkg/gotawrapper/gotawrapper_test.go
package gotawrapper

import (
	"path/filepath"
	"testing"
)

func TestLoadCSV(t *testing.T) {
	csvPath := filepath.Join("..", "..", "DATA", "NU_20848292_01ABR2024_08ABR2024.csv")
	df, err := LoadCSV(csvPath)
	if err != nil {
		t.Errorf("LoadCSV() error = %v, wantErr %v", err, false)
		return
	}
	if df.Nrow() == 0 {
		t.Error("Expected non-empty DataFrame")
	}
}

func TestRowToString(t *testing.T) {
	csvPath := filepath.Join("..", "..", "DATA", "NU_20848292_01ABR2024_08ABR2024.csv")
	df, err := LoadCSV(csvPath)
	if err != nil {
		t.Errorf("LoadCSV() error = %v, wantErr %v", err, false)
		return
	}

	rowIndex := 0
	expectedStr := "Data: 03/04/2024; Valor: 1100; Identificador: 660d6971-8a13-46d9-8244-07f75d90a630; Descrição: Transferência recebida pelo Pix - DEBORA ANACLETO VIEIRA DA CUNHA - •••.850.904-•• - BANCO INTER (0077) Agência: 1 Conta: 2888538-4; "
	rowStr := RowToString(df, rowIndex)
	if rowStr != expectedStr {
		t.Errorf("RowToString() got = %v, want %v", rowStr, expectedStr)
	}
}

func TestSaveCSV(t *testing.T) {
	csvPath := filepath.Join("..", "..", "testdata", "testdata.csv")
	df, err := LoadCSV(csvPath)
	if err != nil {
		t.Errorf("LoadCSV() error = %v, wantErr %v", err, false)
		return
	}

	newCsvPath := filepath.Join("..", "..", "testdata", "testdata_saved.csv")
	err = SaveCSV(df, newCsvPath)
	if err != nil {
		t.Errorf("SaveCSV() error = %v, wantErr %v", err, false)
	}
}
