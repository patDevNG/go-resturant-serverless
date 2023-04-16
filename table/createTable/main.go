package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"github.com/projects/my_resturant_serverless/database"
	"github.com/projects/my_resturant_serverless/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
var validate = validator.New()
var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")
type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, req events.APIGatewayProxyRequest)(*events.APIGatewayProxyResponse, error){
	var table models.Table

	err:= json.Unmarshal([]byte(req.Body), &table)
	if err != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	}

	validationError := validate.Struct(table)
	if validationError !=nil{
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, validationError
	}
	table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.ID = primitive.NewObjectID()
  table.Table_id = table.ID.Hex()
	res, insertErr := tableCollection.InsertOne(ctx, table)

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

func main(){
	lambda.Start(Handler)
}