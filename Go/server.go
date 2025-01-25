package main

import (
	"backend/domain/dto"
	"backend/handler"
	"backend/utils/http_response"

	_ "backend/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupServer(ginEngine *gin.Engine, deps CommonDeps) {
	responseWriter := http_response.NewHttpResponseWriter()
	_ = responseWriter

	// handlers
	invoiceHandler := handler.NewInvoiceHandler(responseWriter, deps.invoiceUcase)

	// register routes
	router := ginEngine
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, dto.BaseJSONResp{
			Code:    200,
			Message: "pong",
		})
	})

	// invoice
	router.POST("/invoices", invoiceHandler.CreateInvoice)
	router.PUT("/invoices/:invoice_uuid", invoiceHandler.UpdateInvoice)
	router.DELETE("/invoices/:invoice_uuid", invoiceHandler.DeleteInvoice)
	router.DELETE("/invoices/no/:invoice_no", invoiceHandler.DeleteInvoiceByInvoiceNo)
	router.GET("/invoices/:invoice_uuid", invoiceHandler.GetInvoiceDetail)
	router.GET("/invoices", invoiceHandler.GetInvoiceList)

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/swagger/index.html")
	})
}
