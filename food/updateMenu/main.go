package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/projects/my_resturant_serverless/database"
	"github.com/projects/my_resturant_serverless/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, req events.APIGatewayProxyRequest)(*events.APIGatewayProxyResponse, error){
	var menu models.Menu
	menuId := req.PathParameters["menuId"]
	filter := bson.M{"menu_id": menuId}
	err := json.Unmarshal([]byte(req.Body), &menu)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body: fmt.Sprintf("check request"),
		}, nil
	}
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{"start_date", menu.Start_Date})
	updateObj = append(updateObj, bson.E{"end_date", menu.End_Date})
	if menu.Name !=""{
		updateObj = append(updateObj, bson.E{"name", menu.Name})
	}

	if menu.Category != ""{
		updateObj = append(updateObj, bson.E{"category", menu.Category})
	}

	menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", menu.Updated_at})

	upsert := true
	opt:= options.UpdateOptions{
		Upsert: &upsert,
	}
	update := bson.D{{"$set", updateObj}}

	res, err := menuCollection.UpdateOne(
		ctx,
		filter,
		update,
		&opt,
	)
	log.Fatal(res, "result")
	if err != nil{
		msg:= fmt.Sprintf("Error occured while updating")
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadGateway,
			Body: msg,
		},err
	}
	jsonBytes, err := json.Marshal(res)

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

func main()  {
	lambda.Start(Handler)
}