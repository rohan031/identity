package services

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"string,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
