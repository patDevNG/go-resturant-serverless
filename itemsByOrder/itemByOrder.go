package ItemsByOrder

import (
	"context"
	"time"

	"github.com/projects/my_resturant_serverless/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
var orderItemCollection  *mongo.Collection = database.OpenCollection(database.Client, "OrderItem")
func ItemsByOrder(id string)(OrderItems []primitive.M, err error){
	var ctx, cancel = context.WithTimeout(context.Background(),100*time.Second)
	matchStage := bson.D{{"$match", bson.M{"order_id":id}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as","food"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}
	lookupOrderStage := bson.D{{"$lookup", bson.D{{"from","order"},{"localField","order_id"}, {"foreignField","order_id"},{"as","order"}}}}
	unwindOrderStage := bson.D{{"$unwind", bson.D{{"path","$order"},{"preserveNullEmptyArrays",true}}}}
	lookupTableStage := bson.D {{"$lookup", bson.D{{"from", "table"},{"localField", "order.table_id"}, {"foreignField","table_id"}, {"as","table"}}}}
	unwindTableStage := bson.D{{"$unwind", bson.D{{"path","table"}, {"preserveNullAndEmptyArrays",true}}}}

	projectStage := bson.D{
		{
			"$project", bson.D{
			{"id",0},
			{"amount","$food.price"},
			{"food_name", "$food.name"},
			{"food_image","food.food_image"},
			{"total_count",1},
			{"table_number","$table.table_number"},
			{"table_id","table.table_id"},
			{"order_id","$order.order_id"},
			{"price", "$food.price"},
			{"qunatity", 1},
			},
		},
	}
	groupStage:= bson.D{{"$group", bson.D{{"_id", bson.D{{"order_id","$order_id"},{"table_id","$table_id"},{"table_number","$table_number"}}},{"pament_due",bson.D{{"$sum","$amount"}}},{"total_count",bson.D{{"$sum",1}}},{"order_items", bson.D{{"$$push",1}}}}}}
	projectStage2:= bson.D{
		{
			"$project", bson.D{
				{"id",0},
				{"payment_due",1},
				{"total_count",1},
				{"table_number", "$_id.table_number"},
				{"order_items",1},
			},
		},
	}
	res, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})
	if err !=nil{
		panic(err)
	}
	if err = res.All(ctx, &OrderItems); err != nil{
		panic(err)
	}
	defer cancel()
	return OrderItems, err
}