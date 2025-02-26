package service

import (
	"errors"
	"go-electroshop/internal/payload/entity"
	"go-electroshop/internal/payload/request"
	"go-electroshop/internal/payload/response"
	"math"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductService struct {
	DB *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{DB: db}
}

func (s *ProductService) GetProducts(filter request.ProductFilter) (*response.ProductListResponse, error) {
	var products []entity.Product

	query := s.DB.Model(&entity.Product{})

	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR category ILIKE ?", searchTerm, searchTerm)
	}

	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		logrus.Errorf("Failed to count products: %v", err)
		return nil, errors.New("failed to count products")
	}

	// Pagination
	offset := (filter.Page - 1) * filter.Limit
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(filter.Limit).
		Find(&products).
		Error; err != nil {
		logrus.Errorf("Failed to get products: %v", err)
		return nil, errors.New("failed to get products")
	}

	// Transform to response
	productResponse := make([]response.ProductResponse, len(products))
	for i, product := range products {
		productResponse[i] = response.ProductResponse{
			ID:        product.ID,
			Thumbnail: product.Thumbnail,
			Category:  product.Category,
			Name:      product.Name,
			Price:     product.Price,
			ImageLink: product.ImageLink,
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
		}
	}

	return &response.ProductListResponse{
		Products: productResponse,
		Pagination: response.Pagination{
			CurrentPage: filter.Page,
			TotalPage:   int(math.Ceil(float64(total) / float64(filter.Limit))),
			TotalItems:  total,
			ItemPerPage: filter.Limit,
		},
	}, nil
}

func (s *ProductService) GetProductByID(productID uint) (*response.ProductResponse, error) {
	var product entity.Product

	if err := s.DB.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		logrus.Errorf("Error getting product by ID: %v", err)
		return nil, errors.New("failed to get product")
	}

	return &response.ProductResponse{
		ID:        product.ID,
		Thumbnail: product.Thumbnail,
		Category:  product.Category,
		Name:      product.Name,
		Price:     product.Price,
		ImageLink: product.ImageLink,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}, nil
}

func (s *ProductService) CreateProduct(req *request.ProductRequest) (*response.ProductResponse, error) {
	product := entity.Product{
		Thumbnail: req.Thumbnail,
		Category:  strings.TrimSpace(req.Category),
		Name:      strings.TrimSpace(req.Name),
		Price:     req.Price,
		ImageLink: req.ImageLink,
	}

	// Save to database
	if err := s.DB.Create(&product).Error; err != nil {
		logrus.Errorf("Error creating product: %v", err)
		return nil, errors.New("failed to create product")
	}

	return &response.ProductResponse{
		ID:        product.ID,
		Thumbnail: product.Thumbnail,
		Category:  product.Category,
		Name:      product.Name,
		Price:     product.Price,
		ImageLink: product.ImageLink,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}, nil
}

func (s *ProductService) UpdateProduct(productID uint, req *request.UpdateProductRequest) (*response.ProductResponse, error) {
	var product entity.Product

	if err := s.DB.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		logrus.Errorf("Error getting product for update: %v", err)
		return nil, errors.New("failed to get product")
	}

	// Update fields
	product.Thumbnail = req.Thumbnail
	product.Category = strings.TrimSpace(req.Category)
	product.Name = strings.TrimSpace(req.Name)
	product.Price = req.Price
	product.ImageLink = req.ImageLink

	if err := s.DB.Save(&product).Error; err != nil {
		logrus.Errorf("Error updating product: %v", err)
		return nil, errors.New("failed to update product")
	}

	return &response.ProductResponse{
		ID:        product.ID,
		Thumbnail: product.Thumbnail,
		Category:  product.Category,
		Name:      product.Name,
		Price:     product.Price,
		ImageLink: product.ImageLink,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}, nil
}

func (s *ProductService) DeleteProduct(productID uint) error {
	result := s.DB.Delete(&entity.Product{}, productID)

	if result.Error != nil {
		logrus.Errorf("Error deleting product: %v", result.Error)
		return errors.New("failed to delete product")
	}

	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (s *ProductService) GetProductCategories() ([]string, error) {
	var categories []string

	err := s.DB.Model(&entity.Product{}).
		Distinct("category").
		Pluck("category", &categories).Error

	if err != nil {
		logrus.Errorf("Error getting product categories: %v", err)
		return nil, errors.New("failed to get product categories")
	}

	return categories, nil
}
