package main

import (
	"backend/domain/dto"
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
	// register routes
	router := ginEngine
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, dto.BaseJSONResp{
			Code:    200,
			Message: "pong",
		})
	})

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/swagger/index.html")
	})
}
