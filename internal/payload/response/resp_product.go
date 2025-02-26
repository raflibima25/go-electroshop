package response

import "time"

type ProductResponse struct {
	ID        uint      `json:"id"`
	Thumbnail string    `json:"thumbnail"`
	Category  string    `json:"category"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	ImageLink string    `json:"image_link"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	Pagination Pagination        `json:"pagination"`
}
