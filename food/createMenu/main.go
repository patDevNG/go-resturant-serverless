package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"github.com/projects/my_resturant_serverless/database"
	"github.com/projects/my_resturant_serverless/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
 var validate =  validator.New()
var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, req events.APIGatewayProxyRequest)(*events.APIGatewayProxyResponse, error){
	var menu models.Menu
	err := json.Unmarshal([]byte(req.Body), &menu)
	if err != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	}

	validationErr := validate.Struct(menu)
	if validationErr != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, validationErr
	}
	menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.ID = primitive.NewObjectID()
	menu.Menu_id = menu.ID.Hex()
	res, insertErr := menuCollection.InsertOne(ctx, menu)

	if insertErr!=nil{
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadGateway,
			Body: fmt.Sprintf("Menu was not created"),
		}, insertErr
	}
	jsonBytes, err := json.Marshal(res);

	if err != nil{
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadGateway,
		},err
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body: string(jsonBytes),
	},nil
}

func main() {
	lambda.Start(Handler)
}