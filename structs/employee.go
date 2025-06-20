package structs

type DocumentResponse struct {
	Id         uint   `json:"id"`
	Title      string `json:"title"`
	FilePath   string `json:"file_path"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type EmployeeResponse struct {
	Id         uint               `json:"id"`
	Name       string             `json:"name"`
	Email      string             `json:"email"`
	Phone      string             `json:"phone"`
	Position   string             `json:"position"`
	Department string             `json:"department"`
	CreatedAt  string             `json:"created_at"`
	UpdatedAt  string             `json:"updated_at"`
	Documents  []DocumentResponse `json:"documents"`
}


type EmployeeCreateRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Phone      string `json:"phone" binding:"omitempty"`
	Position   string `json:"position" binding:"required"`
	Department string `json:"department" binding:"required"`
}


type EmployeeUpdateRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Phone      string `json:"phone" binding:"omitempty"`
	Position   string `json:"position" binding:"required"`
	Department string `json:"department" binding:"required"`
}
