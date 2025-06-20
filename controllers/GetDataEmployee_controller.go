package controllers

import (
	"encoding/json"
	"fmt"
	"kevinmajesta/karyawan/database"
	"kevinmajesta/karyawan/helpers"
	"kevinmajesta/karyawan/models"
	"kevinmajesta/karyawan/structs"

	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllEmployees(c *gin.Context) {
	cacheKey := "employees:all"

	cached, err := helpers.RedisClient.Get(helpers.Ctx, cacheKey).Result()
	if err == nil {
		var employees []structs.EmployeeResponse
		if err := json.Unmarshal([]byte(cached), &employees); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"success": true,
				"data":    employees,
				"source":  "cache",
			})
			return
		}
	}

	var employees []models.Employee
	if err := database.DB.Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get employees"})
		return
	}

	var response []structs.EmployeeResponse
	for _, e := range employees {
		response = append(response, structs.EmployeeResponse{
			Id:         e.Id,
			Name:       e.Name,
			Email:      e.Email,
			Phone:      e.Phone,
			Position:   e.Position,
			Department: e.Department,
			CreatedAt:  e.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  e.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	jsonData, _ := json.Marshal(response)
	helpers.RedisClient.Set(helpers.Ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"data":    response,
		"source":  "database",
	})
}

func GetEmployeeByID(c *gin.Context) {
	idStr := c.Param("id")
	cacheKey := fmt.Sprintf("employee:%s", idStr)

	cached, err := helpers.RedisClient.Get(helpers.Ctx, cacheKey).Result()
	if err == nil {
		var employee structs.EmployeeResponse
		if err := json.Unmarshal([]byte(cached), &employee); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"success": true,
				"data":    employee,
				"source":  "cache",
			})
			return
		}
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.SuccessResponse{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid employee ID",
			Data:    nil,
		})
		return
	}

	var employee models.Employee
	if err := database.DB.Preload("Documents").First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.SuccessResponse{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Employee not found",
			Data:    nil,
		})
		return
	}

	var documents []structs.DocumentResponse
	for _, d := range employee.Documents {
		documents = append(documents, structs.DocumentResponse{
			Id:        d.Id,
			Title:     d.Title,
			FilePath:  d.FilePath,
			CreatedAt: d.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: d.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	response := structs.EmployeeResponse{
		Id:         employee.Id,
		Name:       employee.Name,
		Email:      employee.Email,
		Phone:      employee.Phone,
		Position:   employee.Position,
		Department: employee.Department,
		CreatedAt:  employee.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  employee.UpdatedAt.Format("2006-01-02 15:04:05"),
		Documents:  documents,
	}

	jsonData, _ := json.Marshal(response)
	helpers.RedisClient.Set(helpers.Ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"data":    response,
		"source":  "database",
	})
}
