package dto

import (
	"backend/domain/enum"
	"backend/domain/model"
	"fmt"

	"github.com/google/uuid"
)

type InvoiceRepo_GetListParams struct {
	PaymentType *enum.InvoicePaymentType
	Query       *string
	QueryBy     *string // leave empty to query by all
	Page        *int
	Limit       *int
	SortOrder   *enum.SortOrder
	SortBy      *string
	DoCount     bool
}

func (params *InvoiceRepo_GetListParams) Validate() error {
	if params.PaymentType != nil && !(*params.PaymentType).IsValid() {
		return fmt.Errorf("invalid payment_type")
	}

	if params.SortOrder != nil && !(*params.SortOrder).IsValid() {
		return fmt.Errorf("invalid sort_order")
	}

	tmp := model.Invoice{}
	if params.QueryBy != nil {
		queriableFields := tmp.GetProps().QueriableFields
		contain := false
		for _, field := range queriableFields {
			if *params.QueryBy == field {
				contain = true
				break
			}
		}
		if !contain {
			return fmt.Errorf("invalid query_by")
		}
	}

	if params.SortBy != nil {
		sortableFields := tmp.GetProps().SortableFields
		contain := false
		for _, field := range sortableFields {
			if *params.SortBy == field {
				contain = true
				break
			}
		}
		if !contain {
			return fmt.Errorf("invalid sort_by")
		}
	}

	return nil
}

type CreateInvoiceReq struct {
	CustomerName    string                  `json:"customer_name" binding:"required"`
	SalesPersonName string                  `json:"sales_person_name" binding:"required"`
	PaymentType     enum.InvoicePaymentType `json:"payment_type" binding:"required"`
	Notes           *string                 `json:"notes" binding:"omitempty"`
	ProductUUIDs    []uuid.UUID             `json:"product_uuids" binding:"required"`
}

func (c *CreateInvoiceReq) Validate() error {
	if !c.PaymentType.IsValid() {
		return fmt.Errorf("invalid payment_type")
	}

	return nil
}

type CreateInvoiceResp struct {
	model.BaseInvoiceResp
}

type UpdateInvoiceReq struct {
	CustomerName    *string                  `json:"customer_name" binding:"omitempty"`
	SalesPersonName *string                  `json:"sales_person_name" binding:"omitempty"`
	PaymentType     *enum.InvoicePaymentType `json:"payment_type" binding:"omitempty"`
	Notes           *string                  `json:"notes" binding:"omitempty"`         // use 'null' to set explicitly to null
	ProductUUIDs    *[]uuid.UUID             `json:"product_uuids" binding:"omitempty"` // use 'null' remove all products
}

type GetInvoiceByUUIDRespData struct {
	model.BaseInvoiceResp
}

type GetInvoiceListReq struct {
	PaymentType *enum.InvoicePaymentType `form:"payment_type" binding:"omitempty,oneof=CASH CREDIT"`
	Query       *string                  `form:"query" binding:"omitempty"`
	QueryBy     *string                  `form:"query_by" binding:"omitempty,oneof=invoice_no customer_name sales_person_name"` // leave empty to query by all
	Page        int                      `form:"page" binding:"required" default:"1"`
	Limit       int                      `form:"limit" binding:"required" default:"10"`
	SortOrder   enum.SortOrder           `form:"sort_order" binding:"required,oneof=asc desc" default:"desc"`
	SortBy      string                   `form:"sort_by" binding:"required,oneof=created_at updated_at invoice_no" default:"updated_at"`
}

func (r *GetInvoiceListReq) Validate() error {
	if r.PaymentType != nil && !(*r.PaymentType).IsValid() {
		return fmt.Errorf("invalid payment_type")
	}

	if !r.SortOrder.IsValid() {
		return fmt.Errorf("invalid sort_order")
	}

	return nil
}

type GetInvoiceListRespData struct {
	BasePaginationRespData
	Data []model.BaseInvoiceResp `json:"data"`
}
