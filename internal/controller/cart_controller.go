package controller

import (
	"go-electroshop/internal/payload/entity"
	"go-electroshop/internal/payload/request"
	"go-electroshop/internal/payload/response"
	"go-electroshop/internal/repository"
	"go-electroshop/internal/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type cartController struct {
	db       *gorm.DB
	cartRepo *repository.CartRepository
}

func NewCartController(db *gorm.DB, cartRepo *repository.CartRepository) *cartController {
	return &cartController{db: db, cartRepo: cartRepo}
}

// GetCartHandler godoc
// @Summary     Get user's cart
// @Description Retrieves all items in the user's cart
// @Tags        cart
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} response.SuccessResponse{data=response.CartResponse}
// @Failure     401 {object} response.ErrorResponse
// @Router      /cart [get]
func (c *cartController) GetCartHandler(ctx *gin.Context) {
	// Dapatkan user ID dari context (set oleh middleware auth)
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			ResponseStatus:  false,
			ResponseMessage: "Unauthorized: " + err.Error(),
		})
		return
	}

	// Ambil cart items
	cartItems, err := c.cartRepo.GetUserCart(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			ResponseStatus:  false,
			ResponseMessage: "Failed to get cart: " + err.Error(),
		})
		return
	}

	// Transform to response format
	cartResponse := response.CartResponse{
		Items:      make([]response.CartItemResponse, 0, len(cartItems)),
		TotalItems: 0,
		TotalPrice: 0,
	}

	for _, item := range cartItems {
		// Calculate totals
		cartResponse.TotalItems += item.Quantity
		cartResponse.TotalPrice += item.Product.Price * float64(item.Quantity)

		// Add to items
		cartResponse.Items = append(cartResponse.Items, response.CartItemResponse{
			ID:       item.ID,
			Quantity: item.Quantity,
			Product: response.ProductResponse{
				ID:        item.Product.ID,
				Thumbnail: item.Product.Thumbnail,
				Category:  item.Product.Category,
				Name:      item.Product.Name,
				Price:     item.Product.Price,
				ImageLink: item.Product.ImageLink,
				CreatedAt: item.Product.CreatedAt,
				UpdatedAt: item.Product.UpdatedAt,
			},
		})
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Cart retrieved successfully",
		Data:            cartResponse,
	})
}

// AddToCartHandler godoc
// @Summary     Add to cart
// @Description Add a product to the user's cart
// @Tags        cart
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body AddToCartRequest true "Add to cart request"
// @Success     200 {object} response.SuccessResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Router      /cart [post]
func (c *cartController) AddToCartHandler(ctx *gin.Context) {
	// Dapatkan user ID dari context
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			ResponseStatus:  false,
			ResponseMessage: "Unauthorized: " + err.Error(),
		})
		return
	}

	// Parse request body
	var req request.AddToCartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			ResponseStatus:  false,
			ResponseMessage: "Invalid request: " + err.Error(),
		})
		return
	}

	// Validasi product ada
	var product entity.Product
	if err := c.db.First(&product, req.ProductID).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			ResponseStatus:  false,
			ResponseMessage: "Product not found",
		})
		return
	}

	// Tambahkan ke cart
	if err := c.cartRepo.AddToCart(userID, req.ProductID, req.Quantity); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			ResponseStatus:  false,
			ResponseMessage: "Failed to add to cart: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Item added to cart successfully",
	})
}

// UpdateCartItemHandler godoc
// @Summary     Update cart item
// @Description Update quantity of a cart item
// @Tags        cart
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Cart Item ID"
// @Param       request body UpdateCartItemRequest true "Update cart item request"
// @Success     200 {object} response.SuccessResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Router      /cart/{id} [put]
func (c *cartController) UpdateCartItemHandler(ctx *gin.Context) {
	// Dapatkan user ID dari context
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			ResponseStatus:  false,
			ResponseMessage: "Unauthorized: " + err.Error(),
		})
		return
	}

	// Parse item ID dari path
	itemID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			ResponseStatus:  false,
			ResponseMessage: "Invalid item ID",
		})
		return
	}

	// Parse request body
	var req request.UpdateCartItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			ResponseStatus:  false,
			ResponseMessage: "Invalid request: " + err.Error(),
		})
		return
	}

	// Update item
	if err := c.cartRepo.UpdateCartItemQuantity(uint(itemID), userID, req.Quantity); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			ResponseStatus:  false,
			ResponseMessage: "Failed to update cart item: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Cart item updated successfully",
	})
}

// RemoveFromCartHandler godoc
// @Summary     Remove from cart
// @Description Remove an item from the user's cart
// @Tags        cart
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Cart Item ID"
// @Success     200 {object} response.SuccessResponse
// @Failure     400 {object} response.ErrorResponse
// @Failure     401 {object} response.ErrorResponse
// @Router      /cart/{id} [delete]
func (c *cartController) RemoveFromCartHandler(ctx *gin.Context) {
	// Dapatkan user ID dari context
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			ResponseStatus:  false,
			ResponseMessage: "Unauthorized: " + err.Error(),
		})
		return
	}

	// Parse item ID dari path
	itemID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			ResponseStatus:  false,
			ResponseMessage: "Invalid item ID",
		})
		return
	}

	// Hapus item dari cart
	if err := c.cartRepo.RemoveFromCart(uint(itemID), userID); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			ResponseStatus:  false,
			ResponseMessage: "Failed to remove from cart: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Item removed from cart successfully",
	})
}

// ClearCartHandler godoc
// @Summary     Clear cart
// @Description Remove all items from the user's cart
// @Tags        cart
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} response.SuccessResponse
// @Failure     401 {object} response.ErrorResponse
// @Router      /cart [delete]
func (c *cartController) ClearCartHandler(ctx *gin.Context) {
	// Dapatkan user ID dari context
	userID, err := utility.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			ResponseStatus:  false,
			ResponseMessage: "Unauthorized: " + err.Error(),
		})
		return
	}

	// Kosongkan cart
	if err := c.cartRepo.ClearCart(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			ResponseStatus:  false,
			ResponseMessage: "Failed to clear cart: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Cart cleared successfully",
	})
}
