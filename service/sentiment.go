package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/staciiaz/g5-backend-service/model"
	"github.com/staciiaz/g5-backend-service/repository"
)

type SentimentService struct {
	apiEndpoint string
	repo        *repository.SentimentRepository
}

func NewSentimentService(apiEndpoint string, repo *repository.SentimentRepository) *SentimentService {
	return &SentimentService{apiEndpoint: apiEndpoint, repo: repo}
}

func (s *SentimentService) Poll() *model.RatioForPacMan {
	positive := s.repo.GetSentiment("positive")
	neutral := s.repo.GetSentiment("neutral")
	negative := s.repo.GetSentiment("negative")

	positive += neutral / 2
	negative += neutral / 2
	sum := positive + negative

	return &model.RatioForPacMan{
		Cheer: positive / sum,
		Jeer:  negative / sum,
	}
}

func (s *SentimentService) Refresh() error {
	var err error

	err = s.repo.SetSentiment("positive", 1)
	if err != nil {
		return err
	}

	err = s.repo.SetSentiment("neutral", 0)
	if err != nil {
		return err
	}

	err = s.repo.SetSentiment("negative", 1)
	if err != nil {
		return err
	}

	return nil
}

func (s *SentimentService) getSentiment(input *model.SentimentIO) (*model.SentimentIO, error) {
	apiPath := "/text/analytics/v3.0/sentiment"
	url := fmt.Sprintf("%s%s", s.apiEndpoint, apiPath)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result model.SentimentIO
	body, _ = io.ReadAll(res.Body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *SentimentService) DoAnalytic(input *model.SentimentIO) (*model.SentimentIO, error) {
	result, err := s.getSentiment(input)
	if err != nil {
		return nil, err
	}

	for _, document := range result.Documents {
		sentiment := document.Sentiment
		s.repo.IncrementSentiment(sentiment, 1.0)
	}

	return result, nil
}
