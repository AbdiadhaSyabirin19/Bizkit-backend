package service

import (
	"errors"

	"bizkit-backend/internal/model"
	"bizkit-backend/internal/repository"
)

type ProductRequest struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryID  *uint   `json:"category_id"`
	BrandID     *uint   `json:"brand_id"`
	UnitID      *uint   `json:"unit_id"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Image       string  `json:"image"`
	Status      string  `json:"status"`
	VariantIDs  []uint  `json:"variant_ids"`
	OutletIDs   []uint  `json:"outlet_ids"`
}

func GetAllProducts(search string) ([]model.Product, error) {
	return repository.GetAllProducts(search)
}

func GetProductByID(id uint) (*model.Product, error) {
	product, err := repository.GetProductByID(id)
	if err != nil {
		return nil, errors.New("Produk tidak ditemukan")
	}
	return product, nil
}

func CreateProduct(req ProductRequest) (*model.Product, error) {
	if req.Status == "" {
		req.Status = "active"
	}
	product := model.Product{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		BrandID:     req.BrandID,
		UnitID:      req.UnitID,
		Price:       req.Price,
		Stock:       req.Stock,
		Image:       req.Image,
		Status:      req.Status,
	}
	err := repository.CreateProduct(&product, req.VariantIDs, req.OutletIDs)
	if err != nil {
		return nil, err
	}
	return repository.GetProductByID(product.ID)
}

func UpdateProduct(id uint, req ProductRequest) (*model.Product, error) {
	product, err := repository.GetProductByID(id)
	if err != nil {
		return nil, errors.New("Produk tidak ditemukan")
	}
	product.Code = req.Code
	product.Name = req.Name
	product.Description = req.Description
	product.CategoryID = req.CategoryID
	product.BrandID = req.BrandID
	product.UnitID = req.UnitID
	product.Price = req.Price
	product.Stock = req.Stock
	if req.Image != "" {
		product.Image = req.Image
	}
	if req.Status != "" {
		product.Status = req.Status
	}
	err = repository.UpdateProduct(product, req.VariantIDs, req.OutletIDs)
	if err != nil {
		return nil, err
	}
	return repository.GetProductByID(product.ID)
}

func DeleteProduct(id uint) error {
	_, err := repository.GetProductByID(id)
	if err != nil {
		return errors.New("Produk tidak ditemukan")
	}
	return repository.DeleteProduct(id)
}

// GetProductPrices — ambil harga per kategori untuk 1 produk
func GetProductPrices(productID uint) ([]model.ProductPrice, error) {
	return repository.GetProductPricesByProductID(productID)
}