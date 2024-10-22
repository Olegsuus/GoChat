package handlers

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"new_password" binding:"required,min=6"`
}

type RegisterNewUserDTO struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPasswordDTO struct {
	Email       string `json:"email" binding:"required,email"`
	SecretWord  string `json:"secret_word" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
