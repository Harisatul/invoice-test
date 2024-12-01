package model

type CreateInvoiceRequest struct {
	CustomerName    string `json:"customer_name"`
	SalesPersonName string `json:"sales_person_name"`
	PaymentType     string `json:"payment_type"`
	Notes           string `json:"notes"`
	Product         []struct {
		Name                string `json:"item_name"`
		Quantity            int    `json:"quantity"`
		TotalCostOfGoodSold int64  `json:"total_cogs"`
		TotalPriceSold      int64  `json:"total_price_sold"`
	} `json:"product"`
}

type UpdateInvoiceRequest struct {
	CustomerName    string `json:"customer_name"`
	SalesPersonName string `json:"sales_person_name"`
	PaymentType     string `json:"payment_type"`
	Notes           string `json:"notes"`
	Product         []struct {
		Name                string `json:"item_name"`
		Quantity            int    `json:"quantity"`
		TotalCostOfGoodSold int64  `json:"total_cogs"`
		TotalPriceSold      int64  `json:"total_price_sold"`
	} `json:"product"`
}
