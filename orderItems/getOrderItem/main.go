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
	"github.com/projects/my_resturant_serverless/models"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

var orderItemCollection  *mongo.Collection = database.OpenCollection(database.Client, "OrderItem")

func Handler( ctx context.Context, req Request)(*Response, error){
	var orderItem models.OrderItems

	orderItemId := req.PathParameters["orderItemId"]

	err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)

	if err != nil{
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Body: fmt.Sprintf("Error ffetching orderItems"),
		},nil
	}

	jsonBytes, err := json.Marshal(orderItem)

	if err != nil{
		log.Fatal(err)
	}

	return &Response{
		StatusCode: http.StatusOK,
		Body: string(jsonBytes),
	},nil
}

func main(){
	lambda.Start(Handler)
}