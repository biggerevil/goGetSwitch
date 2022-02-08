package dbFunctions

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"goGetSwitch/producerCode"
	"goGetSwitch/stats"
	"log"
	"strconv"
)

// Это названия полей я сохраняю в отдельных переменных. Таким образом хочу защититься от того, что я где-то
// опечатаюсь в названии поля.
const maBuyColumnName = "MaBuy"
const maSellColumnName = "MaSell"
const tiBuyColumnName = "TiBuy"
const tiSellColumnName = "TiSell"

/*
	Это отдельная функция, которая подключается к БД и возвращает коллекцию.
	Вынес это в отдельную ф-ю, чтобы не писать такой код в начале каждой ф-и, которой надо подключиться к БД.
*/
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

/*
	Эта функция принимает на вход комбинацию и коллекцию. Возвращает она статистику по переданной комбинации.
	TODO: Считать винрейт не относительно всех ставок, а относительно тех ставок, где endPriceMoreThanInitial. Так
	 как в некоторых комбинациях может быть куча ставок, где endPriceMoreThanInitial=0. И там на самом деле винрейт
	 не 33%, а как бы 50% (если не считать endPriceMoreThanInitial=0. Но в таком случае брокер по идее делает просто
	 возврат ставки, но там кто как, не говоря уже о задержке сети при отправке ставки. Но в любом случае, я думаю,
	 так считать винрейт будет более корректным решением).
	TODO: привести функцию к красивому виду, удалить ненужный код, добавить комментарии к тому, что она делает.
*/
func GetCombinationStats(combination producerCode.Combination, collection *mongo.Collection) stats.Stats {
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

	//filter := bson.M{
	//	"Pairname":  bson.M{"$in": pairnameArray},
	//	"Timeframe": bson.M{"$in": timeframeArray},
	//}

	filter := makeFilter(combination)

	fmt.Println("filter = ", filter)

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

	//stats := make(map[string]interface{})
	//stats["Combination"] = combination
	//stats["Stakes at all"] = stakesCount
	//stats["Stakes with end price more than initial"] = stakesWhereEndPriceMoreThanInitialCount
	//stats["Percent of stakes with end price more than initial"] = fmt.Sprintf("%.2f", percentOfStakesWhereEndPriceMoreThanInitial)

	stats := stats.Stats{}
	stats.Combination = combination
	stats.StakesAtAll = stakesCount
	stats.AllRelativeStakesWhereEndPriceMoreThanInitialCount = stakesWhereEndPriceMoreThanInitialCount
	stats.AllRelativePercentOfStakesWhereEndPriceMoreThanInitial = percentOfStakesWhereEndPriceMoreThanInitial

	return stats
}

/*
	Эта функция принимает на вход комбинацию и возвращает filter для MongoDB.
	При помощи этого фильтра можно будет достать из БД только те ставки, которые подходят под эту комбинацию.
*/
func makeFilter(combination producerCode.Combination) bson.M {
	// 1. Сначала мы по очереди достаём из комбинации все пары, все таймфреймы, все индикаторы.
	pairnames := getPairnamesFromCombination(combination)
	timeframes := getTimeframesFromCombination(combination)
	maBuys := getIndicatorFromCombination(combination, maBuyColumnName)
	maSells := getIndicatorFromCombination(combination, maSellColumnName)
	tiBuys := getIndicatorFromCombination(combination, tiBuyColumnName)
	tiSells := getIndicatorFromCombination(combination, tiSellColumnName)

	//fmt.Println("[makeFilter] pairnames = ", pairnames)
	//fmt.Println("[makeFilter] timeframes = ", timeframes)

	// 2. Затем мы создаём фильтр (который в итоге вернём), в который будем
	// добавлять условия (Condition) из комбинации (Combination).
	filter := bson.M{
		//"Pairname":  bson.M{"$in": pairnames},
		//"Timeframe": bson.M{"$in": timeframes},
	}

	// 3. И затем по очереди добавляем в наш фильтр все условия (Condition) из комбинации (Combination).
	// Если эти условия (Condition) вообще есть, конечно (поэтому сначала проверяем при помощи if len(...) != 0)
	if len(pairnames) != 0 {
		filter["Pairname"] = bson.M{"$in": pairnames}
	}

	if len(timeframes) != 0 {
		filter["Timeframe"] = bson.M{"$in": timeframes}
	}

	if len(maBuys) != 0 {
		filter[maBuyColumnName] = bson.M{"$in": maBuys}
	}

	if len(maSells) != 0 {
		filter[maSellColumnName] = bson.M{"$in": maSells}
	}

	if len(tiBuys) != 0 {
		filter[tiBuyColumnName] = bson.M{"$in": tiBuys}
	}

	if len(tiSells) != 0 {
		filter[tiSellColumnName] = bson.M{"$in": tiSells}
	}

	return filter
}

/*
	При помощи этой функции мы достаём все пары из комбинации.
*/
func getPairnamesFromCombination(combination producerCode.Combination) []string {
	var pairnames []string
	for _, condition := range combination.Conditions {
		if condition.ColumnName == "Pairname" {
			pairnames = append(pairnames, condition.Value)
		}
	}
	return pairnames
}

/*
	При помощи этой функции мы достаём все таймфреймы из комбинации.
	В отличии от доставания пары таймфрейм мы превращаем в int (а пару мы превращаем в string).
*/
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

/*
	При помощи этой функции мы достаём все индикаторы из комбинации.
	При этом поскольку у нас 4 вида индикаторов, но по сути они похожи, то мы дополнительно передаём в ф-ю
	название нашего индикатора.
*/
func getIndicatorFromCombination(combination producerCode.Combination, indicatorColumnName string) []int {
	var indicators []int
	for _, condition := range combination.Conditions {
		if condition.ColumnName == indicatorColumnName {
			// Конвертируем в int, так как condition.Value по умолчанию является string
			valueAsInt, _ := strconv.Atoi(condition.Value)

			indicators = append(indicators, valueAsInt)
		}
	}
	return indicators
}
