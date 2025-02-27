package repository

import (
	"errors"
	"go-electroshop/internal/payload/entity"

	"gorm.io/gorm"
)

type CartRepository struct {
	DB *gorm.DB
}

// GetUserCart retrieves cart items for a user
func (r *CartRepository) GetUserCart(userID uint) ([]entity.CartItem, error) {
	var cartItems []entity.CartItem
	err := r.DB.Where("user_id = ?", userID).Preload("Product").Find(&cartItems).Error
	return cartItems, err
}

// AddToCart adds an item to the user's cart
func (r *CartRepository) AddToCart(userID, productID uint, quantity int) error {
	var existingItem entity.CartItem
	err := r.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&existingItem).Error

	if err == nil {
		// Item exists, update quantity
		existingItem.Quantity += quantity
		return r.DB.Save(&existingItem).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Item doesn't exist, create new
		newItem := entity.CartItem{
			UserID:    userID,
			ProductID: productID,
			Quantity:  quantity,
		}
		return r.DB.Create(&newItem).Error
	}

	return err
}

func (r *CartRepository) UpdateCartItemQuantity(itemID, userID uint, quantity int) error {
	// Validasi quantity minimal 1
	if quantity < 1 {
		return errors.New("quantity must be at least 1")
	}

	// Cari item dan pastikan milik user tersebut
	var cartItem entity.CartItem
	if err := r.DB.Where("id = ? AND user_id = ?", itemID, userID).First(&cartItem).Error; err != nil {
		return err
	}

	// Update quantity
	cartItem.Quantity = quantity
	return r.DB.Save(&cartItem).Error
}

// RemoveFromCart removes an item from user's cart
func (r *CartRepository) RemoveFromCart(itemID, userID uint) error {
	result := r.DB.Where("id = ? AND user_id = ?", itemID, userID).Delete(&entity.CartItem{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("cart item not found or doesn't belong to the user")
	}

	return nil
}

// ClearCart removes all items from a user's cart
func (r *CartRepository) ClearCart(userID uint) error {
	return r.DB.Where("user_id = ?", userID).Delete(&entity.CartItem{}).Error
}
