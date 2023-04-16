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

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
func Handler(ctx context.Context, req Request ) (*Response, error) {
var food models.Food
foodId := req.PathParameters["foodId"]

err:= foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
if err != nil {
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Body: fmt.Sprintf("food not found"),
	},nil
}
jsonBytes, err := json.Marshal(food)
if err !=nil{
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Body: fmt.Sprintf("Error Fetching foods"),
	},nil
}
return &Response{
StatusCode: http.StatusOK,
Body: string(jsonBytes),
},nil
}

func main() {
	lambda.Start(Handler)
}
