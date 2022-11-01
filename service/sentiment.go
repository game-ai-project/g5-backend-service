package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/staciiaz/g5-backend-service/model"
)

type SentimentService struct {
	apiEndpoint string
}

func NewSentimentService(apiEndpoint string) *SentimentService {
	return &SentimentService{apiEndpoint: apiEndpoint}
}

func (s *SentimentService) GetSentiment(input *model.SentimentIO) (*model.SentimentIO, error) {
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
