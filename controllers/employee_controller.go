package controllers

import (
	"kevinmajesta/karyawan/database"
	"kevinmajesta/karyawan/helpers"
	"kevinmajesta/karyawan/models"
	"kevinmajesta/karyawan/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateEmployee(c *gin.Context) {
	var req structs.EmployeeCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Success: false,
			Message: "Validasi Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	employee := models.Employee{
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Position:   req.Position,
		Department: req.Department,
	}

	if err := database.DB.Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Gagal membuat data employee",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Code:    http.StatusCreated,
		Success: true,
		Message: "Employee berhasil dibuat",
		Data: structs.EmployeeResponse{
			Id:         employee.Id,
			Name:       employee.Name,
			Email:      employee.Email,
			Phone:      employee.Phone,
			Position:   employee.Position,
			Department: employee.Department,
			CreatedAt:  employee.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  employee.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func UpdateEmployee(c *gin.Context) {
	var req structs.EmployeeUpdateRequest

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "ID tidak valid",
			Errors:  nil,
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Success: false,
			Message: "Validasi Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var employee models.Employee
	if err := database.DB.First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Employee tidak ditemukan",
			Errors:  nil,
		})
		return
	}

	employee.Name = req.Name
	employee.Email = req.Email
	employee.Phone = req.Phone
	employee.Position = req.Position
	employee.Department = req.Department

	if err := database.DB.Save(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Gagal memperbarui data employee",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "Employee berhasil diperbarui",
		Data: structs.EmployeeResponse{
			Id:         employee.Id,
			Name:       employee.Name,
			Email:      employee.Email,
			Phone:      employee.Phone,
			Position:   employee.Position,
			Department: employee.Department,
			CreatedAt:  employee.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  employee.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func DeleteEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "ID tidak valid",
			Errors:  nil,
		})
		return
	}

	var employee models.Employee
	if err := database.DB.First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Code:    http.StatusNotFound,
			Success: false,
			Message: "Employee tidak ditemukan",
			Errors:  nil,
		})
		return
	}

	if err := database.DB.Delete(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Gagal menghapus employee",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "Employee berhasil dihapus",
		Data:    nil,
	})
}
