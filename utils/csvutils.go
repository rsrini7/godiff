package utils

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

//ColumnReorder is to reorder the CSV columns
func ColumnReorder(filePath string, columns []int) {

	buf := bytes.Buffer{}
	//defer buf.Reset()

	file, err := os.Open(filePath)
	if err != nil {
		_ = fmt.Errorf("error in reading file %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var newColumn []string

	/*writer, wFile := getWriter(filePath)
	defer wFile.Close()*/

	for line, err := reader.Read(); err == nil; line, err = reader.Read() {
		for _, v := range columns {
			newColumn = append(newColumn, line[v])
		}

		if _, err = buf.WriteString(strings.Join(newColumn, ",") + "\n"); err != nil {
			fmt.Println("Error:", err)
			break
		}

		/*if err = writer.Write(newColumn); err != nil {
			fmt.Println("Error:", err)
			break
		}
		writer.Flush()*/

		newColumn = newColumn[:0]
	}

	writeToFile(filePath+".colreordered", buf.Bytes())
}

func writeToFile(filePath string, buf []byte) {
	wFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("writeToFile Error: ", err)
		return
	}
	defer wFile.Close()

	outFile := bufio.NewWriter(wFile)
	outFile.Write(buf)
	outFile.Flush()
}

/*func getWriter(filePath string) (*csv.Writer, *os.File) {
	// Creating csv writer
	wFile, err := os.Create(filePath + ".colreordered")
	if err != nil {
		fmt.Println("Error:", err)
		return nil, nil
	}

	writer := csv.NewWriter(wFile)
	return writer, wFile
}*/

//GetColumnCount : return the CSV column count
func GetColumnCount(filePath string) int {
	line := GetHeader(filePath)
	return len(line)
}

//GetHeader : return the CSV header
func GetHeader(filePath string) []string {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		_ = fmt.Errorf("error in reading file %s", err)
	}

	reader := csv.NewReader(file)
	line, err := reader.Read()

	if err != nil {
		_ = fmt.Errorf("error in reading CSV file %s", err)
	}
	return line
}

//HeaderPositionEqual :test whether the given two string slices are equal
func HeaderPositionEqual(a, b []int) bool {

	if (a != nil) && (b != nil) && (len(a) == len(b)) {
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
	}
	return true
}
