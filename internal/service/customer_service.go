package service

import (
	"bizkit-backend/internal/model"
	"bizkit-backend/internal/repository"
	"errors"
)

type CustomerRequest struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

func GetAllCustomers(search string) ([]model.Customer, error) {
	return repository.GetAllCustomers(search)
}

func GetCustomerByID(id uint) (*model.Customer, error) {
	customer, err := repository.GetCustomerByID(id)
	if err != nil {
		return nil, errors.New("Pelanggan tidak ditemukan")
	}
	return customer, nil
}

func CreateCustomer(req CustomerRequest) (*model.Customer, error) {
	customer := model.Customer{
		Name:    req.Name,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
	}
	err := repository.CreateCustomer(&customer)
	return &customer, err
}

func UpdateCustomer(id uint, req CustomerRequest) (*model.Customer, error) {
	customer, err := repository.GetCustomerByID(id)
	if err != nil {
		return nil, errors.New("Pelanggan tidak ditemukan")
	}

	customer.Name = req.Name
	customer.Phone = req.Phone
	customer.Email = req.Email
	customer.Address = req.Address

	err = repository.UpdateCustomer(customer)
	return customer, err
}

func DeleteCustomer(id uint) error {
	return repository.DeleteCustomer(id)
}
