package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest
var validate = validator.New()
var orderCollection *mongo.Collection = database.OpenCollection(database.Client,"order")
var orderItemCollection  *mongo.Collection = database.OpenCollection(database.Client, "OrderItem")
type OrderItemPack struct {
	Table_id *string
	Order_items []models.OrderItems
}

func Handler(ctx context.Context, req Request)(*Response, error){
var orderItemPack OrderItemPack
var order models.Order

if err := json.Unmarshal([]byte(req.Body), &orderItemPack); err !=nil{
	return &Response{
		StatusCode: http.StatusBadRequest,
		Body: "An error occured check input",
	},nil
}
order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
orderItemsToBeInserted := []interface{}{}
order_id := OrderItemOrderCreator(order)
for _, orderItem :=range orderItemPack.Order_items{
	orderItem.Order_id = order_id

	validationErr := validate.Struct(orderItem)
	if validationErr != nil {
		return &Response{
			StatusCode: http.StatusBadRequest,
			Body: fmt.Sprintf("An error occured during input validation"),
		},nil
	}
	orderItem.ID = primitive.NewObjectID()
	orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)) 
	orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderItem.Order_items_id = orderItem.ID.Hex()
	var num = toFixed(*orderItem.Unit_price, 2)
	orderItem.Unit_price = &num
	orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
}
insertedOrderItems, err := orderItemCollection.InsertMany(ctx,orderItemsToBeInserted)
if err != nil{
	log.Fatal(err)
}
jsonBytes, err := json.Marshal(insertedOrderItems)
if err != nil {
	log.Fatal(err)
}
return &Response{
	StatusCode: http.StatusOK,
	Body: string(jsonBytes),
},nil
}

func OrderItemOrderCreator(order models.Order) string {
	var	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()
	
	orderCollection.InsertOne(ctx, order)
	defer cancel()
	return order.Order_id
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
