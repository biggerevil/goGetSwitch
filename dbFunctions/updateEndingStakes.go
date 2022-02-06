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
	// Отбираем все сигналы, где EndUnixTimestamp == текущий timestamp
	log.Println("[UpdateEndingStakes] Отбираем старые сигналы, которые сейчас заканчиваются")
	filter := bson.D{{"EndUnixTimestamp", currentUnixTimestamp}}
	findOptions := options.Find()
	cur, err := collectionOfStakes.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	// ====2=========================================================

	// ====3=========================================================
	// Добавляем информацию к "старым" сигналам
	log.Println("[UpdateEndingStakes] Начинаем итерацию по каждому старому сигналу")

	// Создаём словарь, в котором будем хранить пары типа
	// "НазваниеВалютнойПары":ЗначениеВалютнойПарыНаДанныйМомент
	var memorizedPairsAndPrices map[string]float64

	// Итерируемся по каждому сигналу, у которого EndUnixTimestamp == текущий timestamp
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

/* Эта функция принимает на вход название пары и массив сигналов. И ищет среди этого массива
 сигнал с такой же парой. Это нужно для того, чтобы:
	1. Мы находим в БД сигнал, который вот сейчас уже закончился;
	2. Мы хотим добавить этому сигналу цену при его "окончании";
	3. Мы используем эту функцию, передавая в неё название пары сигнала и массив свежих сигналов.
	4. Ф-я возвращает цену сигнала.
*/
func findPriceOfPairnameInAllSignals(requiredPairname string, allSignals []signal.Signal) float64 {
	for _, signal := range allSignals {
		if signal.Pairname == requiredPairname {
			return signal.CurrentPrice
		}
	}

	// TODO: И log.Fatalln, и panic - оба заканчивают программу. Поэтому panic не сработает, так как log.Fatalln
	// уже закончит программу. Я, вроде, написал так для перестраховки, но после добавления тестов можно будет
	// убрать этот panic (или убрать log.Fatalln и оставить panic).
	log.Fatalln("[findPriceOfPairnameInAllSignals] Не нашёл сигнал с парой = ", requiredPairname, "среди всех сигналов, делаю panic()")
	panic("panic that shouldn't ever work because of previous log.Fatalln, right?")
}

/*
	Эта ф-я принимает на вход текущую цену сигнала и предыдущую (изначальную) цену.
	И в зависимости от переданных значений возвращает 1, -1 или 0. Я мог не использовать отдельную функцию
	для такого кода, но мне показалось это хорошей идеей вынести такой важный функционал в отдельную ф-ю.
*/
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

/*
	Эта ф-я делает фактическое обновление "закончившейся" ставки в БД.
	Нам нужно добавлять только два поля - финальную цену и больше ли финальная цена начальной.
	На вход получает:
	1. Ссылку на коллекцию (грубо говоря ссылку на БД);
	2. ID сигнала, который нужно обновить;
	3. Финальную цену;
	4. Больше ли конечная цена, чем начальная.
*/
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
