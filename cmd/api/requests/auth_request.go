package requests

type RegisterUserRequest struct{
	FirstName string `json:"first_name" validate:"required,alpha,max=50"`
	LastName string `json:"last_name" validate:"required,alpha,max=50"`
	Email string `json:"email" validate:"required,email,max=100"`
	Phone     string `json:"phone" validate:"omitempty,numeric,len=11"`
	Password string `json:"password" validate:"required,min=2,max=8"`
}

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=2,max=8"`
}