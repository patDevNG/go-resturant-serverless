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

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")
func Handler(ctx context.Context, _ events.APIGatewayProxyRequest)(*events.APIGatewayProxyResponse,error){
	var allTable []bson.M
	res,err := tableCollection.Find(context.TODO(), bson.M{})
	if err != nil{
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadGateway,
			Body: fmt.Sprintf("error fetching tables"),
		},nil
	}
	if err = res.All(ctx, &allTable); err != nil{
		log.Fatal(err)
	}

	jsonBytes, err := json.Marshal(allTable)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadGateway,
			Body: fmt.Sprintf("error fetching table"),
		},nil
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body: string(jsonBytes),
	},nil
}

func main(){
	lambda.Start(Handler)
}