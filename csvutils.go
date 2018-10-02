package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

//ColumnReorder is to reorder the CSV columns
func ColumnReorder(filePath string, columns []int) {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		_ = fmt.Errorf("error in reading file %s", err)
	}

	reader := csv.NewReader(file)
	var newColumn []string

	writer, wFile := getWriter()
	defer wFile.Close()

	for line, err := reader.Read(); err == nil; line, err = reader.Read() {
		for _, v := range columns {
			newColumn = append(newColumn, line[v])
		}
		//fmt.Println(strings.Join(newColumn, ","))

		if err = writer.Write(newColumn); err != nil {
			fmt.Println("Error:", err)
			break
		}

		writer.Flush()
		newColumn = newColumn[:0]
	}
}

func getWriter() (*csv.Writer, *os.File) {
	// Creating csv writer

	wFile, err := os.Create("./temp000.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return nil, nil
	}

	writer := csv.NewWriter(wFile)
	return writer, wFile
}

//GetColumnCount - return the CSV column count
func GetColumnCount(filePath string) int {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		_ = fmt.Errorf("error in reading file %s", err)
	}

	reader := csv.NewReader(file)
	line, err := reader.Read()

	if err != nil {
		_ = fmt.Errorf("error in reading file %s", err)
	}
	return len(line)
}
