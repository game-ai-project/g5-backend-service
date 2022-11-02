package model

type MessagePayload struct {
	Message string `json:"message"`
}

type SentimentPayload struct {
	Result     string  `json:"result"`
	Confidence float64 `json:"confidence"`
}
