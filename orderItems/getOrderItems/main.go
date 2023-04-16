package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/projects/my_resturant_serverless/database"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest
var orderItemCollection  *mongo.Collection = database.OpenCollection(database.Client, "OrderItem")
func Handler(ctx context.Context, _ Request )(*Response, error){
	res, err := orderItemCollection.Find(context.TODO(), bson.M{})

	if err !=nil{
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Body: fmt.Sprintf("Error fetching orderItems"),
		},nil
	}

	var allOrderItem []bson.M
	if err = res.All(ctx, &allOrderItem); err != nil {
		log.Fatal(err)
	}
	jsonBytes, err := json.Marshal(allOrderItem)

	if err !=nil {
		log.Fatal(err)
	}

	return &Response{
		StatusCode: http.StatusOK,
		Body: string(jsonBytes),
	},nil
}

func main () {
	lambda.Start(Handler)
}