package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kevinmajesta/karyawan/config"
	"kevinmajesta/karyawan/controllers"
	"kevinmajesta/karyawan/database"
	"kevinmajesta/karyawan/helpers"
	"kevinmajesta/karyawan/routes"
	"kevinmajesta/karyawan/worker"
)

var internalNotificationChan chan string

func main() {
	config.LoadEnv()

	helpers.InitRedis()
	db, err := database.InitPostgres()
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	database.DB = db
	helpers.SeedAdmin()
	helpers.SeedEmployees()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	controllers.EmailSenderInstance = helpers.NewEmailSender()

	r := routes.SetupRouter()

	// Inisialisasi Go Channel untuk notifikasi internal
	internalNotificationChan = make(chan string, 50)

	// Jalankan Worker Goroutine
	go worker.RunUploadWorker(ctx, internalNotificationChan)
	go worker.RunEmailWorkerFromRedis(ctx, internalNotificationChan)
	go worker.RunInternalNotificationWorker(ctx, internalNotificationChan) // Panggil fungsi baru ini

	// Jalankan Server Gin dengan Graceful Shutdown
	appPort := config.GetEnv("APP_PORT", "3000")
	srv := &http.Server{
		Addr:    ":" + appPort,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Tutup channel notifikasi saat shutdown
	close(internalNotificationChan)

	log.Println("Waiting for workers to finish...")
	time.Sleep(2 * time.Second)
	log.Println("Application exiting.")
}