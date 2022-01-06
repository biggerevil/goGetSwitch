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
	columnNameOfFirstCondition := combination.Conditions[0].ColumnName
	valueOfFirstCondition := combination.Conditions[0].Value

	//filter := bson.D{{columnNameOfFirstCondition, valueOfFirstCondition}}
	filter := bson.M{
		"$and": []bson.M{
			{"Pairname": "EUR/JPY"},
			{"Timeframe": 300},
		},
	}

	fmt.Println("columnNameOfFirstCondition = ", columnNameOfFirstCondition)
	fmt.Println("valueOfFirstCondition = ", valueOfFirstCondition)
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
