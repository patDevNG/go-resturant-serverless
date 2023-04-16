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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")
func Handler(ctx context.Context, req events.APIGatewayProxyRequest)(*events.APIGatewayProxyResponse, error){
var table models.Table

tableId := req.PathParameters["tableId"]
filter:= bson.M{"table_id": tableId}
var updateObj primitive.D

		if table.Number_of_guest != nil {
			updateObj = append(updateObj, bson.E{"number_of_guest", table.Number_of_guest})
		}

		if table.Table_number != nil {
			updateObj = append(updateObj, bson.E{"table_number", table.Table_number})
		}
		upsert :=true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		res, updateErr := tableCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set",updateObj},
			},
			&opt,
		)
		if updateErr != nil {
			log.Fatal(updateErr)
			msg:=fmt.Sprintf("table update failed")
			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body: msg,
			},nil
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