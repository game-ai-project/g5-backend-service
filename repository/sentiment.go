package repository

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/staciiaz/g5-backend-service/model"
)

type SentimentRepository struct {
	redisClient *redis.Client
}

var ctx = context.Background()

func NewSentimentRepository(redisClient *redis.Client) *SentimentRepository {
	return &SentimentRepository{redisClient: redisClient}
}

func (r *SentimentRepository) SentimentRatio() *model.SentimentRatio {
	positive := r.GetSentiment("positive")
	neutral := r.GetSentiment("neutral")
	negative := r.GetSentiment("negative")

	sum := positive + neutral + negative
	if sum == 0 {
		sum = 1
	}

	return &model.SentimentRatio{
		Positive: positive / sum,
		Neutral:  neutral / sum,
		Negative: negative / sum,
	}
}

func (r *SentimentRepository) GetSentiment(sentiment string) float64 {
	c, _ := r.redisClient.Get(ctx, sentiment).Float64()
	return c
}

func (r *SentimentRepository) IncrementSentiment(sentiment string, value float64) error {
	_, err := r.redisClient.IncrByFloat(ctx, sentiment, value).Result()
	return err
}

func (r *SentimentRepository) ClearSentiment() error {
	var err error

	_, err = r.redisClient.Set(ctx, "positive", 0, 0).Result()
	if err != nil {
		return err
	}

	_, err = r.redisClient.Set(ctx, "neutral", 0, 0).Result()
	if err != nil {
		return err
	}

	_, err = r.redisClient.Set(ctx, "negative", 0, 0).Result()
	if err != nil {
		return err
	}

	return nil
}
