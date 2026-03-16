package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Code        string            `json:"code"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	CategoryID  *uint             `json:"category_id"`
	BrandID     *uint             `json:"brand_id"`
	UnitID      *uint             `json:"unit_id"`
	Price       float64           `json:"price"`
	Stock       int               `json:"stock"`
	Image       string            `json:"image"`
	Status      string            `json:"status"`
	Category    *Category         `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Brand       *Brand            `json:"brand,omitempty" gorm:"foreignKey:BrandID"`
	Unit        *Unit             `json:"unit,omitempty" gorm:"foreignKey:UnitID"`
	Variants    []VariantCategory `json:"variants,omitempty" gorm:"many2many:product_variant_categories;"`
	Outlets     []Outlet          `json:"outlets,omitempty" gorm:"many2many:product_outlets;"`
	Prices      []ProductPrice    `json:"prices,omitempty" gorm:"foreignKey:ProductID"`
}