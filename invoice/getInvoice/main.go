package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/projects/my_resturant_serverless/database"
	"github.com/projects/my_resturant_serverless/itemsByOrder"
	"github.com/projects/my_resturant_serverless/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func Handler(ctx context.Context, req Request)(*Response, error){
invoiceId := req.PathParameters["invoiceId"]
var invoice models.Invoice
var invoiceView models.InvoiceViewFormat
err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
if err != nil {
	return &Response{
		StatusCode: http.StatusNotFound,
		Body: fmt.Sprintf("Invoice not found"),
	},nil
}
allOrderItems, err := ItemsByOrder.ItemsByOrder(invoice.Order_id)
invoiceView.Order_id = invoice.Order_id
invoiceView.Payment_due_date = invoice.Payment_due_date
invoice.Payment_method = nil
if invoice.Payment_method != nil {
	invoiceView.Payment_method = *invoice.Payment_method
}

invoiceView.Invoice_id = invoice.Invoice_id

invoiceView.Payment_status = *&invoice.Payment_status
invoiceView.Payment_due = allOrderItems[0]["payment_due"]
invoiceView.Table_number = allOrderItems[0]["table_number"]
invoiceView.Order_details = allOrderItems[0]["order_items"]

jsonBytes, err := json.Marshal(invoiceView)
return &Response{
	StatusCode: http.StatusOK,
	Body: string(jsonBytes),
},nil
}

func main(){
	lambda.Start(Handler)
}