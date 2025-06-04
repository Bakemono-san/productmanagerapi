package types

type Response struct {
	Status  int    `json:"status`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type ID struct {
	ID string `json:"id"`
}

type ProductSale struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type SaleRequest struct {
	Products []ProductSale `json:"products"`
}
