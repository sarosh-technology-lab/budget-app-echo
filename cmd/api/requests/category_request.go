package requests

// CategoryRequestable defines the interface for any request that can be used to create a category
type CategoryRequestable interface {
    GetName() string
    GetIsCustom() bool
}

// Make CategoryRequest implement CategoryRequestable
func (r CategoryRequest) GetName() string {
    return r.Name
}

func (r CategoryRequest) GetIsCustom() bool {
    return r.IsCustom
}

// Make CategoryFormRequest implement CategoryRequestable
func (r CategoryFormRequest) GetName() string {
    return r.Name
}

func (r CategoryFormRequest) GetIsCustom() bool {
    return true // Default value since CategoryFormRequest doesn't have IsCustom field
}


type CategoryRequest struct {
	Name string `json:"name" validate:"required,max=255"`
	IsCustom bool `json:"is_custom"`
}

type IDParamRequest struct {
	ID uint `param:"id" binding:"required"`
}

type CategoryFormRequest struct {
	Name string `form:"name" validate:"required,max=255"`
}