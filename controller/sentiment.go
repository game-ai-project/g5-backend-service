package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/staciiaz/g5-backend-service/model"
	"github.com/staciiaz/g5-backend-service/service"
)

type SentimentController struct {
	sentimentService *service.SentimentService

	server *socketio.Server
	router *gin.Engine
}

func NewSentimentController(sentimentService *service.SentimentService) *SentimentController {
	return &SentimentController{sentimentService: sentimentService}
}

func (c *SentimentController) Provide(server *socketio.Server, router *gin.Engine) {
	c.server = server
	c.router = router
}

func (c *SentimentController) setupSocketIO() {
	c.server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext(context.Background())
		s.RemoteHeader().Add("Access-Control-Allow-Origin", "*")
		s.RemoteHeader().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		s.RemoteHeader().Add("Access-Control-Allow-Headers", "Content-Type")
		return nil
	})

	c.server.OnError("/", func(s socketio.Conn, e error) {

	})

	c.server.OnDisconnect("/", func(s socketio.Conn, reason string) {

	})
}

func (c *SentimentController) setupRouter() {
	c.router.GET("/poll", func(ctx *gin.Context) {
		query := ctx.Request.URL.Query()
		data := c.sentimentService.Poll()
		if query.Get("silent") != "true" {
			c.server.BroadcastToNamespace("/", "message", data)
			c.sentimentService.Refresh()
		}

		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.JSON(http.StatusOK, data)
	})

	c.router.POST("/sentiment", func(ctx *gin.Context) {
		var payload model.MessagePayload
		if err := ctx.BindJSON(&payload); err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		result, err := c.sentimentService.DoAnalytic(&model.SentimentIO{
			Documents: []model.Document{
				{
					Language: "en",
					ID:       "1",
					Text:     payload.Message,
				},
			},
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.JSON(http.StatusOK, result)
	})

	c.router.StaticFS("/public", http.Dir("../asset"))
}

func (c *SentimentController) Setup() {
	c.setupSocketIO()
	c.setupRouter()
}
