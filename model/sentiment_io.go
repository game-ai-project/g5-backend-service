package model

type SentimentIO struct {
	Documents []Document `json:"documents"`
}

type Document struct {
	Language         string           `json:"language"`
	ID               string           `json:"id"`
	Text             string           `json:"text"`
	Sentiment        string           `json:"sentiment"`
	ConfidenceScores ConfidenceScores `json:"confidenceScores"`
	Sentences        []Sentence       `json:"sentences"`
}

type ConfidenceScores struct {
	Positive float64 `json:"positive"`
	Neutral  float64 `json:"neutral"`
	Negative float64 `json:"negative"`
}

type Sentence struct {
	Sentiment        string           `json:"sentiment"`
	ConfidenceScores ConfidenceScores `json:"confidenceScores"`
	Offset           int              `json:"offset"`
	Length           int              `json:"length"`
	Text             string           `json:"text"`
}

func (cs *ConfidenceScores) GetMax() float64 {
	max := cs.Positive
	if cs.Negative > max {
		max = cs.Negative
	}
	if cs.Neutral > max {
		max = cs.Neutral
	}
	return max
}
