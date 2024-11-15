// internal/models/student.go
package models

type Student struct {
	ID            int    `json:"id" db:"id"`
	StudentNumber string `json:"student_number" db:"student_number"`
	FirstName     string `json:"first_name" db:"first_name"`
	LastName      string `json:"last_name" db:"last_name"`
	Classroom     string `json:"classroom" db:"classroom"`
	Coach         string `json:"coach" db:"coach"`
}
