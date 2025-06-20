package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"kevinmajesta/karyawan/database"
	"kevinmajesta/karyawan/helpers"
	"kevinmajesta/karyawan/models"
	"kevinmajesta/karyawan/structs"

	"github.com/gin-gonic/gin"
)

func GetDocuments(c *gin.Context) {
	cacheKey := "documents:all"

	cached, err := helpers.RedisClient.Get(helpers.Ctx, cacheKey).Result()
	if err == nil {
		var docs []models.Document
		if err := json.Unmarshal([]byte(cached), &docs); err == nil {
			c.JSON(http.StatusOK, structs.SuccessResponse{
				Code:    http.StatusOK,
				Success: true,
				Message: "Documents fetched from cache",
				Data:    docs,
			})
			return
		}
	}

	var docs []models.Document
	if err := database.DB.Preload("Employee").Find(&docs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.SuccessResponse{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to fetch documents",
			Data:    nil,
		})
		return
	}

	jsonData, _ := json.Marshal(docs)
	helpers.RedisClient.Set(helpers.Ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "Documents fetched from database",
		Data:    docs,
	})
}

func GetDocumentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.SuccessResponse{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid document ID",
			Data:    nil,
		})
		return
	}

	cacheKey := "document:" + idStr

	cached, err := helpers.RedisClient.Get(helpers.Ctx, cacheKey).Result()
	if err == nil {
		var doc models.Document
		if err := json.Unmarshal([]byte(cached), &doc); err == nil {
			c.JSON(http.StatusOK, structs.SuccessResponse{
				Code:    http.StatusOK,
				Success: true,
				Message: "Document fetched from cache",
				Data:    doc,
			})
			return
		}
	}

	var doc models.Document
	if err := database.DB.Preload("Employee").First(&doc, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.SuccessResponse{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Document not found",
			Data:    nil,
		})
		return
	}

	jsonData, _ := json.Marshal(doc)
	helpers.RedisClient.Set(helpers.Ctx, cacheKey, jsonData, 10*time.Minute)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "Document fetched from database",
		Data:    doc,
	})
}
