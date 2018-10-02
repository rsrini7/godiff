package utils

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestColumnReorder(t *testing.T) {
	filePath := filepath.Join("data", "base-small.csv")

	t.Run("TestColumnReorder", func(t *testing.T) {
		columnCount := GetColumnCount(filePath)

		var reorderData []int
		for i := 0; i < columnCount; i++ {
			reorderData = append(reorderData, i)
		}
		utils.RandShuffle(reorderData)
		//[]int{0, 2, 1, 4, 3, 5, 6, 7, 8, 9, 10}
		ColumnReorder(filePath, reorderData)
	})
}

func TestGetColumnCount(t *testing.T) {
	t.Run("Get CSV Column Count", func(t *testing.T) {
		if got := GetColumnCount(filepath.Join("data", "base-small.csv")); got == 0 {
			t.Errorf("GetColumnCount() = %v", got)
		} else {
			fmt.Printf("Given CSV Column Count : %d", got)
		}
	})
}
