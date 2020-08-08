package model

type Article struct {
	Id    string
	Score float32
}

type QueryIDsRes struct {
	A []Article
}

type QueryLog struct {
	Owner     string `json:"owner"`
	Payload   string `json:"payload"`
	TimeStamp string `json:"timestamp"`
}
