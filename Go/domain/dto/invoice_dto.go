package dto

import (
	"backend/domain/enum"
	"backend/domain/model"
	"fmt"
	"time"
)

type InvoiceRepo_GetListParams struct {
	CreatedAt_gte   *time.Time
	CreatedAt_lte   *time.Time
	Date_gte        *time.Time
	Date_lte        *time.Time
	PaymentType     *enum.InvoicePaymentType
	Query           *string
	QueryBy         *string // leave empty to query by all
	Page            *int
	Limit           *int
	SortOrder       *enum.SortOrder
	SortBy          *string
	PreloadProducts bool
	DoCount         bool
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
	ProductUUIDs    []string                `json:"product_uuids" binding:"required"`
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
	ProductUUIDs    *[]string                `json:"product_uuids" binding:"omitempty"` // use empty array [] remove all products
}

func (r *UpdateInvoiceReq) Validate() error {
	if r.PaymentType != nil && !(*r.PaymentType).IsValid() {
		return fmt.Errorf("invalid payment_type")
	}
	return nil
}

type UpdateInvoiceRespData struct {
	model.BaseInvoiceResp
}

type GetInvoiceByUUIDRespData struct {
	model.BaseInvoiceResp
}

type GetInvoiceListReq struct {
	DateFrom    *string                  `form:"date_from" binding:"omitempty"`                      // DD-MM-YYYY, fill up date_from and date_to to get profit_total and cash_transaction_total. WARNING: if set, pagination will be ignored for profit calculation.
	DateTo      *string                  `form:"date_to" binding:"omitempty"`                        // DD-MM-YYYY, fill up date_from and date_to to get profit_total and cash_transaction_total. WARNING: if set, pagination will be ignored for profit calculation.
	PaymentType *enum.InvoicePaymentType `form:"payment_type" binding:"omitempty,oneof=CASH CREDIT"` // leave empty to query all payment types
	Query       *string                  `form:"query" binding:"omitempty"`
	QueryBy     *string                  `form:"query_by" binding:"omitempty,oneof=invoice_no customer_name sales_person_name"` // leave empty to query by all queriable fields
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

type GetInvoiceListRespData_DataItem struct {
	model.BaseInvoiceResp
	ProductTotal int `json:"product_total"`
}

type GetInvoiceListRespData struct {
	BasePaginationRespData
	ProfitTotal          int64                             `json:"profit_total"`
	CashTransactionTotal int64                             `json:"cash_transaction_total"`
	Data                 []GetInvoiceListRespData_DataItem `json:"data"`
}

type GetInvoiceDetailRespData struct {
	model.BaseInvoiceResp
	Products []model.BaseProductResp `json:"products"`
}
