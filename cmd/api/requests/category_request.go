package requests

type CategoryRequest struct {
	Name string `json:"name" validate:"required,max=255"`
	IsCustom bool `json:"is_custom"`
}

type IDParamRequest struct {
	ID uint `param:"id" binding:"required"`
}