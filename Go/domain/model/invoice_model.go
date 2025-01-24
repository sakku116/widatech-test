package model

import (
	"backend/domain/enum"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	UUID            uuid.UUID               `gorm:"type:uuid;not null" json:"uuid"`
	InvoiceNo       string                  `gorm:"type:varchar(50);not null;unique" json:"invoice_no"`
	Date            time.Time               `gorm:"type:timestamp;not null" json:"date"`
	CustomerName    string                  `gorm:"type:varchar(100);not null" json:"customer_name"`
	SalesPersonName string                  `gorm:"type:varchar(100);not null" json:"sales_person_name"`
	PaymentType     enum.InvoicePaymentType `gorm:"type:varchar(50);not null" json:"payment_type"`
	Notes           *string                 `gorm:"type:text" json:"notes"`
	Products        []Product               `gorm:"foreignKey:InvoiceID;references:ID" json:"products"`
}

func (i *Invoice) GetProps() ModelProps {
	return ModelProps{
		QueriableFields: []string{
			"invoice_no",
			"customer_name",
			"sales_person_name",
		},
		SortableFields: []string{
			"created_at",
			"updated_at",
			"invoice_no",
		},
	}
}

func (i *Invoice) Validate() error {
	if len(i.InvoiceNo) < 1 {
		return errors.New("invoice_no cannot be empty")
	}

	if len(i.CustomerName) < 2 {
		return errors.New("invalid customer_name length")
	}

	if len(i.SalesPersonName) < 2 {
		return errors.New("invalid sales_person_name length")
	}

	if !i.PaymentType.IsValid() {
		return errors.New("invalid payment_type")
	}

	if i.Notes != nil && len(*i.Notes) < 5 {
		return errors.New("invalid notes length")
	}

	return nil
}

type BaseInvoiceResp struct {
	UUID            uuid.UUID               `json:"uuid"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
	InvoiceNo       string                  `json:"invoice_no"`
	Date            string                  `json:"date"`
	CustomerName    string                  `json:"customer_name"`
	SalesPersonName string                  `json:"sales_person_name"`
	PaymentType     enum.InvoicePaymentType `json:"payment_type"`
	Notes           *string                 `json:"notes"`
}

func (i *Invoice) ToBaseResp() BaseInvoiceResp {
	return BaseInvoiceResp{
		UUID:            i.UUID,
		CreatedAt:       i.CreatedAt,
		UpdatedAt:       i.UpdatedAt,
		InvoiceNo:       i.InvoiceNo,
		Date:            i.Date.Format(time.RFC3339),
		CustomerName:    i.CustomerName,
		SalesPersonName: i.SalesPersonName,
		PaymentType:     i.PaymentType,
		Notes:           i.Notes,
	}
}
