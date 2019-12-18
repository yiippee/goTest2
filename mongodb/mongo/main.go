package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
	"time"
)

type Info struct {
	Id   int64
	Time time.Time
	Lat  float64
	Long float64
	info string
}

func Image() {
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://172.20.200.17:40000"))
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	coll := client.Database("test").Collection("log")

	{
		// 插入结构体测试
		info := Info{
			Id:   999,
			Time: time.Now(),
			Lat:  23.33,
			Long: 116.43,
			info: "test",
		}

		result, err := coll.InsertOne(context.Background(), info) // 可直接orm

		if err != nil {
			panic(err)
		}
		fmt.Println(result.InsertedID)
		//
		info2 := Info{
			Id: 999,
		}
		cursor, err := coll.Find(context.Background(), info2)

		for cursor.Next(context.Background()) {
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)
		}
	}
	{
		// Start Example 1

		result, err := coll.InsertOne(
			context.Background(),
			bson.D{
				{"id", 3333333333333},
				{"time", 24345354},
				{"planeID", "123456"},
				{"lat", 22.6},
				{"long", 116.8},
				{"info", bson.D{
					{"imageUrl", "./Images/123.jpg"},
					{"h", 28},
					{"w", 35.5},
					{"point", bson.D{
						{"x", 123},
						{"y", 456},
					}},
				}},
			})
		if err != nil {
			panic(err)
		}
		fmt.Println(result.InsertedID)
		//
		cursor, err := coll.Find(context.Background(),
			bson.D{{"planeID", "123456"}})

		for cursor.Next(context.Background()) {
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)
		}
	}
}
func main() {
	Image()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	// export excel test
	//coll := client.Database("excel").Collection("stock")
	//result, err := coll.InsertOne(
	//	context.Background(),
	//	bson.D{
	//		{"item", "canvas"},
	//		{"qty", 100},
	//		{"tags", bson.A{"cotton"}},
	//		{"size", bson.D{
	//			{"h", 28},
	//			{"w", 35.5},
	//			{"uom", "cm"},
	//		}},
	//	})

	//
	db := client.Database("myTest")
	t := new(testing.T)

	UpdateExamples(t, db)

	// QueryToplevelFieldsExamples(t, db)
	// InsertExamples(t, db)

	collection := client.Database("blog").Collection("student")

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(result)
		// do something with result....
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	var result struct {
		Sage int
	}
	filter := bson.M{"sname": "lisi"}
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
}
func requireCursorLength(t *testing.T, cursor *mongo.Cursor, length int) {
	i := 0
	for cursor.Next(context.Background()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
		i++
	}

	require.NoError(t, cursor.Err())
	require.Equal(t, i, length)
}

func containsKey(doc bson.Raw, key ...string) bool {
	_, err := doc.LookupErr(key...)
	if err != nil {
		return false
	}
	return true
}

// InsertExamples contains examples for insert operations.
func InsertExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_insert")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		// Start Example 1

		result, err := coll.InsertOne(
			context.Background(),
			bson.D{
				{"item", "canvas"},
				{"qty", 100},
				{"tags", bson.A{"cotton"}},
				{"size", bson.D{
					{"h", 28},
					{"w", 35.5},
					{"uom", "cm"},
				}},
			})

		// End Example 1

		require.NoError(t, err)
		require.NotNil(t, result.InsertedID)
	}

	{
		// Start Example 2

		cursor, err := coll.Find(
			context.Background(),
			bson.D{{"item", "canvas"}},
		)

		// End Example 2

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)

	}

	{
		// Start Example 3

		result, err := coll.InsertMany(
			context.Background(),
			[]interface{}{
				bson.D{
					{"item", "journal"},
					{"qty", int32(25)},
					{"tags", bson.A{"blank", "red"}},
					{"size", bson.D{
						{"h", 14},
						{"w", 21},
						{"uom", "cm"},
					}},
				},
				bson.D{
					{"item", "mat"},
					{"qty", int32(25)},
					{"tags", bson.A{"gray"}},
					{"size", bson.D{
						{"h", 27.9},
						{"w", 35.5},
						{"uom", "cm"},
					}},
				},
				bson.D{
					{"item", "mousepad"},
					{"qty", 25},
					{"tags", bson.A{"gel", "blue"}},
					{"size", bson.D{
						{"h", 19},
						{"w", 22.85},
						{"uom", "cm"},
					}},
				},
			})

		// End Example 3

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 3)
	}
}

