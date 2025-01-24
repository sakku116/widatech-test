package ucase

import (
	"backend/domain/dto"
	"backend/domain/model"
	"backend/repository"
	error_utils "backend/utils/error"
	"backend/utils/helper"
	"strconv"

	"github.com/google/uuid"
)

type InvoiceUcase struct {
	invoiceRepo repository.IInvoiceRepo
	productRepo repository.IProductRepo
}

type IInvoiceUcase interface {
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

func (u *InvoiceUcase) CreateInvoice(payload dto.CreateInvoiceReq) (*dto.CreateInvoiceResp, error) {
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
		return nil, &error_utils.CustomErr{
			HttpCode: 400,
			Message:  "invalid request",
			Detail:   "invalid product uuids",
			Data:     nil,
		}
	}

	// create new invoice
	newInvoice := &model.Invoice{
		Products:        products,
		UUID:            uuid.New(),
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
