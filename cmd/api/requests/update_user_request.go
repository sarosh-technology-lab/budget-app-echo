package requests

type UpdateUserRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,max=50"`
    LastName  string `json:"last_name" validate:"omitempty,max=50"`
    Email     string `json:"email" validate:"omitempty,email,max=100"`
    Phone     string `json:"phone" validate:"omitempty,numeric,len=11"`
    Address   string `json:"address" validate:"omitempty"`
    Gender    string `json:"gender" validate:"omitempty,oneof=M F O"` // Matches ENUM
}