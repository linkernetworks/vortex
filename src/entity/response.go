package entity

type Response struct {
	Status string      `json:"status"`
	Info   interface{} `json:"info"`
}
