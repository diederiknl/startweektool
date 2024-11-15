// internal/handlers/admin.go
package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"startweektool/internal/database"
	"startweektool/internal/services"
	"startweektool/internal/services/search"
)

type AdminHandler struct {
	db            *database.DB
	searchService *search.Service
}

func NewAdminHandler(db *database.DB) *AdminHandler {
	return &AdminHandler{
		db:            db,
		searchService: search.NewService(db),
	}
}

func (h *AdminHandler) UploadCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geen bestand ontvangen"})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kan bestand niet openen"})
		return
	}
	defer openedFile.Close()

	students, err := services.ParseStudentsCSV(openedFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fout bij verwerken CSV: " + err.Error()})
		return
	}

	processed, err := h.searchService.SaveStudents(students)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fout bij opslaan: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Upload succesvol",
		"processed": processed,
	})
}

func (h *AdminHandler) GetPreview(c *gin.Context) {
	students, err := h.searchService.GetAllStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fout bij ophalen data: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, students)
}

func (h *AdminHandler) SearchStudents(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geen zoekterm opgegeven"})
		return
	}

	students, err := h.searchService.SearchStudents(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fout bij zoeken: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}
