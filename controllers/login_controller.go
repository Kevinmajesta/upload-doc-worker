package controllers

import (
	"net/http"
	"kevinmajesta/karyawan/database"
	"kevinmajesta/karyawan/helpers"
	"kevinmajesta/karyawan/models"
	"kevinmajesta/karyawan/structs"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {

	var req = structs.AdminLoginRequest{}
	var user = models.Admin{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "User Not Found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Invalid Password",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	token := helpers.GenerateToken(user.Username)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "Login Success",
		Data: structs.AdminResponse{
			Id:        user.Id,
			Username:  user.Username,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
			Token:     &token,
		},
	})
}