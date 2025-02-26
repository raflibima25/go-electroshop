package controller

import (
	"go-electroshop/internal/payload/request"
	"go-electroshop/internal/payload/response"
	"go-electroshop/internal/service"
	"go-electroshop/internal/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ProductController struct {
	ProductService *service.ProductService
}

func NewProductController(productService *service.ProductService) *ProductController {
	return &ProductController{ProductService: productService}
}

// GetProductsHandler godoc
// @Summary 	Get all of products
// @Description Get all of products with filter by category, min price, max price, page, and limit
// @Tags 		product
// @Accept 		json
// @Produce 	json
// @Param 		category query string false "Filter by category"
// @Param 		search query string false "Search term"
// @Param 		min_price query number false "Minimum price"
// @Param 		max_price query number false "Maximum price"
// @Param 		page query int false "Page number"
// @Param 		limit query int false "Items per page"
// @Success 	200 {object} response.SuccessResponse{data=response.ProductListResponse}
// @Failure 	400 {object} response.SuccessResponse
// @Router 		/product [get]
func (c *ProductController) GetProductsHandler(ctx *gin.Context) {
	var filter request.ProductFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		logrus.Errorf("Error binding query params: %v", err)
		utility.ErrorResponse(ctx, 400, "Invalid filter parameters: "+err.Error(), nil)
		return
	}

	// Set default values
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	products, err := c.ProductService.GetProducts(filter)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Get products success",
		Data:            products,
	})
}

// GetProductByIDHandler godoc
// @Summary 	Get product by ID
// @Description Get a specific product by ID
// @Tags 		products
// @Accept 		json
// @Produce 	json
// @Param 		id path int true "Product ID"
// @Success 	200 {object} response.SuccessResponse{data=response.ProductResponse}
// @Failure 	400 {object} response.SuccessResponse
// @Failure 	404 {object} response.SuccessResponse
// @Router 		/product/{id} [get]
func (c *ProductController) GetProductByIDHandler(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	product, err := c.ProductService.GetProductByID(uint(id))
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Get product successful",
		Data:            product,
	})
}

// CreateProductHandler godoc
// @Summary 	Create product
// @Description Create a new product
// @Tags 		products
// @Accept 		json
// @Produce 	json
// @Param 		request body request.ProductRequest true "Product data"
// @Success 	201 {object} response.SuccessResponse{data=response.ProductResponse}
// @Failure 	400 {object} response.SuccessResponse
// @Router 		/product [post]
func (c *ProductController) CreateProductHandler(ctx *gin.Context) {
	var req request.ProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utility.FormatValidationError(validationErrors)
			message := make([]string, len(formattedErrors))
			for i, err := range formattedErrors {
				message[i] = utility.GetReadableErrorMessage(err)
			}
			utility.ErrorResponse(ctx, http.StatusBadRequest, message[0], []response.ErrorDetail{
				{
					Field:   "validation",
					Message: message,
				},
			})
			return
		}
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid input format", nil)
		return
	}

	product, err := c.ProductService.CreateProduct(&req)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusCreated, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Product created successfully",
		Data:            product,
	})
}

// UpdateProductHandler godoc
// @Summary 	Update product
// @Description Update an existing product
// @Tags 		products
// @Accept 		json
// @Produce 	json
// @Param 		id path int true "Product ID"
// @Param 		request body request.UpdateProductRequest true "Product data"
// @Success 	200 {object} response.SuccessResponse{data=response.ProductResponse}
// @Failure 	400 {object} response.SuccessResponse
// @Failure 	404 {object} response.SuccessResponse
// @Router 		/product/{id} [put]
func (c *ProductController) UpdateProductHandler(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	var req request.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utility.FormatValidationError(validationErrors)
			messages := make([]string, len(formattedErrors))
			for i, err := range formattedErrors {
				messages[i] = utility.GetReadableErrorMessage(err)
			}
			utility.ErrorResponse(ctx, http.StatusBadRequest, messages[0], []response.ErrorDetail{
				{
					Field:   "validation",
					Message: messages,
				},
			})
			return
		}
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid input format", nil)
		return
	}

	products, err := c.ProductService.UpdateProduct(uint(id), &req)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Product updated successfully",
		Data:            products,
	})
}

// DeleteProductHandler godoc
// @Summary 	Delete product
// @Description Delete a product by ID
// @Tags 		products
// @Accept 		json
// @Produce 	json
// @Param 		id path int true "Product ID"
// @Success 	200 {object} response.SuccessResponse
// @Failure 	400 {object} response.SuccessResponse
// @Failure 	404 {object} response.SuccessResponse
// @Router 		/product/{id} [delete]
func (c *ProductController) DeleteProductHandler(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	if err := c.ProductService.DeleteProduct(uint(id)); err != nil {
		utility.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Product deleted successfully",
		Data:            nil,
	})
}

// GetProductCategoriesHandler godoc
// @Summary 	Get all product categories
// @Description Get a list of all unique product categories
// @Tags 		products
// @Accept 		json
// @Produce 	json
// @Success 	200 {object} response.SuccessResponse
// @Failure 	400 {object} response.SuccessResponse
// @Router 		/product/categories [get]
func (c *ProductController) GetProductCategoriesHandler(ctx *gin.Context) {
	categories, err := c.ProductService.GetProductCategories()
	if err != nil {
		utility.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		ResponseStatus:  true,
		ResponseMessage: "Get product categories successful",
		Data: gin.H{
			"categories": categories,
		},
	})
}
