package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/projects/my_resturant_serverless/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func Handler(ctx context.Context, _ Request)(*Response, error){
	res, err := invoiceCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Body: fmt.Sprintf("Error fetching invoices"),
		},nil
	}

	var allInvoices []bson.M
	if err = res.All(ctx, &allInvoices); err != nil {
		log.Fatal(err)
	}
	jsonBytes, err := json.Marshal(allInvoices)
	if err != nil{
		log.Fatal(err)
	}
	return &Response{
		StatusCode: http.StatusOK,
		Body:string(jsonBytes) ,
	},nil
}

func main(){
	lambda.Start(Handler)
}