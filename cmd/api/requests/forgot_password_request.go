package requests

type ForgotPasswordRequest struct{
	Email string `json:"email" validate:"required,email"`
	FrontendURL string `json:"frontend_url" validate:"required,url"`
}