// QueryToplevelFieldsExamples contains examples for querying top-level fields.
func QueryToplevelFieldsExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_query_top")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		// Start Example 6

		docs := []interface{}{
			bson.D{
				{"item", "journal"},
				{"qty", 25},
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "notebook"},
				{"qty", 50},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "paper"},
				{"qty", 100},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "planner"},
				{"qty", 75},
				{"size", bson.D{
					{"h", 22.85},
					{"w", 30},
					{"uom", "cm"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "postcard"},
				{"qty", 45},
				{"size", bson.D{
					{"h", 10},
					{"w", 15.25},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
		}

		result, err := coll.InsertMany(context.Background(), docs)

		// End Example 6

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 5)
	}

	{
		// Start Example 7

		cursor, err := coll.Find(
			context.Background(),
			bson.D{},
		)

		// End Example 7

		require.NoError(t, err)
		requireCursorLength(t, cursor, 5)
	}

	{
		// Start Example 9

		cursor, err := coll.Find(
			context.Background(),
			bson.D{{"status", "D"}},
		)

		// End Example 9

		require.NoError(t, err)
		requireCursorLength(t, cursor, 2)
	}

	{
		// Start Example 10

		cursor, err := coll.Find(
			context.Background(),
			bson.D{{"status", bson.D{{"$in", bson.A{"A", "D"}}}}})

		// End Example 10

		require.NoError(t, err)
		requireCursorLength(t, cursor, 5)
	}

	{
		// Start Example 11

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"status", "A"},
				{"qty", bson.D{{"$lt", 30}}},
			})

		// End Example 11

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		// Start Example 12

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"$or",
					bson.A{
						bson.D{{"status", "A"}},
						bson.D{{"qty", bson.D{{"$lt", 30}}}},
					}},
			})

		// End Example 12

		require.NoError(t, err)
		requireCursorLength(t, cursor, 3)
	}

	{
		// Start Example 13

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"status", "A"},
				{"$or", bson.A{
					bson.D{{"qty", bson.D{{"$lt", 30}}}},
					bson.D{{"item", primitive.Regex{Pattern: "^p", Options: ""}}},
				}},
			})

		// End Example 13

		require.NoError(t, err)
		requireCursorLength(t, cursor, 2)
	}

}

