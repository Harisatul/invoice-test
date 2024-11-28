package model

import "time"

type PaymentStatus string

const (
	PaymentStatusCASH   PaymentStatus = "CASH"
	PaymentStatusCREDIT PaymentStatus = "CREDIT"
)

type Invoice struct {
	InvoiceNumber string
	Date          time.Time
	CustomerName  string
	Salesperson   string
	Notes         string
	PaymentType   PaymentStatus
	Products      []Product
}

func Validate() {

}
