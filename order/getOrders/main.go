package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/projects/my_resturant_serverless/database"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)
type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

var orderCollection *mongo.Collection = database.OpenCollection(database.Client,"order")
func Handler(ctx context.Context, _ Request)(*Response, error){
res, err:= orderCollection.Find(context.TODO(), bson.M{})
if err !=nil {
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Body: fmt.Sprintf("Error Fetching order"),
	},nil
}
var allOrder []bson.M
if err = res.All(ctx,&allOrder); err != nil{
	log.Fatal(err)
}
return &Response{},nil
}

func main(){
	lambda.Start(Handler)
}