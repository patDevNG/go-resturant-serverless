package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/projects/my_resturant_serverless/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
func Handler(ctx context.Context, req Request) (*Response, error) {
	recordPerPage, err := strconv.Atoi(req.QueryStringParameters["recordPerPage"])
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}
	page, err := strconv.Atoi(req.QueryStringParameters["page"])
	if err != nil || page < 1 {
		page = 1
	}

	startIndex, err := strconv.Atoi(req.QueryStringParameters["startIndex"])
	if err != nil || startIndex < 0 {
		startIndex = 0
	}

	matchStage := bson.D{{"$match", bson.D{{}}}}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}},
		}},
	}
	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"food_item", bson.D{{"$slice", []interface{}{
				"$data", startIndex, recordPerPage,
			}}}},
		}},
	}

	cursor, err := foodCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
	if err != nil {
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error fetching foods: %s", err),
		}, nil
	}

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error decoding cursor results: %s", err),
		}, nil
	}

	jsonBytes, err := json.Marshal(result[0])
	if err != nil {
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error encoding response: %s", err),
		}, nil
	}

	return &Response{
		StatusCode: http.StatusOK,
		Body:       string(jsonBytes),
	}, nil
}


func main() {
	lambda.Start(Handler)
}
