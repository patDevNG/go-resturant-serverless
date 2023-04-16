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
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, _ events.APIGatewayProxyRequest)(*events.APIGatewayProxyResponse, error){

	res,err := menuCollection.Find(context.TODO(), bson.M{})
	if err != nil{
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadGateway,
			Body: fmt.Sprintf("error fetching menus"),
		},nil
	}
  var allMenu []bson.M
	if err = res.All(ctx, &allMenu); err != nil {
		log.Fatal(err)
	}
	jsonbytes, err := json.Marshal(allMenu)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadGateway,
			Body: fmt.Sprintf("error fetching menus"),
		},nil
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body: string(jsonbytes),
	},nil
}

func main()  {
	lambda.Start(Handler)
}