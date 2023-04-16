package models

import "time"

type InvoiceViewFormat struct {
	Invoice_id           string
	Payment_method       string
	Order_id             string
	Payment_status       *string
	Payment_due          interface{}
	Table_number         interface{}
	Payment_due_date     time.Time
	Order_details        interface{}
}