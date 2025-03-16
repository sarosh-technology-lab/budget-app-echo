package requests

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=2,max=8"`
	Password string `json:"password" validate:"required,min=2,max=8"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
}

type ResetPasswordRequest struct {
	Password string `json:"password" validate:"required,min=2,max=8"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
	Token string `json:"token" validate:"required,min=5"`
	Meta string `json:"meta" validate:"required"`
}