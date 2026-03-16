package model

import "gorm.io/gorm"

type Sale struct {
	gorm.Model
	InvoiceNumber   string        `json:"invoice_number" gorm:"unique;not null"`
	UserID          uint          `json:"user_id"`
	CustomerName    string        `json:"customer_name"`
	PaymentMethodID uint          `json:"payment_method_id"`
	PriceCategoryID *uint         `json:"price_category_id"`
	PromoID         *uint         `json:"promo_id"`
	Subtotal        float64       `json:"subtotal"`
	DiscountTotal   float64       `json:"discount_total"`
	ManualDiscount  float64       `json:"manual_discount" gorm:"default:0"`
	AdditionalFee   float64       `json:"additional_fee" gorm:"default:0"`
	GrandTotal      float64       `json:"grand_total"`
	Source          string        `json:"source" gorm:"default:dashboard"`
	User            User          `json:"user" gorm:"foreignKey:UserID"`
	PaymentMethod   PaymentMethod `json:"payment_method" gorm:"foreignKey:PaymentMethodID"`
	PriceCategory   *PriceCategory `json:"price_category" gorm:"foreignKey:PriceCategoryID"`
	Promo           *Promo        `json:"promo" gorm:"foreignKey:PromoID"`
	Items           []SaleItem    `json:"items" gorm:"foreignKey:SaleID"`
}

type SaleItem struct {
	gorm.Model
	SaleID    uint              `json:"sale_id"`
	ProductID uint              `json:"product_id"`
	Quantity  int               `json:"quantity"`
	BasePrice float64           `json:"base_price"`
	Discount  float64           `json:"discount" gorm:"default:0"`
	Subtotal  float64           `json:"subtotal"`
	Product   Product           `json:"product" gorm:"foreignKey:ProductID"`
	Variants  []SaleItemVariant `json:"variants" gorm:"foreignKey:SaleItemID"`
}

type SaleItemVariant struct {
	gorm.Model
	SaleItemID      uint          `json:"sale_item_id"`
	VariantOptionID uint          `json:"variant_option_id"`
	AdditionalPrice float64       `json:"additional_price"`
	VariantOption   VariantOption `json:"variant_option" gorm:"foreignKey:VariantOptionID"`
}