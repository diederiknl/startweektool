// internal/services/parser.go
package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"startweektool/internal/models"
)

func ParseStudentsCSV(reader io.Reader) ([]models.Student, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %v", err)
	}

	// Validate header
	expectedHeader := []string{"studentnummer", "voornaam", "achternaam", "lokaal", "coach"}
	if len(header) != len(expectedHeader) {
		return nil, fmt.Errorf("invalid CSV format: expected %d columns, got %d", len(expectedHeader), len(header))
	}

	var students []models.Student
	lineNumber := 1 // Start at 1 because we already read the header

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV line %d: %v", lineNumber, err)
		}

		if len(record) != 5 {
			return nil, fmt.Errorf("invalid number of columns at line %d: expected 5, got %d", lineNumber, len(record))
		}

		student := models.Student{
			StudentNumber: record[0],
			FirstName:     record[1],
			LastName:      record[2],
			Classroom:     record[3],
			Coach:         record[4],
		}

		// Basic validation
		if student.StudentNumber == "" {
			return nil, fmt.Errorf("empty student number at line %d", lineNumber)
		}

		students = append(students, student)
		lineNumber++
	}

	return students, nil
}
