package model

type InvoiceAggregateResponse struct {
	TotalProfit            int `json:"total_profit"`
	TotalOfCashTransaction int `json:"total_of_cash_transaction"`
}

type GetInvoiceResponse struct {
	Invoice                  []Invoice
	InvoiceAggregateResponse InvoiceAggregateResponse
}
