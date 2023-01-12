package controller

import (
	"log"
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
		s.SetContext("")
		return nil
	})

	c.server.OnEvent("/", "message", func(s socketio.Conn, data *model.MessagePayload) {
		result, err := c.sentimentService.DoAnalytic(&model.SentimentIO{
			Documents: []model.Document{
				{
					Language: "en",
					ID:       "1",
					Text:     data.Message,
				},
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		s.Emit("result", result)
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
		ctx.JSON(http.StatusOK, data)
	})

	c.router.GET("/socket.io/*any", gin.WrapH(c.server))
	c.router.POST("/socket.io/*any", gin.WrapH(c.server))
	c.router.StaticFS("/public", http.Dir("../asset"))
}

func (c *SentimentController) Setup() {
	c.setupSocketIO()
	c.setupRouter()
}
