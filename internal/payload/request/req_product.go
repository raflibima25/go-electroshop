package request

type ProductRequest struct {
	Thumbnail string  `json:"thumbnail"`
	Category  string  `json:"category" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Price     float64 `json:"price" binding:"required,gt=0"`
	ImageLink string  `json:"image_link"`
}

type UpdateProductRequest struct {
	Thumbnail string  `json:"thumbnail"`
	Category  string  `json:"category" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Price     float64 `json:"price" binding:"required,gt=0"`
	ImageLink string  `json:"image_link"`
}

type ProductFilter struct {
	Category string  `form:"category"`
	Search   string  `form:"search"`
	MinPrice float64 `form:"min_price"`
	MaxPrice float64 `form:"max_price"`
	Page     int     `form:"page,default=1"`
	Limit    int     `form:"limit,default=10"`
}
