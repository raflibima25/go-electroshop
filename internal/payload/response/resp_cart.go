package response

// CartItemResponse represents a single item in the cart
type CartItemResponse struct {
	ID       uint            `json:"id"`
	Quantity int             `json:"quantity"`
	Product  ProductResponse `json:"product"`
}

// CartResponse represents the cart with multiple items
type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalItems int                `json:"total_items"`
	TotalPrice float64            `json:"total_price"`
}