// UpdateExamples contains examples of update operations.
func UpdateExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_update")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		// Start Example 51

		docs := []interface{}{
			bson.D{
				{"item", "canvas"},
				{"qty", 100},
				{"size", bson.D{
					{"h", 28},
					{"w", 35.5},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "journal"},
				{"qty", 25},
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "mat"},
				{"qty", 85},
				{"size", bson.D{
					{"h", 27.9},
					{"w", 35.5},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "mousepad"},
				{"qty", 25},
				{"size", bson.D{
					{"h", 19},
					{"w", 22.85},
					{"uom", "in"},
				}},
				{"status", "P"},
			},
			bson.D{
				{"item", "notebook"},
				{"qty", 50},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "P"},
			},
			bson.D{
				{"item", "paper"},
				{"qty", 100},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "D"},
				{"array", bson.A{1, 2, 3}},
			},
			bson.D{
				{"item", "planner"},
				{"qty", 75},
				{"size", bson.D{
					{"h", 22.85},
					{"w", 30},
					{"uom", "cm"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "postcard"},
				{"qty", 45},
				{"size", bson.D{
					{"h", 10},
					{"w", 15.25},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "sketchbook"},
				{"qty", 80},
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "sketch pad"},
				{"qty", 95},
				{"size", bson.D{
					{"h", 22.85},
					{"w", 30.5},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
		}

		result, err := coll.InsertMany(context.Background(), docs)

		// End Example 51

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 10)
	}

	{
		// Start Example 52

		result, err := coll.UpdateOne(
			context.Background(),
			bson.D{
				{"item", "paper"},
			},
			bson.D{
				{"$set", bson.D{
					{"size.uom", "cm"},
					{"status", "P"},
				}},
				{"$currentDate", bson.D{
					{"lastModified", true},
				}},
			},
		)

		// update push array
		result, err = coll.UpdateOne(
			context.Background(),
			bson.D{
				{"item", "paper"},
			},
			bson.D{ // bson.D 代表一个对象， bson.A 代表一个数组
				{"$push", bson.D{{"array", bson.D{{"$each", bson.A{4, 5, 6}}}}}},
				{"$currentDate", bson.D{
					{"lastModified", true},
				}},
			},
		)
		// End Example 52

		require.NoError(t, err)
		require.Equal(t, int64(1), result.MatchedCount)
		require.Equal(t, int64(1), result.ModifiedCount)

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"item", "paper"},
			})

		require.NoError(t, err)

		for cursor.Next(context.Background()) {
			doc := cursor.Current

			uom, err := doc.LookupErr("size", "uom")
			require.NoError(t, err)
			require.Equal(t, uom.StringValue(), "cm")

			status, err := doc.LookupErr("status")
			require.NoError(t, err)
			require.Equal(t, status.StringValue(), "P")

			require.True(t, containsKey(doc, "lastModified"))
		}

		require.NoError(t, cursor.Err())
	}

	{
		// Start Example 53

		result, err := coll.UpdateMany(
			context.Background(),
			bson.D{
				{"qty", bson.D{
					{"$lt", 50},
				}},
			},
			bson.D{
				{"$set", bson.D{
					{"size.uom", "cm"},
					{"status", "P"},
				}},
				{"$currentDate", bson.D{
					{"lastModified", true},
				}},
			},
		)

		// End Example 53

		require.NoError(t, err)
		require.Equal(t, int64(3), result.MatchedCount)
		require.Equal(t, int64(3), result.ModifiedCount)

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"qty", bson.D{
					{"$lt", 50},
				}},
			})

		require.NoError(t, err)

		for cursor.Next(context.Background()) {
			doc := cursor.Current

			uom, err := doc.LookupErr("size", "uom")
			require.NoError(t, err)
			require.Equal(t, uom.StringValue(), "cm")

			status, err := doc.LookupErr("status")
			require.NoError(t, err)
			require.Equal(t, status.StringValue(), "P")

			require.True(t, containsKey(doc, "lastModified"))
		}

		require.NoError(t, cursor.Err())
	}

	{
		// Start Example 54

		result, err := coll.ReplaceOne(
			context.Background(),
			bson.D{
				{"item", "paper"},
			},
			bson.D{
				{"item", "paper"},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "A"},
						{"qty", 60},
					},
					bson.D{
						{"warehouse", "B"},
						{"qty", 40},
					},
				}},
			},
		)

		// End Example 54

		require.NoError(t, err)
		require.Equal(t, int64(1), result.MatchedCount)
		require.Equal(t, int64(1), result.ModifiedCount)

		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"item", "paper"},
			})

		require.NoError(t, err)

		for cursor.Next(context.Background()) {
			require.True(t, containsKey(cursor.Current, "_id"))
			require.True(t, containsKey(cursor.Current, "item"))
			require.True(t, containsKey(cursor.Current, "instock"))

			instock, err := cursor.Current.LookupErr("instock")
			require.NoError(t, err)
			vals, err := instock.Array().Values()
			require.NoError(t, err)
			require.Equal(t, len(vals), 2)

		}

		require.NoError(t, cursor.Err())
	}

}
