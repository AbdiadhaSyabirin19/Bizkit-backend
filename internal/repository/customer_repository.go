package repository

import (
	"bizkit-backend/config"
	"bizkit-backend/internal/model"
)

func GetAllCustomers(search string) ([]model.Customer, error) {
	var customers []model.Customer
	query := config.DB.Model(&model.Customer{})
	if search != "" {
		query = query.Where("name LIKE ? OR phone LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	return customers, query.Find(&customers).Error
}

func GetCustomerByID(id uint) (*model.Customer, error) {
	var customer model.Customer
	err := config.DB.First(&customer, id).Error
	return &customer, err
}

func CreateCustomer(customer *model.Customer) error {
	return config.DB.Create(customer).Error
}

func UpdateCustomer(customer *model.Customer) error {
	return config.DB.Save(customer).Error
}

func DeleteCustomer(id uint) error {
	return config.DB.Delete(&model.Customer{}, id).Error
}
