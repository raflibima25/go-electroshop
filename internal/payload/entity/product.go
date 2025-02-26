package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Thumbnail string  `gorm:"type:varchar(255)"`
	Category  string  `gorm:"type:varchar(100);not null;index"`
	Name      string  `gorm:"type:varchar(255);not null"`
	Price     float64 `gorm:"type:decimal(15,2);not null"`
	ImageLink string  `gorm:"type:varchar(255)"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.Name == "" {
		return gorm.ErrModelValueRequired
	}
	if p.Price <= 0 {
		return gorm.ErrInvalidData
	}
	return nil
}
