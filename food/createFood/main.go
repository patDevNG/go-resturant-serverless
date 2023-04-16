package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"github.com/projects/my_resturant_serverless/database"
	"github.com/projects/my_resturant_serverless/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)
var validate =  validator.New()
type Response events.APIGatewayProxyResponse
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
func Handler(ctx context.Context, req events.APIGatewayProxyRequest)(*events.APIGatewayProxyResponse, error){
	var food models.Food
	var menu models.Menu
	err:= json.Unmarshal([]byte(req.Body), &food)

	if err != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	}
	
	validationErr := validate.Struct(food)
	if validationErr != nil {
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, validationErr
	}
	foundMenuErr := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)

	if foundMenuErr !=nil{
		return &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: fmt.Sprintf("menu not found")}, nil
	}
	food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.ID = primitive.NewObjectID()
	food.Food_Id = food.ID.Hex()
	var num = toFixed(*food.Price,2)
		food.Price = &num
		res, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr !=nil{
			return &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: fmt.Sprintf("error creating food")}, nil
		}

		jsonBytes, err := json.Marshal(res)
		
		if err != nil{
			return &events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: fmt.Sprintf("An error occured")}, nil
		}
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body: string(jsonBytes),
		},nil
}


func toFixed(num float64, precision int) float64{
	output := math.Pow(10, float64(precision))
	return float64(round(num *output))/output
} 

func round(num float64) int{
	return int(num + math.Copysign(0.5, num))
	
	}

	func main(){
		lambda.Start(Handler)
	}