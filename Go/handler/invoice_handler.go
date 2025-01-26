package handler

import (
	"backend/domain/dto"
	ucase "backend/usecase"
	"backend/utils/http_response"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
	respWriter   http_response.IHttpResponseWriter
	invoiceUcase ucase.IInvoiceUcase
}

type IInvoiceHandler interface {
	CreateInvoice(c *gin.Context)
	UpdateInvoice(c *gin.Context)
	DeleteInvoice(c *gin.Context)
	DeleteInvoiceByInvoiceNo(c *gin.Context)
	GetInvoiceDetail(c *gin.Context)
	GetInvoiceList(c *gin.Context)
	ImportFromXlsx(c *gin.Context)
}

func NewInvoiceHandler(
	respWriter http_response.IHttpResponseWriter,
	invoiceUcase ucase.IInvoiceUcase,
) IInvoiceHandler {
	return &InvoiceHandler{
		respWriter:   respWriter,
		invoiceUcase: invoiceUcase,
	}
}

// @Summary create new invoice
// @Tags invoice
// @Accept json
// @Produce json
// @Param payload body dto.CreateInvoiceReq true "create invoice request"
// @Success 200 {object} dto.BaseJSONResp{data=dto.CreateInvoiceResp}
// @Router /invoices [post]
func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	var payload dto.CreateInvoiceReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.respWriter.HTTPJson(
			c, 400, "invalid request", err.Error(), nil,
		)
		return
	}

	resp, err := h.invoiceUcase.CreateInvoice(payload)
	if err != nil {
		h.respWriter.HTTPCustomErr(c, err)
		return
	}

	h.respWriter.HTTPJson(
		c, 200, "success", "", resp,
	)
}

// @Summary update invoice
// @Tags invoice
// @Accept json
// @Produce json
// @Param invoice_uuid path string true "invoice uuid"
// @Param payload body dto.UpdateInvoiceReq true "update invoice request"
// @Success 200 {object} dto.BaseJSONResp{data=dto.UpdateInvoiceRespData}
// @Router /invoices/{invoice_uuid} [patch]
func (h *InvoiceHandler) UpdateInvoice(c *gin.Context) {
	invoiceUUID := c.Param("invoice_uuid")
	var payload dto.UpdateInvoiceReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.respWriter.HTTPJson(
			c, 400, "invalid request", err.Error(), nil,
		)
		return
	}

	resp, err := h.invoiceUcase.UpdateInvoice(invoiceUUID, payload)
	if err != nil {
		h.respWriter.HTTPCustomErr(c, err)
		return
	}

	h.respWriter.HTTPJson(
		c, 200, "success", "", resp,
	)
}

// @Summary delete invoice
// @Tags invoice
// @Accept json
// @Produce json
// @Param invoice_uuid path string true "invoice uuid"
// @Success 200 {object} dto.BaseJSONResp
// @Router /invoices/{invoice_uuid} [delete]
func (h *InvoiceHandler) DeleteInvoice(c *gin.Context) {
	invoiceUUID := c.Param("invoice_uuid")

	err := h.invoiceUcase.DeleteInvoice(invoiceUUID)
	if err != nil {
		h.respWriter.HTTPCustomErr(c, err)
		return
	}

	h.respWriter.HTTPJson(
		c, 200, "success", "", nil,
	)
}

// @Summary delete invoice by invoice no
// @Tags invoice
// @Accept json
// @Produce json
// @Param invoice_no path string true "invoice no"
// @Success 200 {object} dto.BaseJSONResp
// @Router /invoices/no/{invoice_no} [delete]
func (h *InvoiceHandler) DeleteInvoiceByInvoiceNo(c *gin.Context) {
	invoiceNo := c.Param("invoice_no")

	err := h.invoiceUcase.DeleteByInvoiceNo(invoiceNo)
	if err != nil {
		h.respWriter.HTTPCustomErr(c, err)
		return
	}

	h.respWriter.HTTPJson(
		c, 200, "success", "", nil,
	)
}

// @Summary get invoice detail
// @Tags invoice
// @Accept json
// @Produce json
// @Param invoice_uuid path string true "invoice uuid"
// @Success 200 {object} dto.BaseJSONResp{data=dto.GetInvoiceDetailRespData}
// @Router /invoices/{invoice_uuid} [get]
func (h *InvoiceHandler) GetInvoiceDetail(c *gin.Context) {
	invoiceUUID := c.Param("invoice_uuid")

	resp, err := h.invoiceUcase.GetInvoiceDetail(invoiceUUID)
	if err != nil {
		h.respWriter.HTTPCustomErr(c, err)
		return
	}

	h.respWriter.HTTPJson(
		c, 200, "success", "", resp,
	)
}

// @Summary get invoice list
// @Description set date_from and date_to to get profit_total and cash_transaction_total.\nif date_from and date_to is set, pagination will be ignored for profit calculation.
// @Tags invoice
// @Accept json
// @Produce json
// @Param payload query dto.GetInvoiceListReq true "get invoice list request"
// @Success 200 {object} dto.BaseJSONResp{data=dto.GetInvoiceListRespData}
// @Router /invoices [get]
func (h *InvoiceHandler) GetInvoiceList(c *gin.Context) {
	var payload dto.GetInvoiceListReq
	if err := c.ShouldBindQuery(&payload); err != nil {
		h.respWriter.HTTPJson(
			c, 400, "invalid request", err.Error(), nil,
		)
		return
	}

	resp, err := h.invoiceUcase.GetInvoiceList(payload)
	if err != nil {
		h.respWriter.HTTPCustomErr(c, err)
		return
	}

	h.respWriter.HTTPJson(
		c, 200, "success", "", resp,
	)
}

// @Summary import invoice from xlsx
// @Tags invoice
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "xlsx file"
// @Success 200 {object} dto.BaseJSONResp
// @Router /invoices/import [post]
func (h *InvoiceHandler) ImportFromXlsx(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		h.respWriter.HTTPJson(
			c, 400, "invalid request", err.Error(), nil,
		)
		return
	}

	// validate file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".xlsx" {
		h.respWriter.HTTPJson(
			c, 400, "invalid request", "invalid file extension, must be xlsx", nil,
		)
		return
	}

	// open file
	uploadedFile, err := file.Open()
	if err != nil {
		h.respWriter.HTTPJson(
			c, 500, "internal server error", err.Error(), nil,
		)
		return
	}
	defer uploadedFile.Close()

	// save to tmp folder
	tempFilePath := "./tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		h.respWriter.HTTPJson(
			c, 500, "internal server error", err.Error(), nil,
		)
		return
	}

	err = h.invoiceUcase.ImportXlsx(tempFilePath)
	if err != nil {
		h.respWriter.HTTPCustomErr(c, err)
		return
	}

	// clean up tmp file
	if err := os.Remove(tempFilePath); err != nil {
		h.respWriter.HTTPJson(
			c, 500, "internal server error", err.Error(), nil,
		)
		return
	}

	h.respWriter.HTTPJson(
		c, 200, "success", "", nil,
	)
}
