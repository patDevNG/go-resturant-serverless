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

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, req events.APIGatewayProxyRequest)(*events.APIGatewayProxyResponse, error){
	var menu models.Menu
	menuId := req.PathParameters["menuId"]
	err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
	if err !=nil{
		msg := fmt.Sprintf("Error fetching menu")
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadGateway,
			Body: msg,
		},nil
	}
	jsonBytes, err := json.Marshal(menu)

	if err != nil{
		msg := fmt.Sprintf("Error fetching menu")
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadGateway,
			Body: msg,
		},err
	}
	return &events.APIGatewayProxyResponse{	
	StatusCode: http.StatusOK,
	Body: string (jsonBytes),
	}, nil
}

func main()  {
	lambda.Start(Handler)
}