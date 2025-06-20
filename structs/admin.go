package structs


type AdminResponse struct {
	Id        uint    `json:"id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"` // <<< Add this line
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Token     *string `json:"token,omitempty"`
}

type AdminCreateRequest struct {
	Username string `json:"username" binding:"required" gorm:"unique;not null"`
	Email    string `json:"email" binding:"required,email"` // <<< Add this line and the 'email' binding
	Password string `json:"password" binding:"required"`
}

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}