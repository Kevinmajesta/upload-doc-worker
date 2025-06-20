package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"kevinmajesta/karyawan/database"
	"kevinmajesta/karyawan/models"
	"kevinmajesta/karyawan/structs"

	"github.com/gin-gonic/gin"
)

func UpdateDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var doc models.Document
	if err := database.DB.First(&doc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	title := c.PostForm("title")
	employeeIDStr := c.PostForm("employee_id")

	if title != "" {
		doc.Title = title
	}

	if employeeIDStr != "" {
		if employeeID, err := strconv.ParseUint(employeeIDStr, 10, 64); err == nil {
			doc.EmployeeID = uint(employeeID)
		}
	}

	// Cek jika ada file baru
	file, err := c.FormFile("file")
	if err == nil {

		os.Remove(doc.FilePath)

		uploadDir := "./uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.Mkdir(uploadDir, os.ModePerm)
		}

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		fullPath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, fullPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new file"})
			return
		}

		doc.FilePath = fullPath
	}

	doc.UpdatedAt = time.Now()

	if err := database.DB.Save(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document"})
		return
	}

	var result models.Document
	if err := database.DB.Preload("Employee").First(&result, doc.Id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to load employee data",
			Errors:  nil,
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "Document updated successfully",
		Data:    result,
	})
}

func DeleteDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var doc models.Document
	if err := database.DB.Preload("Employee").First(&doc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	if err := os.Remove(doc.FilePath); err != nil {
		fmt.Println("Failed to delete file:", err)
	}

	if err := database.DB.Delete(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to delete document",
			Errors:  nil,
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "Document deleted",
		Data:    nil,
	})
}
