package model

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	Name    string `json:"name" gorm:"not null"`
	Phone   string `json:"phone" gorm:"uniqueIndex"`
	Email   string `json:"email"`
	Address string `json:"address"`
}
