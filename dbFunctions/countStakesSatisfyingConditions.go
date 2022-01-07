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
	"strconv"
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

func GetCombinationStats(combination producerCode.Combination, collection *mongo.Collection) map[string]interface{} {
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

	//filter := bson.M{
	//	"Pairname":  bson.M{"$in": pairnameArray},
	//	"Timeframe": bson.M{"$in": timeframeArray},
	//}

	filter := makeFilter(combination)

	fmt.Println("filter = ", filter)

	makeFilter(combination)

	// Запрос без фильтров (при нежелании использовать какой-либо фильтр надо передавать
	// bson.D{{}}, а не nil, иначе будет ошибка)
	//stakesCount, err := collection.CountDocuments(context.TODO(), bson.D{{}})
	stakesCount, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Fatalln("err = ", err)
	}

	//filterEndPriceMoreThanInitial := bson.M{
	//	"$and": []bson.M{
	//		{"Pairname": "EUR/JPY"},
	//		{"Timeframe": 300},
	//		{"endPriceMoreThanInitial": 1},
	//	},
	//}

	filterEndPriceMoreThanInitial := filter
	filterEndPriceMoreThanInitial["endPriceMoreThanInitial"] = 1

	stakesWhereEndPriceMoreThanInitialCount, err := collection.CountDocuments(context.TODO(), filterEndPriceMoreThanInitial)
	if err != nil {
		log.Fatalln("err = ", err)
	}

	fmt.Println("stakesCount = ", stakesCount)
	fmt.Println("stakesWhereEndPriceMoreThanInitialCount = ", stakesWhereEndPriceMoreThanInitialCount)
	percentOfStakesWhereEndPriceMoreThanInitial := (float64(stakesWhereEndPriceMoreThanInitialCount) / float64(stakesCount)) * 100
	fmt.Printf("priceMore / allStakes = %.1f\n", percentOfStakesWhereEndPriceMoreThanInitial)

	stats := make(map[string]interface{})
	stats["Combination"] = combination
	stats["Stakes at all"] = stakesCount
	stats["Stakes with end price more than initial"] = stakesWhereEndPriceMoreThanInitialCount
	stats["Percent of stakes with end price more than initial"] = fmt.Sprintf("%.2f", percentOfStakesWhereEndPriceMoreThanInitial)

	return stats
}

func makeFilter(combination producerCode.Combination) bson.M {
	pairnames := getPairnamesFromCombination(combination)
	timeframes := getTimeframesFromCombination(combination)

	//fmt.Println("[makeFilter] pairnames = ", pairnames)
	//fmt.Println("[makeFilter] timeframes = ", timeframes)

	filter := bson.M{
		//"Pairname":  bson.M{"$in": pairnames},
		//"Timeframe": bson.M{"$in": timeframes},
	}

	if len(pairnames) != 0 {
		filter["Pairname"] = bson.M{"$in": pairnames}
	}

	if len(timeframes) != 0 {
		filter["Timeframe"] = bson.M{"$in": timeframes}
	}

	return filter
}

func getPairnamesFromCombination(combination producerCode.Combination) []string {
	var pairnames []string
	for _, condition := range combination.Conditions {
		if condition.ColumnName == "Pairname" {
			pairnames = append(pairnames, condition.Value)
		}
	}
	return pairnames
}

func getTimeframesFromCombination(combination producerCode.Combination) []int {
	var timeframes []int
	for _, condition := range combination.Conditions {
		if condition.ColumnName == "Timeframe" {
			// Конвертируем в int, так как Timeframe в БД хранится как числовое значение, и в фильтре
			// тоже должно быть числовое значение
			valueAsInt, _ := strconv.Atoi(condition.Value)

			timeframes = append(timeframes, valueAsInt)
		}
	}
	return timeframes
}
