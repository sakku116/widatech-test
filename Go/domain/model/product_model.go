package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	UUID                 uuid.UUID `gorm:"type:uuid;not null" json:"uuid"`
	InvoiceID            uint      `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	InvoiceUUID          uuid.UUID `gorm:"type:uuid;not null;index" json:"invoice_uuid"`
	ItemName             string    `gorm:"type:varchar(100);not null" json:"item_name"`
	Quantity             int       `gorm:"not null" json:"quantity"`
	TotalCostOfGoodsSold int64     `gorm:"not null" json:"total_cost_of_goods_sold"`
	TotalPriceSold       int64     `gorm:"not null" json:"total_price_sold"`
}

func (i *Product) GetProps() ModelProps {
	return ModelProps{
		QueriableFields: []string{
			"item_name",
		},
		SortableFields: []string{
			"created_at",
			"updated_at",
			"quantity",
		},
	}
}

func (i *Product) Validate() error {
	if len(i.ItemName) < 5 {
		return errors.New("invalid item_name length")
	}

	if i.Quantity < 1 {
		return errors.New("invalid quantity")
	}

	return nil
}

type BaseProductResp struct {
	UUID                 uuid.UUID `json:"uuid"`
	InvoiceUUID          uuid.UUID `json:"invoice_uuid"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	ItemName             string    `json:"item_name"`
	Quantity             int       `json:"quantity"`
	TotalCostOfGoodsSold int64     `json:"total_cost_of_goods_sold"`
	TotalPriceSold       int64     `json:"total_price_sold"`
}

func (i *Product) ToBaseResp() BaseProductResp {
	return BaseProductResp{
		UUID:                 i.UUID,
		InvoiceUUID:          i.InvoiceUUID,
		CreatedAt:            i.CreatedAt,
		UpdatedAt:            i.UpdatedAt,
		ItemName:             i.ItemName,
		Quantity:             i.Quantity,
		TotalCostOfGoodsSold: i.TotalCostOfGoodsSold,
		TotalPriceSold:       i.TotalPriceSold,
	}
}
