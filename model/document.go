package model

type Document struct {
	Language         string           `json:"language"`
	ID               string           `json:"id"`
	Text             string           `json:"text"`
	Sentiment        string           `json:"sentiment"`
	ConfidenceScores ConfidenceScores `json:"confidenceScores"`
	Sentences        []Sentence       `json:"sentences"`
}
