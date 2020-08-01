package model

type QueryIDsRes struct {
	A []string
}

type QueryLog struct {
	Owner     string `json:"owner"`
	Payload   string `json:"payload"`
	TimeStamp string `json:"timestamp"`
}
