// internal/database/db.go
package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"startweektool/internal/models"
)

type DB struct {
	*sql.DB
}

func InitDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// Create tables
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS students (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            student_number TEXT UNIQUE NOT NULL,
            first_name TEXT NOT NULL,
            last_name TEXT NOT NULL,
            classroom TEXT NOT NULL,
            coach TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );
        CREATE INDEX IF NOT EXISTS idx_student_number ON students(student_number);
        CREATE INDEX IF NOT EXISTS idx_first_name ON students(first_name);
        CREATE INDEX IF NOT EXISTS idx_last_name ON students(last_name);
    `)

	if err != nil {
		return nil, fmt.Errorf("error creating tables: %v", err)
	}

	return &DB{db}, nil
}

func (db *DB) SaveStudents(students []models.Student) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
        INSERT OR REPLACE INTO students 
        (student_number, first_name, last_name, classroom, coach) 
        VALUES (?, ?, ?, ?, ?)
    `)
	if err != nil {
		return 0, fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	processed := 0
	for _, student := range students {
		_, err := stmt.Exec(
			student.StudentNumber,
			student.FirstName,
			student.LastName,
			student.Classroom,
			student.Coach,
		)
		if err != nil {
			return processed, fmt.Errorf("error inserting student %s: %v", student.StudentNumber, err)
		}
		processed++
	}

	if err := tx.Commit(); err != nil {
		return processed, fmt.Errorf("error committing transaction: %v", err)
	}

	return processed, nil
}

func (db *DB) GetAllStudents() ([]models.Student, error) {
	rows, err := db.Query(`
        SELECT id, student_number, first_name, last_name, classroom, coach 
        FROM students 
        ORDER BY student_number 
        LIMIT 100
    `)
	if err != nil {
		return nil, fmt.Errorf("error querying students: %v", err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		err := rows.Scan(&s.ID, &s.StudentNumber, &s.FirstName, &s.LastName, &s.Classroom, &s.Coach)
		if err != nil {
			return nil, fmt.Errorf("error scanning student row: %v", err)
		}
		students = append(students, s)
	}
	return students, nil
}

func (db *DB) SearchStudents(query string) ([]models.Student, error) {
	rows, err := db.Query(`
        SELECT id, student_number, first_name, last_name, classroom, coach 
        FROM students 
        WHERE student_number LIKE ? 
           OR first_name LIKE ? 
           OR last_name LIKE ?
        ORDER BY student_number
        LIMIT 50
    `, "%"+query+"%", "%"+query+"%", "%"+query+"%")

	if err != nil {
		return nil, fmt.Errorf("error searching students: %v", err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		err := rows.Scan(&s.ID, &s.StudentNumber, &s.FirstName, &s.LastName, &s.Classroom, &s.Coach)
		if err != nil {
			return nil, fmt.Errorf("error scanning search result: %v", err)
		}
		students = append(students, s)
	}
	return students, nil
}

func (db *DB) Close() error {
	if db.DB != nil {
		return db.DB.Close()
	}
	return nil
}
