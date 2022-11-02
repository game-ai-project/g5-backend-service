package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/staciiaz/g5-backend-service/service"
)

type RouterController struct {
	sentimentService *service.SentimentService
}

func NewRouterController(sentimentService *service.SentimentService) *RouterController {
	return &RouterController{sentimentService: sentimentService}
}

func (c *RouterController) Setup(router *gin.Engine) {

	router.GET("/poll", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, c.sentimentService.Poll())
	})

	router.GET("/refresh", func(ctx *gin.Context) {
		err := c.sentimentService.Refresh()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
}
