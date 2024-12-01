package model

type InvoiceAggregateResponse struct {
	TotalProfit            int64 `json:"total_profit"`
	TotalOfCashTransaction int64 `json:"total_of_cash_transaction"`
}

type GetInvoiceResponse struct {
	Invoice                  []Invoice
	InvoiceAggregateResponse InvoiceAggregateResponse
}
