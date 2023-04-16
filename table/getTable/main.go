package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/projects/my_resturant_serverless/database"
	"github.com/projects/my_resturant_serverless/models"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
)


var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")
func Handler(ctx context.Context, req events.APIGatewayProxyRequest )(*events.APIGatewayProxyResponse, error){
var table models.Table

tableId := req.PathParameters["tableId"]
err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
if err !=nil{
	msg := fmt.Sprintf("Error fetching menu")
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadGateway,
		Body: msg,
	},nil
}
jsonBytes, err := json.Marshal(table)

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