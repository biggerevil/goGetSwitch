package dbFunctions

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"goGetSwitch/signal"
	"log"
)

func UpdateEndingStakes(currentUnixTimestamp int64, allSignals []signal.Signal) {
	log.Println("[UpdateEndingStakes] Начало")

	// ====1=========================================================
	// Подключаемся к БД
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	log.Println("Собираюсь сделать Ping")
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("Закончил с вызовом Ping")
	fmt.Println("Successfully connected and pinged.")

	collectionOfStakes := client.Database("history_stakes").Collection("stakes")
	// ====1=========================================================

	// ====2=========================================================
	// Отбираем все сигналы, где UnixTimestamp == текущий timestamp
	log.Println("[UpdateEndingStakes] Отбираем старые сигналы, которые сейчас заканчиваются")
	filter := bson.D{{"EndUnixTimestamp", currentUnixTimestamp}}
	findOptions := options.Find()
	cur, err := collectionOfStakes.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	// ====2=========================================================

	// ====3=========================================================
	log.Println("[UpdateEndingStakes] Начинаем итерацию по каждому старому сигналу")

	// Создаём словарь, в котором будем хранить пары типа
	// "НазваниеВалютнойПары":ЗначениеВалютнойПарыНаДанныйМомент
	var memorizedPairsAndPrices map[string]float64

	// Итерируемся по каждому сигналу, у которого UnixTimestamp == текущий timestamp
	for cur.Next(context.TODO()) {
		// Достаём сигнал
		var oldSignal signal.Signal
		err := cur.Decode(&oldSignal)
		if err != nil {
			fmt.Println("[UpdateMultipleStakesOutcomes] Ошибка при попытке Find. #2")
			log.Fatal(err)
		}

		log.Println("[UpdateEndingStakes] Получили сигнал с ID = ", oldSignal.ID)

		// 3.1 и 3.2.
		// Находим текущую цену сигнала
		var currentPrice float64
		requiredPairname := oldSignal.Pairname
		if memorizedPrice, ok := memorizedPairsAndPrices[requiredPairname]; ok {
			currentPrice = memorizedPrice
		} else {
			currentPrice = findPriceOfPairnameInAllSignals(requiredPairname, allSignals)
		}

		// 3.3
		// Определяем, больше ли цена, чем предыдущая
		currentPriceBiggerThanPreviousPrice := currentPriceBiggerThanPreviousPriceFunction(currentPrice, oldSignal.CurrentPrice)
		//_ = currentPriceBiggerThanPreviousPriceFunction(currentPrice, oldSignal.CurrentPrice)

		// 3.4
		// Делаем update сигнала, добавляя новую цену и рез-т сделки
		updateStakeInDB(collectionOfStakes, oldSignal.ID, currentPrice, currentPriceBiggerThanPreviousPrice)
	}
	// ====3=========================================================

	log.Println("[UpdateEndingStakes] Конец")
}

func findPriceOfPairnameInAllSignals(requiredPairname string, allSignals []signal.Signal) float64 {
	for _, signal := range allSignals {
		if signal.Pairname == requiredPairname {
			return signal.CurrentPrice
		}
	}

	log.Fatalln("[findPriceOfPairnameInAllSignals] Не нашёл сигнал с парой = ", requiredPairname, "среди всех сигналов, делаю panic()")
	panic("panic that shouldn't ever work because of previous log.Fatalln, right?")
}

func currentPriceBiggerThanPreviousPriceFunction(currentPrice float64, previousPrice float64) int {
	if currentPrice > previousPrice {
		return 1
	}

	if currentPrice < previousPrice {
		return -1
	}

	if currentPrice == previousPrice {
		return 0
	}

	panic("[currentPriceBiggerThanPreviousPrice] Код никогда не должен доходить досюда!")
}

func updateStakeInDB(stakesCollection *mongo.Collection, idOfsignaltoupdate primitive.ObjectID, endPrice float64, endPriceBiggerThanInitial int) {
	result, err := stakesCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": idOfsignaltoupdate},
		bson.D{
			{"$set", bson.D{{"endPrice", endPrice}}},
			{"$set", bson.D{{"endPriceMoreThanInitial", endPriceBiggerThanInitial}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
}
