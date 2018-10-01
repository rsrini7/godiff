package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func columnReorder(rFile *os.File, fromColumn int, toColumn int) {
	reader := csv.NewReader(rFile)

	writer := getWriter()
	for line, err := reader.Read(); err == nil; line, err = reader.Read() {
		//TODO: all columns should be specified with correct order based on input csv column length
		// need to refactor to used one time column re-order instead of swap just two columns
		if err = writer.Write([]string{line[fromColumn], line[toColumn]}); err != nil {
			fmt.Println("Error:", err)
			break
		}

		writer.Flush()
	}
}

func getWriter() *csv.Writer {
	// Creating csv writer
	wFile, err := os.Create("./temp000.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer wFile.Close()
	writer := csv.NewWriter(wFile)
	return writer
}
