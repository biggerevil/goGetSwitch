package dbFunctions

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"goGetSwitch/producerCode"
	"log"
)

func ConnectToDB() *mongo.Collection {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	// Ping the primary
	log.Println("Собираюсь сделать Ping")
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("Закончил с вызовом Ping")

	fmt.Println("Successfully connected and pinged.")

	collection := client.Database("history_stakes").Collection("stakes")

	return collection
}

func GetCombinationStats(combination producerCode.Combination, collection *mongo.Collection) {
	// CORRECT in its own way. Just doesn't enough, because only one timeframe
	//filter := bson.M{
	//	"$and": []bson.M{
	//		bson.M{"$or": []bson.M{
	//			bson.M{"Pairname": "EUR/USD"},
	//			bson.M{"Pairname": "EUR/JPY"},
	//		}},
	//		{"Timeframe": 900},
	//	}}

	// Запрос для mongoshell с такими же критериями, как в этом фильтре:
	// db.stakes.find({$and: [{$or : [{"Pairname":"EUR/JPY"},{"Pairname":"EUR/USD"}]},{$or : [{"Timeframe":900}]}] }).count()
	//filter := bson.M{
	//	"$and": []bson.M{
	//		bson.M{"$or": []bson.M{
	//			bson.M{"Pairname": "EUR/USD"},
	//			bson.M{"Pairname": "EUR/JPY"},
	//		}},
	//		bson.M{"$or": []bson.M{
	//			bson.M{"Timeframe": 300},
	//			bson.M{"Timeframe": 900},
	//		}},
	//	}}

	var pairnameArray []string
	pairnameArray = append(pairnameArray, "EUR/JPY")
	pairnameArray = append(pairnameArray, "EUR/USD")
	//pairnameArray = append(pairnameArray, "USD/JPY")

	var timeframeArray []int
	timeframeArray = append(timeframeArray, 300)
	//timeframeArray = append(timeframeArray, 900)
	//timeframeArray = append(timeframeArray, 1800)

	filter := bson.M{
		"Pairname":  bson.M{"$in": pairnameArray},
		"Timeframe": bson.M{"$in": timeframeArray},
	}

	fmt.Println("filter = ", filter)

	// Запрос без фильтров (при нежелании использовать какой-либо фильтр надо передавать
	// bson.D{{}}, а не nil, иначе будет ошибка)
	//stakesCount, err := collection.CountDocuments(context.TODO(), bson.D{{}})
	stakesCount, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Fatalln("err = ", err)
	}

	filterEndPriceMoreThanInitial := bson.M{
		"$and": []bson.M{
			{"Pairname": "EUR/JPY"},
			{"Timeframe": 300},
			{"endPriceMoreThanInitial": 1},
		},
	}
	stakesWhereEndPriceMoreThanInitialCount, err := collection.CountDocuments(context.TODO(), filterEndPriceMoreThanInitial)
	if err != nil {
		log.Fatalln("err = ", err)
	}

	fmt.Println("stakesCount = ", stakesCount)
	fmt.Println("stakesWhereEndPriceMoreThanInitialCount = ", stakesWhereEndPriceMoreThanInitialCount)
	percentOfStakesWhereEndPriceMoreThanInitial := (float64(stakesWhereEndPriceMoreThanInitialCount) / float64(stakesCount)) * 100
	fmt.Printf("priceMore / allStakes = %.1f\n", percentOfStakesWhereEndPriceMoreThanInitial)
}
