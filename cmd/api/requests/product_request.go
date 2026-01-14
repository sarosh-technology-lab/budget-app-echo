package requests

// ProductRequestable defines the interface for any request that can be used to create a category
type ProductRequestable interface {
    GetName() string
    GetDescription() string
    GetImage() string
    GetCategoryId() uint
    GetSubCategoryId() uint
    GetId() uint
}

// Make ProductRequest implement ProductRequestable
func (r ProductRequest) GetName() string {
    return r.Name
}

func (r ProductRequest) GetDescription() string {
    return r.Description
}

func (r ProductRequest) GetImage() string {
    return r.Image
}

func (r ProductRequest) GetCategoryId() uint {
    return r.CategoryId
}

func (r ProductRequest) GetSubCategoryId() uint {
    return r.SubCategoryId
}

func (r ProductRequest) GetId() uint {
    return r.ID
}

// Make ProductFormRequest implement ProductRequestable
// func (r ProductFormRequest) GetName() string {
//     return r.Name
// }

// func (r ProductFormRequest) GetDescription() string {
//     return r.Description
// }

// func (r ProductFormRequest) GetImage() string {
//     return r.Image
// }

type ProductRequest struct {
	Name string `form:"name" validate:"required,max=255"`
	CategoryId uint `form:"category_id" validate:"required"`
	SubCategoryId uint `form:"sub_category_id" validate:"required"`
    Description string `form:"description"`
    Image string `form:"image"`
	ID uint `form:"id"`
}

type ProductIDParamRequest struct {
	ID uint `param:"id" binding:"required"`
}

// type ProductFormRequest struct {
// 	Name string `form:"name" validate:"required,max=255"`
//     Description string `form:"description"`
//     Image string `form:"image"`
// }