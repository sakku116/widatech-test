package ucase

import (
	"backend/domain/dto"
	"backend/domain/model"
	"backend/repository"
	error_utils "backend/utils/error"
	"backend/utils/helper"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type InvoiceUcase struct {
	invoiceRepo repository.IInvoiceRepo
	productRepo repository.IProductRepo
}

type IInvoiceUcase interface {
	CreateInvoice(
		payload dto.CreateInvoiceReq,
	) (*dto.CreateInvoiceResp, error)
	UpdateInvoice(
		invoiceUUID string,
		payload dto.UpdateInvoiceReq,
	) (*dto.UpdateInvoiceRespData, error)
	DeleteInvoice(
		invoiceUUID string,
	) error
	DeleteByInvoiceNo(invoiceNo string) error
	GetInvoiceDetail(invoiceUUID string) (*dto.GetInvoiceDetailRespData, error)
	GetInvoiceList(
		payload dto.GetInvoiceListReq,
	) (*dto.GetInvoiceListRespData, error)
}

func NewInvoiceUcase(
	invoiceRepo repository.IInvoiceRepo,
	productRepo repository.IProductRepo,
) IInvoiceUcase {
	return &InvoiceUcase{
		invoiceRepo: invoiceRepo,
		productRepo: productRepo,
	}
}

func (u *InvoiceUcase) CreateInvoice(
	payload dto.CreateInvoiceReq,
) (*dto.CreateInvoiceResp, error) {
	// validate
	err := payload.Validate()
	if err != nil {
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  "invalid request",
			Detail:   err.Error(),
			Data:     nil,
		}
	}

	// validate products
	var products []model.Product
	if len(payload.ProductUUIDs) > 0 {
		logger.Debugf("getting products by uuids: %v", payload.ProductUUIDs)
		products, err := u.productRepo.GetListByUUIDs(payload.ProductUUIDs)
		if err != nil {
			return nil, &error_utils.CustomErr{
				HttpCode: 500,
				Message:  "internal server error",
				Detail:   err.Error(),
				Data:     nil,
			}
		}

		if len(products) != len(payload.ProductUUIDs) {
			// get the not found products
			var notFoundUUIDs []string
			productsMap := make(map[string]model.Product)
			for _, product := range products {
				productsMap[product.UUID] = product
			}
			for _, uuid := range payload.ProductUUIDs {
				if _, ok := productsMap[uuid]; !ok {
					notFoundUUIDs = append(notFoundUUIDs, uuid)
				}
			}
			return nil, &error_utils.CustomErr{
				HttpCode: 400,
				Message:  "invalid request",
				Detail:   fmt.Sprintf("invalid product uuids, these inputs are not found: %v", notFoundUUIDs),
				Data:     nil,
			}
		}
	}

	// create new invoice
	newInvoice := &model.Invoice{
		Products:        products,
		UUID:            uuid.New().String(),
		InvoiceNo:       strconv.FormatInt(helper.TimeNowEpochUtc(), 10),
		Date:            helper.TimeNowUTC(),
		CustomerName:    payload.CustomerName,
		SalesPersonName: payload.SalesPersonName,
		PaymentType:     payload.PaymentType,
		Notes:           payload.Notes,
	}

	// validate
	err = newInvoice.Validate()
	if err != nil {
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  "invalid request",
			Detail:   err.Error(),
			Data:     nil,
		}
	}

	// insert
	err = u.invoiceRepo.Create(newInvoice)
	if err != nil {
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			Message:  "internal server error",
			Detail:   err.Error(),
			Data:     nil,
		}
	}

	return &dto.CreateInvoiceResp{
		BaseInvoiceResp: newInvoice.ToBaseResp(),
	}, nil
}

func (u *InvoiceUcase) UpdateInvoice(
	invoiceUUID string,
	payload dto.UpdateInvoiceReq,
) (*dto.UpdateInvoiceRespData, error) {
	err := payload.Validate()
	if err != nil {
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  "invalid request",
			Detail:   err.Error(),
			Data:     nil,
		}
	}

	// find invoice
	invoice, err := u.invoiceRepo.GetByUUID(invoiceUUID, true)
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				Message:  "invoice not found",
				Detail:   err.Error(),
				Data:     nil,
			}
		}
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			Message:  "internal server error",
			Detail:   err.Error(),
			Data:     nil,
		}
	}

	// update fields
	if payload.CustomerName != nil {
		invoice.CustomerName = *payload.CustomerName
	}
	if payload.SalesPersonName != nil {
		invoice.SalesPersonName = *payload.SalesPersonName
	}
	if payload.PaymentType != nil {
		invoice.PaymentType = *payload.PaymentType
	}
	if payload.Notes != nil {
		if *payload.Notes == "null" {
			invoice.Notes = nil
		} else {
			invoice.Notes = payload.Notes
		}
	}
	if payload.ProductUUIDs != nil {
		if len(*payload.ProductUUIDs) == 0 {
			invoice.Products = nil
		} else {
			products, err := u.productRepo.GetListByUUIDs(*payload.ProductUUIDs)
			if err != nil {
				return nil, &error_utils.CustomErr{
					HttpCode: 500,
					Message:  "internal server error",
					Detail:   err.Error(),
					Data:     nil,
				}
			}
			invoice.Products = products
		}
	}

	// validate
	err = invoice.Validate()
	if err != nil {
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  "invalid request",
			Detail:   err.Error(),
			Data:     nil,
		}
	}

	// update
	err = u.invoiceRepo.Update(invoice)
	if err != nil {
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			Message:  "internal server error",
			Detail:   err.Error(),
			Data:     nil,
		}
	}

	return &dto.UpdateInvoiceRespData{
		BaseInvoiceResp: invoice.ToBaseResp(),
	}, nil
}

