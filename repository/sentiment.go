package repository

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

type SentimentRepository struct {
	redisClient *redis.Client
}

var ctx = context.Background()
var key = "sentiment"

func NewSentimentRepository(redisClient *redis.Client) *SentimentRepository {
	return &SentimentRepository{redisClient: redisClient}
}

func (r *SentimentRepository) GetSentiment(sentiment string) float64 {
	c, _ := r.redisClient.Get(ctx, fmt.Sprintf("%s:%s", key, sentiment)).Float64()
	return c
}

func (r *SentimentRepository) IncrementSentiment(sentiment string, value float64) error {
	_, err := r.redisClient.IncrByFloat(ctx, fmt.Sprintf("%s:%s", key, sentiment), value).Result()
	return err
}

func (r *SentimentRepository) SetSentiment(sentiment string, value float64) error {
	_, err := r.redisClient.Set(ctx, fmt.Sprintf("%s:%s", key, sentiment), value, 0).Result()
	return err
}
