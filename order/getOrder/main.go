package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/projects/my_resturant_serverless/database"
	"github.com/projects/my_resturant_serverless/models"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)
var orderCollection *mongo.Collection = database.OpenCollection(database.Client,"order")
type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest
func Handler(ctx context.Context, req Request)(*Response, error){
orderId := req.PathParameters["orderId"]
var order models.Order
err:= orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)

if err != nil{
	return &Response{
		StatusCode: http.StatusNotFound,
		Body: fmt.Sprintf("Requested order is not found"),
	},nil
}
jsonBytes, err := json.Marshal(order)
if err != nil {
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Body: fmt.Sprintf("An error occured"),
	},nil
}
return &Response{
	StatusCode: http.StatusOK,
	Body: string(jsonBytes),
},nil
}

func main(){
	lambda.Start(Handler)
}