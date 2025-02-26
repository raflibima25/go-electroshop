package controller

import (
	"go-electroshop/internal/payload/response"
	"go-electroshop/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AdminDashboardController struct {
	ProductService *service.ProductService
	UserService    *service.UserService
}

func NewAdminDashboardController(productService *service.ProductService, userService *service.UserService) *AdminDashboardController {
	return &AdminDashboardController{
		ProductService: productService,
		UserService:    userService,
	}
}

// GetDashboardStatsHandler godoc
// @Summary 	Get admin dashboard statistics
// @Description Get statistics for admin dashboard including product counts and categories
// @Tags 		admin
// @Accept 		json
// @Produce 	json
// @Security 	BearerAuth
// @Success 	200 {object} response.SuccessResponse
// @Failure 	401 {object} response.SuccessResponse
// @Router 		/admin/dashboard/stats [get]
func (c *AdminDashboardController) GetDashboardStatsHandler(ctx *gin.Context) {
	// Implementation with mocked data for now
	// TODO: In a real implementation, these would come from the service layer

	categories, err := c.ProductService.GetProductCategories()
	if err != nil {
		logrus.Errorf("Error getting product categories: %v", err)
		categories = []string{}
	}

	// Mock data for stats
	stats := gin.H{
		"total_products": 3,
		"categories":     categories,
		"category_counts": []gin.H{
			{"name": "Iphone", "count": 1},
			{"name": "Samsung", "count": 1},
			{"name": "Xiaomi", "count": 1},
		},
		"recent_products": []gin.H{
			{"id": 1, "name": "Iphone 13 Pro", "price": 12000000, "category": "Iphone"},
			{"id": 2, "name": "Samsung X flip", "price": 20000000, "category": "Samsung"},
			{"id": 3, "name": "Xiaomi Redmi Note 11 Pro", "price": 3200000, "category": "Xiaomi"},
		},
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Dashboard stats retrieved successfully",
		Data:            stats,
	})

}
