// internal/services/search/service.go
package search

import (
	"startweektool/internal/database"
	"startweektool/internal/models"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) SearchStudents(query string) ([]models.Student, error) {
	if query == "" {
		return nil, nil
	}
	return s.db.SearchStudents(query)
}

func (s *Service) GetAllStudents() ([]models.Student, error) {
	return s.db.GetAllStudents()
}

func (s *Service) SaveStudents(students []models.Student) (int, error) {
	return s.db.SaveStudents(students)
}
