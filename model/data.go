package model

type QueryLog struct {
	Owner     string `json:"owner"`
	Payload   string `json:"payload"`
	TimeStamp string `json:"timestamp"`
}
