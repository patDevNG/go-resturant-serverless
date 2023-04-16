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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate =  validator.New()
var orderCollection *mongo.Collection = database.OpenCollection(database.Client,"order")
var tableCollection *mongo.Collection = database.OpenCollection(database.Client,"order")
type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest
func Handler(ctx context.Context, req Request)(*Response, error){
	var table models.Table
	var order models.Order

	err := json.Unmarshal([]byte(req.Body),&order)
	if err !=nil{
		return &Response{
			StatusCode: http.StatusBadRequest,
			Body: fmt.Sprintf("Invalid input"),
		}, nil
	}
	validationErr := validate.Struct(order)
	if validationErr != nil {
		return &Response{StatusCode: http.StatusBadRequest}, validationErr
	}
	if order.Table_id != nil {
		err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode((&table))

		if err != nil {
			return &Response{
				StatusCode: http.StatusNotFound,
				Body: fmt.Sprintf("Table not found"),
			},nil
		}
	}
	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	res, insertErr := orderCollection.InsertOne(ctx, order)

	if insertErr != nil {
		msg:= fmt.Sprintf("order not created")
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Body: msg,
		},nil
	}

	jsonBytes, err:= json.Marshal(res)

	if err!=nil{
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