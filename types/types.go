package types

type Response struct {
	Status  int    `json:"status`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type ID struct {
	ID string `json:"id"`
}
