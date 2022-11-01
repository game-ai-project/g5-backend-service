package model

type Sentence struct {
	Sentiment        string           `json:"sentiment"`
	ConfidenceScores ConfidenceScores `json:"confidenceScores"`
	Offset           int              `json:"offset"`
	Length           int              `json:"length"`
	Text             string           `json:"text"`
}
