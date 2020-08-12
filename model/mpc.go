package model

type MpcTask struct {
	Nonce   string `json:"nonce"`
	Owner   string `json:"owner"`
	Data    string `json:"data"`
	Sponsor string `json:"sponsor"`
}
