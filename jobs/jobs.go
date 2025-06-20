package jobs // <<< Make sure this package name is 'jobs'

// Only import what's needed for these structs

// EmailJob struct untuk payload email yang akan diantrekan di Redis
type EmailJob struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Name    string `json:"name"`
}

// UploadJob struct untuk payload upload yang akan diantrekan di Redis
type UploadJob struct {
	Title      string `json:"title"`
	EmployeeID uint   `json:"employee_id"`
	FilePath   string `json:"file_path"`
}
