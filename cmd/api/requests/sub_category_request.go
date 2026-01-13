package requests

// SubCategoryRequestable defines the interface for any request that can be used to create a category
type SubCategoryRequestable interface {
    GetName() string
    GetIsCustom() bool
    GetCategoryId() uint
    GetId() uint
}

// Make SubCategoryRequest implement SubCategoryRequestable
func (r SubCategoryRequest) GetName() string {
    return r.Name
}

func (r SubCategoryRequest) GetIsCustom() bool {
    return r.IsCustom
}

func (r SubCategoryRequest) GetCategoryId() uint {
    return r.CategoryId
}

func (r SubCategoryRequest) GetId() uint {
    return r.ID
}

// Make SubCategoryFormRequest implement SubCategoryRequestable
func (r SubCategoryFormRequest) GetName() string {
    return r.Name
}

func (r SubCategoryFormRequest) GetIsCustom() bool {
    return true // Default value since SubCategoryFormRequest doesn't have IsCustom field
}

type SubCategoryRequest struct {
	Name string `json:"name" validate:"required,max=255"`
	IsCustom bool `json:"is_custom"`
	CategoryId uint `json:"category_id" validate:"required"`
	ID uint `json:"id"`
}

type SubCategoryIDParamRequest struct {
	ID uint `param:"id" binding:"required"`
}

type SubCategoryFormRequest struct {
	Name string `form:"name" validate:"required,max=255"`
}