func (u *InvoiceUcase) DeleteByInvoiceNo(invoiceNo string) error {
	err := u.invoiceRepo.DeleteByInvoiceNo(invoiceNo)
	if err != nil {
		if err.Error() == "not found" {
			return &error_utils.CustomErr{
				HttpCode: 404,
				Message:  "invoice not found",
				Detail:   err.Error(),
				Data:     nil,
			}
		}
		return &error_utils.CustomErr{
			HttpCode: 500,
			Message:  "internal server error",
			Detail:   err.Error(),
			Data:     nil,
		}
	}

	return nil
}

func (u *InvoiceUcase) DeleteInvoice(
	invoiceUUID string,
) error {
	// find
	err := u.invoiceRepo.Delete(invoiceUUID)
	if err != nil {
		if err.Error() == "not found" {
			return &error_utils.CustomErr{
				HttpCode: 404,
				Message:  "invoice not found",
				Detail:   err.Error(),
				Data:     nil,
			}
		}
		return &error_utils.CustomErr{
			HttpCode: 500,
			Message:  "internal server error",
			Detail:   err.Error(),
			Data:     nil,
		}
	}

	return nil
}

func (u *InvoiceUcase) GetInvoiceDetail(invoiceUUID string) (*dto.GetInvoiceDetailRespData, error) {
	invoice, err := u.invoiceRepo.GetByUUID(invoiceUUID, true)
	if err != nil {
		if err.Error() == "not found" {
			return nil, &error_utils.CustomErr{
				HttpCode: 404,
				Message:  "invoice not found",
				Detail:   err.Error(),
				Data:     nil,
			}
		}
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			Message:  "internal server error",
			Detail:   err.Error(),
			Data:     nil,
		}
	}

	// response
	resp := &dto.GetInvoiceDetailRespData{
		BaseInvoiceResp: invoice.ToBaseResp(),
	}
	for _, product := range invoice.Products {
		resp.Products = append(resp.Products, product.ToBaseResp())
	}
	return resp, nil
}

func (u *InvoiceUcase) GetInvoiceList(
	payload dto.GetInvoiceListReq,
) (*dto.GetInvoiceListRespData, error) {
	// validate
	err := payload.Validate()
	if err != nil {
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  "invalid request",
			Detail:   err.Error(),
			Data:     nil,
		}
	}

	// parse date range
	var date_gte, date_lte *time.Time
	if payload.DateFrom != nil {
		date_gte, err = helper.ParseDateString(*payload.DateFrom)
		if err != nil {
			return nil, &error_utils.CustomErr{
				HttpCode: 400,
				Message:  "invalid request",
				Detail:   fmt.Sprintf("invalid from date: %v", err.Error()),
				Data:     nil,
			}
		}
	}
	if payload.DateTo != nil {
		date_lte, err = helper.ParseDateString(*payload.DateTo)
		if err != nil {
			return nil, &error_utils.CustomErr{
				HttpCode: 400,
				Message:  "invalid request",
				Detail:   fmt.Sprintf("invalid to date: %v", err.Error()),
				Data:     nil,
			}
		}
	}

	// find
	invoices, count, err := u.invoiceRepo.GetList(
		dto.InvoiceRepo_GetListParams{
			Date_gte:    date_gte,
			Date_lte:    date_lte,
			PaymentType: payload.PaymentType,
			Query:       payload.Query,
			QueryBy:     payload.QueryBy,
			Page:        &payload.Page,
			Limit:       &payload.Limit,
			SortOrder:   &payload.SortOrder,
			SortBy:      &payload.SortBy,
			DoCount:     true,
		},
	)
	if err != nil {
		return nil, &error_utils.CustomErr{
			HttpCode: 500,
			Message:  "internal server error",
			Detail:   err.Error(),
			Data:     nil,
		}
	}

	// resp
	resp := &dto.GetInvoiceListRespData{}
	logger.Debugf("count: %v", count)
	resp.SetPagination(count, payload.Page, payload.Limit)
	for _, invoice := range invoices {
		resp.Data = append(resp.Data, invoice.ToBaseResp())
	}

	return resp, nil
}
