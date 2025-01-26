package ucase

import (
	"backend/domain/dto"
	"backend/domain/enum"
	"backend/domain/model"
	"backend/repository"
	error_utils "backend/utils/error"
	"backend/utils/helper"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
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
	ImportXlsx(filepath string) error
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

			if len(products) != len(*payload.ProductUUIDs) {
				// get the not found products
				var notFoundUUIDs []string
				productsMap := make(map[string]model.Product)
				for _, product := range products {
					productsMap[product.UUID] = product
				}
				for _, uuid := range *payload.ProductUUIDs {
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
	logger.Debugf("payload: %v", payload)
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

	// ignore pagination when date_from and date_to is set (for profit calculation)
	var page, limit *int
	if payload.DateFrom != nil && payload.DateTo != nil {
		logger.Debugf("pagination ignored because date_from is set")
		page = nil
		limit = nil
	} else {
		page = &payload.Page
		limit = &payload.Limit
	}

	// find
	invoices, count, err := u.invoiceRepo.GetList(
		dto.InvoiceRepo_GetListParams{
			Date_gte:        date_gte,
			Date_lte:        date_lte,
			PaymentType:     payload.PaymentType,
			Query:           payload.Query,
			QueryBy:         payload.QueryBy,
			Page:            page,
			Limit:           limit,
			SortOrder:       &payload.SortOrder,
			SortBy:          &payload.SortBy,
			DoCount:         true,
			PreloadProducts: true,
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
	resp.SetPagination(count, payload.Page, payload.Limit)
	for _, invoice := range invoices {
		resp.Data = append(resp.Data, dto.GetInvoiceListRespData_DataItem{
			BaseInvoiceResp: invoice.ToBaseResp(),
			ProductTotal:    len(invoice.Products),
		})
	}

	// calculate profit_total and cash_transaction_total
	if payload.DateFrom != nil && payload.DateTo != nil {
		logger.Debugf("calculating profit_total and cash_transaction_total")
		var profitTotal, cashTransactionTotal int64
		for _, invoice := range invoices {
			for _, product := range invoice.Products {
				profitTotal += product.TotalPriceSold - product.TotalCostOfGoodsSold
				cashTransactionTotal += product.TotalPriceSold + product.TotalCostOfGoodsSold
			}
		}
		resp.ProfitTotal = profitTotal
		resp.CashTransactionTotal = cashTransactionTotal
	}

	return resp, nil
}

func (u *InvoiceUcase) ImportXlsx(filepath string) error {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	// search for product sheet
	productStatusesByInvoiceNo := make(map[string][]map[string]interface{})
	/*
		{
			[invoice_no]: [
				{
					"row_index": [row_index],
					"data":      &model.Product{},
					"error":     "",
				}
			]
		}
	*/

	for _, sheetName := range f.GetSheetList() {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return err
		}
		if strings.Contains(strings.ToLower(sheetName), "product") {
			for i, columns := range rows {
				if i == 0 { // header
					// ensure columns total
					if len(columns) != 5 {
						return &error_utils.CustomErr{
							HttpCode: 400,
							Message:  "invalid xlsx file",
							Detail:   "columns total for product sheet must be 5",
							Data:     nil,
						}
					}

					// skip header
					continue
				}

				logger.Debugf("product columns: %v", columns)

				invoiceNo := columns[0]
				if invoiceNo == "" {
					invoiceNo = "unknown"
				}

				// prepare status
				if invoiceProductStatues, ok := productStatusesByInvoiceNo[invoiceNo]; !ok || len(invoiceProductStatues) == 0 {
					productStatusesByInvoiceNo[invoiceNo] = make([]map[string]interface{}, 0)
				}
				status := map[string]interface{}{
					"row_index": i,
					"data":      nil,
					"error":     "",
				}
				itemName := columns[1]
				quantityRaw := columns[2]
				if quantityRaw == "" {
					logger.Debugf("quantity is required")
					status["error"] = "quantity is required"
					productStatusesByInvoiceNo[invoiceNo] = append(productStatusesByInvoiceNo[invoiceNo], status)
					continue
				}
				quantity, err := strconv.Atoi(quantityRaw)
				if err != nil {
					logger.Debugf("invalid quantity: %s", quantityRaw)
					status["error"] = fmt.Sprintf("invalid quantity: %s", quantityRaw)
					productStatusesByInvoiceNo[invoiceNo] = append(productStatusesByInvoiceNo[invoiceNo], status)
					continue
				}
				totalCOGSRaw := columns[3]
				if totalCOGSRaw == "" {
					logger.Debugf("total cogs is required")
					status["error"] = "total cogs is required"
					productStatusesByInvoiceNo[invoiceNo] = append(productStatusesByInvoiceNo[invoiceNo], status)
					continue
				}
				totalCOGS, err := strconv.ParseInt(totalCOGSRaw, 10, 64)
				if err != nil {
					logger.Debugf("invalid total cogs: %s", totalCOGSRaw)
					status["error"] = fmt.Sprintf("invalid total cogs: %s", totalCOGSRaw)
					productStatusesByInvoiceNo[invoiceNo] = append(productStatusesByInvoiceNo[invoiceNo], status)
					continue
				}
				totalPriceRaw := columns[4]
				if totalPriceRaw == "" {
					logger.Debugf("total price is required")
					status["error"] = "total price is required"
					productStatusesByInvoiceNo[invoiceNo] = append(productStatusesByInvoiceNo[invoiceNo], status)
					continue
				}
				totalPrice, err := strconv.ParseInt(totalPriceRaw, 10, 64)
				if err != nil {
					logger.Debugf("invalid total price: %s", totalPriceRaw)
					status["error"] = fmt.Sprintf("invalid total price: %s", totalPriceRaw)
					productStatusesByInvoiceNo[invoiceNo] = append(productStatusesByInvoiceNo[invoiceNo], status)
					continue
				}

				// new product
				newProduct := &model.Product{
					UUID:                 uuid.New().String(),
					ItemName:             itemName,
					Quantity:             quantity,
					TotalCostOfGoodsSold: totalCOGS,
					TotalPriceSold:       totalPrice,
					// fill up InvoiceID, InvoiceUUID, InvoiceNo later on
				}

				// add to statuses even if there is an error for next invoice sheet proccess
				status["data"] = newProduct
				productStatusesByInvoiceNo[invoiceNo] = append(productStatusesByInvoiceNo[invoiceNo], status)
			}
		}
	}

	// search for invoice sheet
	statuses := make(map[string]map[string]interface{})
	/*
		{
			[invoice_no]: {
				"invoice_no":       [invoice_no],
				"row_index":        [row_index],
				"error":            "",
				"product_statuses": []
			}
		}
	*/
	for _, sheetName := range f.GetSheetList() {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return err
		}
		if strings.Contains(strings.ToLower(sheetName), "invoice") {
			for i, columns := range rows {
				if i == 0 { // header
					// ensure columns total
					if len(columns) != 6 {
						return &error_utils.CustomErr{
							HttpCode: 400,
							Message:  "invalid xlsx file",
							Detail:   "columns total for invoice sheet must be 6",
							Data:     nil,
						}
					}

					// skip header
					continue
				}
				logger.Debugf("invoice columns: %v", columns)

				invoiceNo := columns[0]
				invoiceStatus := map[string]interface{}{
					"invoice_no":       invoiceNo,
					"row_index":        i,
					"error":            "",
					"product_statuses": make([]map[string]interface{}, 0),
				}

				// extract fields
				dateRaw := columns[1]
				logger.Debugf("raw date: %s", dateRaw)
				if dateRaw == "" {
					logger.Debugf("date required")
					invoiceStatus["error"] = fmt.Sprintf("date required")
					statuses[invoiceNo] = invoiceStatus
					continue
				}
				date, err := time.Parse("02-01-2006", dateRaw)
				layoutAttempts := []string{"02-01-2006", "02/01/2006", "01-02-06"}
				for i, layoutAttempt := range layoutAttempts {
					date, err = time.Parse(layoutAttempt, dateRaw)
					if err == nil {
						break
					}
					if i == len(layoutAttempts)-1 {
						logger.Debugf("error parsing date: %v", err)
						invoiceStatus["error"] = fmt.Sprintf("invalid date, date format must be one of these: %v", layoutAttempts)
						statuses[invoiceNo] = invoiceStatus
						continue
					}
				}
				customer := columns[2]
				salesPerson := columns[3]
				paymentTypeRaw := columns[4]
				if paymentTypeRaw == "" {
					logger.Debugf("payment type required")
					invoiceStatus["error"] = fmt.Sprintf("payment_type required")
					statuses[invoiceNo] = invoiceStatus
					continue
				}
				paymentType := enum.InvoicePaymentType(paymentTypeRaw)
				ok := paymentType.IsValid()
				if !ok {
					logger.Debugf("error parsing payment_type: %v", err)
					invoiceStatus["error"] = fmt.Sprintf("invalid payment_type")
					statuses[invoiceNo] = invoiceStatus
					continue
				}
				notesRaw := columns[5]
				var notes *string
				if notesRaw == "" {
					notes = nil
				} else {
					notes = &notesRaw
				}

				// new invoice obj
				newInvoice := &model.Invoice{
					UUID:            uuid.New().String(),
					InvoiceNo:       invoiceNo,
					Date:            date,
					CustomerName:    customer,
					SalesPersonName: salesPerson,
					PaymentType:     paymentType,
					Notes:           notes,
				}

				// validate
				err = newInvoice.Validate()
				if err != nil {
					logger.Debugf("error validating invoice: %v", err)
					invoiceStatus["error"] = err.Error()
					statuses[invoiceNo] = invoiceStatus
					continue
				}

				// prepare products for newInvoice
				productStatuses := productStatusesByInvoiceNo[invoiceNo]
				erroredProductStatuses := make([]map[string]interface{}, 0)
				invoiceProducts := make([]model.Product, 0)
				if len(productStatuses) != 0 {
					for _, status := range productStatuses {
						productData, ok := status["data"].(*model.Product)
						if ok {
							// fill up InvoiceID and InvoiceUUID
							productData.InvoiceID = newInvoice.ID
							productData.InvoiceUUID = newInvoice.UUID
							productData.InvoiceNo = newInvoice.InvoiceNo
							invoiceProducts = append(invoiceProducts, *productData)
						}

						if status["error"] != "" { // only append if there is an error in product
							erroredProductStatuses = append(erroredProductStatuses, status)
						}
					}
				}

				// set newInvoice.Products
				newInvoice.Products = invoiceProducts

				// mark invoice as errored if related errors are found
				if len(erroredProductStatuses) > 0 {
					logger.Debugf("error on products: %v", erroredProductStatuses)
					invoiceStatus["product_statuses"] = erroredProductStatuses
					invoiceStatus["error"] = fmt.Sprintf("error on products")
					statuses[invoiceNo] = invoiceStatus
					continue
				}

				// check if invoiceNo is already exists
				existingInvoice, err := u.invoiceRepo.GetByInvoiceNo(newInvoice.InvoiceNo)
				if existingInvoice != nil {
					logger.Debugf("invoice %s already exists", newInvoice.InvoiceNo)
					invoiceStatus["error"] = fmt.Sprintf("invoice %s already exists", newInvoice.InvoiceNo)
					statuses[invoiceNo] = invoiceStatus
					continue
				}

				// insert
				err = u.invoiceRepo.Create(newInvoice)
				if err != nil {
					invoiceStatus["error"] = fmt.Sprintf("error when inserting invoice: %v", err.Error())
					statuses[invoiceNo] = invoiceStatus
					continue
				}
			}
		}
	}

	if len(statuses) > 0 {
		statusResp := make([]map[string]interface{}, 0)
		for _, status := range statuses {
			statusResp = append(statusResp, status)
		}
		return &error_utils.CustomErr{
			HttpCode: 207,
			Message:  "partial success, some rows failed to proccess",
			Detail:   "some rows failed to proccess",
			Data:     statusResp,
		}
	}

	return nil
}
