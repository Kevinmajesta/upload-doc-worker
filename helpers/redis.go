// kevinmajesta/karyawan/helpers/redis.go
package helpers

import (
	"context"
	"encoding/json"
	"fmt" // Import config untuk GetEnv
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"kevinmajesta/karyawan/jobs" 
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {
	host := GetEnv("REDIS_HOST", "redis")
	port := GetEnv("REDIS_PORT", "6379")
	password := GetEnv("REDIS_PASSWORD", "")

	addr := fmt.Sprintf("%s:%s", host, port)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// CreateDocumentAsync - fungsi handler Gin ini tetap di helpers
func CreateDocumentAsync(c *gin.Context) {
	title := c.PostForm("title")
	employeeIDStr := c.PostForm("employee_id")

	employeeID64, err := strconv.ParseUint(employeeIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	uploadDir := "./uploads_temp"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	fullPath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save temporary file"})
		return
	}

	// Gunakan jobs.UploadJob
	job := jobs.UploadJob{ // <<< UBAH INI: Gunakan jobs.UploadJob
		Title:      title,
		EmployeeID: uint(employeeID64),
		FilePath:   fullPath,
	}

	jobJson, err := json.Marshal(job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode job data"})
		return
	}

	if err := RedisClient.RPush(Ctx, "upload_jobs", jobJson).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue job"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"success": true, "message": "Upload job queued"})
}
