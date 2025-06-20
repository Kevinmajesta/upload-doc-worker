// --- package controllers ---
package controllers

import (
	"kevinmajesta/karyawan/database"
	"kevinmajesta/karyawan/helpers"
	"kevinmajesta/karyawan/models"
	"kevinmajesta/karyawan/structs"
	"log" // Tetap pakai log
	"net/http"

	"github.com/gin-gonic/gin"
)

// EmailSenderInstance: Pertahankan ini jika Anda masih menggunakannya di tempat lain.
// Namun, untuk Register, kita akan memanggil helper langsung.
// Jika EmailSenderInstance hanya digunakan di sini, pertimbangkan untuk menghapusnya
// dan langsung memanggil NewEmailSender di main lalu dilewatkan sebagai dependency
// atau inisialisasi di dalam fungsi yang membutuhkan, tapi untuk saat ini, kita ikuti strukur Anda.
var EmailSenderInstance *helpers.EmailSender // Tetap ada untuk konsistensi dengan kode Anda

func Register(c *gin.Context) {
	var req = structs.AdminCreateRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Success: false,
			Message: "Validasi Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	user := models.Admin{
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(req.Password),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		if helpers.IsDuplicateEntryError(err) {
			c.JSON(http.StatusConflict, structs.ErrorResponse{
				Code:    http.StatusConflict,
				Success: false,
				Message: "Duplicate entry error",
				Errors:  helpers.TranslateErrorMessage(err),
			})
		} else {
			c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Success: false,
				Message: "Failed to create Admin",
				Errors:  helpers.TranslateErrorMessage(err),
			})
		}
		return
	}

	// --- MODIFIED: Memanggil SendWelcomeEmail yang sekarang mendorong ke Redis Queue ---
	if EmailSenderInstance != nil {
		err := EmailSenderInstance.SendWelcomeEmail(req.Email, req.Username)
		if err != nil {
			// Penting: Jika gagal mendorong ke Redis, itu adalah masalah serius.
			// Log errornya, tapi mungkin masih bisa mengembalikan 201 OK ke user
			// karena registrasi berhasil, hanya pengiriman email yang antre gagal.
			// Tergantung seberapa kritis pengiriman email di alur ini.
			log.Printf("ERROR: Failed to queue welcome email for %s: %v", req.Email, err)
			// Anda bisa memilih untuk mengembalikan 503 jika antrean sangat penting:
			// c.JSON(http.StatusServiceUnavailable, structs.ErrorResponse{
			// 	Code:    http.StatusServiceUnavailable,
			// 	Success: false,
			// 	Message: "Admin created, but failed to queue welcome email. Please contact support.",
			// 	Errors:  []string{err.Error()},
			// })
			// return
		}
	} else {
		log.Println("WARNING: EmailSenderInstance is not initialized. Welcome email will not be queued.")
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Code:    http.StatusCreated,
		Success: true,
		Message: "Admin created successfully",
		Data: structs.AdminResponse{
			Id:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}