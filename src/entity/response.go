package entity

// Response is the structure for Response
type Response struct {
	Status string      `json:"status"`
	Info   interface{} `json:"info"`
}
