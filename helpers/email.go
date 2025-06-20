package helpers

import (
	"context"       // Tambahkan import context
	"encoding/json" // Tambahkan import json
	"fmt"
	"kevinmajesta/karyawan/config"
	"log"

	"gopkg.in/gomail.v2"
	"kevinmajesta/karyawan/jobs" 
)

// EmailSender struct tetap sama
type EmailSender struct {
	SMTPHost       string
	SMTPPort       int
	SenderEmail    string
	SenderPassword string
}

// NewEmailSender tetap sama
func NewEmailSender() *EmailSender {
	smtpHost := config.GetEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPortStr := config.GetEnv("SMTP_PORT", "587")
	senderEmail := config.GetEnv("SMTP_EMAIL", "")
	senderPassword := config.GetEnv("SMTP_PASSWORD", "")

	smtpPort, err := parsePort(smtpPortStr)
	if err != nil {
		log.Fatalf("Failed to parse SMTP port: %v", err)
	}

	return &EmailSender{
		SMTPHost:       smtpHost,
		SMTPPort:       smtpPort,
		SenderEmail:    senderEmail,
		SenderPassword: senderPassword,
	}
}

// parsePort tetap sama
func parsePort(portStr string) (int, error) {
	var port int
	_, err := fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		return 0, fmt.Errorf("invalid port format: %s", portStr)
	}
	return port, nil
}

// SendEmail tetap sama (ini adalah fungsi yang akan dipanggil oleh worker)
func (e *EmailSender) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.SenderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(e.SMTPHost, e.SMTPPort, e.SenderEmail, e.SenderPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email to %s: %v", to, err)
	}
	return nil
}

func (e *EmailSender) SendWelcomeEmail(to, name string) error {
	subject := "Selamat Datang di Karyawan App!"
	body := fmt.Sprintf("Halo %s,\n\nSelamat datang di aplikasi Karyawan! Kami sangat senang Anda bergabung dengan kami.\n\nSalam Hormat,\nTim Karyawan App", name)

	// Gunakan jobs.EmailJob
	emailJob := jobs.EmailJob{ // <<< UBAH INI: Gunakan jobs.EmailJob
		To:      to,
		Subject: subject,
		Body:    body,
		Name:    name,
	}

	jobJSON, err := json.Marshal(emailJob)
	if err != nil {
		return fmt.Errorf("failed to serialize email job: %w", err)
	}
	err = RedisClient.RPush(context.Background(), "email_jobs", jobJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to push email job to Redis: %w", err)
	}

	log.Printf("INFO: Welcome email job for %s queued successfully in Redis.\n", to)
	return nil
}