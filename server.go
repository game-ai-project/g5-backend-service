package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	socketio "github.com/googollee/go-socket.io"
	"github.com/staciiaz/g5-backend-service/controller"
	"github.com/staciiaz/g5-backend-service/repository"
	"github.com/staciiaz/g5-backend-service/service"
	"go.uber.org/dig"
)

var ctx = context.Background()

func ProvideRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASS"),
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func ProvideSentimentService(repo *repository.SentimentRepository) *service.SentimentService {
	return service.NewSentimentService(os.Getenv("SENTIMENT_SERVICE"), repo)
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
		if err := container.Provide(c); err != nil {
			panic(err)
		}
	}

	server := socketio.NewServer(nil)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
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
