package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"kevinmajesta/karyawan/database"
	"kevinmajesta/karyawan/helpers"
	"kevinmajesta/karyawan/jobs"
	"kevinmajesta/karyawan/models"
)

func RunUploadWorker(ctx context.Context, notificationChan chan<- string) {
	log.Println("Upload worker started.")
	for {
		select {
		case <-ctx.Done():
			log.Println("Upload worker received shutdown signal from context.")
			return
		default:
			result, err := helpers.RedisClient.BLPop(ctx, 0*time.Second, "upload_jobs").Result()
			if err != nil {
				if err == context.Canceled {
					log.Println("Upload worker context canceled, exiting BLPop.")
					return
				}
				log.Printf("ERROR: Failed to pop upload job from Redis: %v\n", err)
				time.Sleep(1 * time.Second)
				continue
			}

			jobJSON := result[1]
			var job jobs.UploadJob
			if err := json.Unmarshal([]byte(jobJSON), &job); err != nil {
				log.Printf("ERROR: Failed to unmarshal upload job: %v. Job data: %s\n", err, jobJSON)
				continue
			}

			log.Printf("INFO: Processing upload job for file: %s (Title: %s)\n", job.FilePath, job.Title)

			destDir := "./uploads"
			if _, err := os.Stat(destDir); os.IsNotExist(err) {
				os.Mkdir(destDir, os.ModePerm)
			}

			filename := filepath.Base(job.FilePath)
			destPath := filepath.Join(destDir, filename)

			err = os.Rename(job.FilePath, destPath)
			if err != nil {
				log.Printf("ERROR: Failed to move file: %v\n", err)
				select {
				case notificationChan <- fmt.Sprintf("Upload FAILED for %s: %v", job.Title, err):
				default:
					log.Println("WARN: Notification channel full for failed upload.")
				}
				continue
			}

			doc := models.Document{
				Title:      job.Title,
				FilePath:   destPath,
				EmployeeID: job.EmployeeID,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			if err := database.DB.Create(&doc).Error; err != nil {
				log.Printf("ERROR: Failed to save document to DB: %v\n", err)
				select {
				case notificationChan <- fmt.Sprintf("DB save FAILED for %s: %v", job.Title, err):
				default:
					log.Println("WARN: Notification channel full for DB save failure.")
				}
				continue
			}

			log.Printf("INFO: Upload job processed successfully for file: %s\n", destPath)
			select {
			case notificationChan <- fmt.Sprintf("Upload SUCCESS: '%s' by Employee %d", job.Title, job.EmployeeID):
			case <-ctx.Done():
				return
			default:
				log.Println("WARN: Notification channel is full, skipping success notification for upload.")
			}
		}
	}
}

func RunEmailWorkerFromRedis(ctx context.Context, notificationChan chan<- string) {
	sender := helpers.NewEmailSender() // Initialize sender once per worker
	log.Println("Email worker started.")
	for {
		select {
		case <-ctx.Done():
			log.Println("Email worker received shutdown signal from context.")
			return
		default:
			result, err := helpers.RedisClient.BLPop(ctx, 0*time.Second, "email_jobs").Result()
			if err != nil {
				if err == context.Canceled {
					log.Println("Email worker context canceled, exiting BLPop.")
					return
				}
				log.Printf("ERROR: Failed to pop email job from Redis: %v\n", err)
				time.Sleep(1 * time.Second)
				continue
			}

			jobJson := result[1]
			var job jobs.EmailJob
			if err := json.Unmarshal([]byte(jobJson), &job); err != nil {
				log.Printf("ERROR: Failed to parse email job from Redis: %v. Job data: %s\n", err, jobJson)
				continue
			}

			log.Printf("INFO: Attempting to send welcome email to %s for user %s\n", job.To, job.Name)

			// Assuming EmailSenderInstance is correctly initialized in main and
			// NewEmailSender() is suitable for being called here (it is, as it creates a new sender)
			err = sender.SendEmail(job.To, job.Subject, job.Body)
			if err != nil {
				log.Printf("ERROR: Failed to send welcome email to %s: %v\n", job.To, err)
				select {
				case notificationChan <- fmt.Sprintf("Email FAILED to %s (Subject: %s): %v", job.To, job.Subject, err):
				default:
					log.Println("WARN: Notification channel full for failed email.")
				}
			} else {
				log.Printf("INFO: Welcome email sent successfully to %s for user %s\n", job.To, job.Name)
				select {
				case notificationChan <- fmt.Sprintf("Email SUCCESS: To %s (Subject: %s)", job.To, job.Subject):
				default:
					log.Println("WARN: Notification channel full for successful email.")
				}
			}
		}
	}
}

// runInternalNotificationWorker adalah goroutine worker baru untuk memproses notifikasi internal
func RunInternalNotificationWorker(ctx context.Context, notifications <-chan string) {
	log.Println("Internal notification worker started.")
	for {
		select {
		case <-ctx.Done():
			log.Println("Internal notification worker received shutdown signal from context. Exiting.")
			return
		case msg, ok := <-notifications:
			if !ok { // Channel ditutup
				log.Println("Internal notification channel closed. Exiting.")
				return
			}
			log.Printf("NOTIFICATION: %s\n", msg) // Menulis notifikasi ke log
		}
	}
}