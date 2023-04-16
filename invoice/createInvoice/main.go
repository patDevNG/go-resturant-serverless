package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"github.com/projects/my_resturant_serverless/database"
	"github.com/projects/my_resturant_serverless/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

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

var validate = validator.New()
var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")
var orderCollection *mongo.Collection = database.OpenCollection(database.Client,"order")
func Handler(ctx context.Context, req Request)(*Response, error){
var invoice models.Invoice
var order models.Order

if err := json.Unmarshal([]byte(req.Body), &invoice); err != nil {
	return &Response{
		StatusCode: http.StatusBadRequest,
		Body: fmt.Sprintf("An input error has occurred"),
	}, nil
}
err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)
if err != nil {
	return &Response{
		StatusCode: http.StatusNotFound,
		Body: fmt.Sprintf("Order not found"),
	},nil
}
status := "PAID"
invoice.Payment_status = &status
invoice.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().AddDate(0,0,1).Format(time.RFC3339))
invoice.Created_at, _ = time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
invoice.ID = primitive.NewObjectID()
invoice.Invoice_id = invoice.ID.Hex()

validationErr := validate.Struct(invoice)
if validationErr != nil {
	return &Response{
		StatusCode: http.StatusBadRequest,
		Body: fmt.Sprintf("Validfation error occured"),
	},nil
}
res, insertErr := invoiceCollection.InsertOne(ctx, invoice);
if insertErr !=nil{
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Body: fmt.Sprintf("Error inserting in database"),
	},nil
}
jsonBytes, err := json.Marshal(res)
if err != nil{
	log.Fatal(err)
}
return &Response{
	StatusCode: http.StatusOK,
	Body: string(jsonBytes),
},nil
}

func main (){
	lambda.Start(Handler)
}