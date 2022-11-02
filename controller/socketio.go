package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/staciiaz/g5-backend-service/model"
	"github.com/staciiaz/g5-backend-service/service"
)

type SocketIOController struct {
	sentimentService *service.SentimentService
}

func NewSocketIOController(sentimentService *service.SentimentService) *SocketIOController {
	return &SocketIOController{sentimentService: sentimentService}
}

func (c *SocketIOController) Setup(router *gin.Engine, server *socketio.Server) {
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		return nil
	})

	server.OnEvent("/", "message", func(s socketio.Conn, data *model.MessagePayload) {
		result, _ := c.sentimentService.DoSingleAnalytic(data.Message)
		s.Emit("message", result)
	})

	server.OnError("/", func(s socketio.Conn, e error) {

	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {

	})

	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	router.StaticFS("/public", http.Dir("../asset"))
}
