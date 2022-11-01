package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	socketio "github.com/googollee/go-socket.io"
	"github.com/staciiaz/g5-backend-service/model"
	"github.com/staciiaz/g5-backend-service/repository"
	"github.com/staciiaz/g5-backend-service/service"
	"go.uber.org/dig"
)

func main() {
	c := dig.New()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis-18767.c302.asia-northeast1-1.gce.cloud.redislabs.com:18767",
		Username: "default",
		Password: "uyyKIkPB6o9QR1J81g1Xruy2vFKIlKWM",
	})

	server := socketio.NewServer(nil)
	sentimentRepository := repository.NewSentimentRepository(redisClient)
	sentimentService := service.NewSentimentService("http://104.215.1.237:5000")

	c.Provide(redisClient)
	c.Provide(server)
	c.Provide(sentimentRepository)
	c.Provide(sentimentService)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		return nil
	})

	server.OnEvent("/", "sentiment", func(s socketio.Conn, msg string) {
		result, _ := sentimentService.GetSentiment(&model.SentimentIO{
			Documents: []model.Document{
				{
					Language: "en",
					ID:       "1",
					Text:     msg,
				},
			},
		})
		for _, document := range result.Documents {
			sentiment := document.Sentiment
			sentimentRepository.IncrementSentiment(sentiment, 1.0)
			s.Emit("sentiment", sentiment)
		}
	})

	server.OnError("/", func(s socketio.Conn, e error) {

	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {

	})

	go server.Serve()
	defer server.Close()

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	router.StaticFS("/public", http.Dir("../asset"))

	router.GET("/poll", func(c *gin.Context) {
		ratio := sentimentRepository.SentimentRatio()
		sentimentRepository.ClearSentiment()
		c.JSON(http.StatusOK, ratio)
	})

	log.Fatal(router.Run(":8000"))
}
