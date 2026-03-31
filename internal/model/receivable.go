package model

import (
	"time"

	"gorm.io/gorm"
)

type ReceivablePayment struct {
	gorm.Model
	PaymentDate     time.Time               `json:"payment_date"`
	CustomerName    string                  `json:"customer_name"`
	PaymentMethodID uint                    `json:"payment_method_id"`
	Notes           string                  `json:"notes"`
	Amount          float64                 `json:"amount"`
	UserID          uint                    `json:"user_id"`
	User            User                    `json:"user" gorm:"foreignKey:UserID"`
	PaymentMethod   PaymentMethod           `json:"payment_method" gorm:"foreignKey:PaymentMethodID"`
	Items           []ReceivablePaymentItem `json:"items" gorm:"foreignKey:ReceivablePaymentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ReceivablePaymentItem struct {
	gorm.Model
	ReceivablePaymentID uint    `json:"receivable_payment_id"`
	SaleID              uint    `json:"sale_id"`
	AmountPaid          float64 `json:"amount_paid"`

	Sale              *Sale              `json:"sale,omitempty" gorm:"foreignKey:SaleID"`
	ReceivablePayment *ReceivablePayment `json:"receivable_payment,omitempty" gorm:"foreignKey:ReceivablePaymentID"`
}
