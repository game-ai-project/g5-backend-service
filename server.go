package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	socketio "github.com/googollee/go-socket.io"
	"github.com/staciiaz/g5-backend-service/controller"
	"github.com/staciiaz/g5-backend-service/repository"
	"github.com/staciiaz/g5-backend-service/service"
	"go.uber.org/dig"
)

func ProvideRedisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis-18767.c302.asia-northeast1-1.gce.cloud.redislabs.com:18767",
		Username: "default",
		Password: "uyyKIkPB6o9QR1J81g1Xruy2vFKIlKWM",
	})
	return redisClient
}

func ProvideSentimentService(repo *repository.SentimentRepository) *service.SentimentService {
	return service.NewSentimentService("http://localhost:5000", repo)
}

func main() {
	container := dig.New()

	constructors := []interface{}{
		ProvideRedisClient,
		repository.NewSentimentRepository,
		ProvideSentimentService,
		controller.NewSentimentController,
	}
	for _, c := range constructors {
		container.Provide(c)
	}

	server := socketio.NewServer(nil)
	router := gin.Default()
	router.Use(gin.Logger())

	container.Invoke(func(service *service.SentimentService) {
		service.Refresh()
	})

	container.Invoke(func(controller *controller.SentimentController) {
		controller.Provide(server, router)
		controller.Setup()
	})

	go server.Serve()
	defer server.Close()

	log.Fatal(router.Run(":8000"))
}
