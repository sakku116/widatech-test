package handler

import (
	"backend/domain/dto"
	ucase "backend/usecase"
	"backend/utils/http_response"

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
// @Router /invoices/{invoice_uuid} [put]
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
// @Router /invoices/{invoice_no} [delete]
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
