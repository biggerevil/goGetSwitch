package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goGetSwitch/dbFunctions"
	"goGetSwitch/producerCode"
	"goGetSwitch/signal"
	"log"
	"math"
	"strconv"
	"time"
)

func addStartDateStringToStake(stakesCollection *mongo.Collection, idOfsignaltoupdate primitive.ObjectID, formattedDateString string) {
	result, err := stakesCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": idOfsignaltoupdate},
		bson.D{
			{"$set", bson.D{{"StartDateString", formattedDateString}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
}

const dateAndTimeFormat = "2006-01-02T15:04:05"
const dateOnlyFormat = "2006-01-02"

func formattedDateFromUnixTimestamp(incomingUnixTimestamp string, dateFormat string) string {
	fmt.Println("incomingUnixTimestamp =", incomingUnixTimestamp)

	i, err := strconv.ParseInt(incomingUnixTimestamp, 10, 64)
	if err != nil {
		panic(err)
	}

	tm := time.Unix(i, 0)
	// Превращаем строку в UTC
	tm = tm.UTC()
	formattedDate := tm.Format(dateFormat)
	fmt.Println("formattedDate = ", formattedDate)

	return formattedDate
}

func addStartDateStringToStakesWithoutIt() {
	start := time.Now()

	// Для тестирования
	//unixTimestamp := "1641020400"
	//fmt.Println("unixTimestamp = ", unixTimestamp)
	//formattedDateFromUnixTimestamp(unixTimestamp)

	collection := dbFunctions.ConnectToDB()

	filter := bson.D{{"StartDateString", nil}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var stake signal.Signal
		if err = cursor.Decode(&stake); err != nil {
			log.Fatal(err)
		}

		fmt.Println("stake = ", signal.SignalDataInOneStringWithComments(stake))
		//fmt.Println("stake.StartUnixTimestamp = ", stake.StartUnixTimestamp)
		startUnixTimestampInString := strconv.FormatInt(stake.StartUnixTimestamp, 10)
		//fmt.Println("startUnixTimestampInString = ", startUnixTimestampInString)

		formattedStartDateString := formattedDateFromUnixTimestamp(startUnixTimestampInString, dateAndTimeFormat)
		addStartDateStringToStake(collection, stake.ID, formattedStartDateString)
	}

	elapsed := time.Since(start)
	log.Printf("Done in %s", elapsed)
}

/*
	Код для проверки, что unixTimestamp корректно превращается моей функцией в дату:

	currentUTCTimestamp := time.Now().UTC().Unix()
	fmt.Println("currentTimestamp = ", currentUTCTimestamp)
	currentUTCTimestampAsString := strconv.FormatInt(currentUTCTimestamp, 10)
	fmt.Println("currentUTCTimestampAsString = ", currentUTCTimestampAsString)
	currentUTCDate := formattedDateFromUnixTimestamp(currentUTCTimestampAsString, dateOnlyFormat)
	fmt.Println("currentUTCDate = ", currentUTCDate)

*/

func composeCombinationFromCondition(incomingConditions ...producerCode.Condition) producerCode.Combination {
	var combination producerCode.Combination
	for _, condition := range incomingConditions {
		//fmt.Println("Adding condition = ", condition)
		combination.Conditions = append(combination.Conditions, producerCode.Condition{ColumnName: condition.ColumnName, Value: condition.Value})
	}
	return combination
}

func getEarliestAndLatestSignalsDates() (int64, int64) {
	collection := dbFunctions.ConnectToDB()

	//filter := bson.D{{"", nil}}

	// $natural: 1 => сортируем по возрастанию
	opts := options.FindOne().SetSort(bson.M{"$natural": 1})
	var firstRecord signal.Signal
	if err := collection.FindOne(context.TODO(), bson.M{}, opts).Decode(&firstRecord); err != nil {
		log.Fatal(err)
	}
	//fmt.Println(firstRecord)

	earliestStartUnixTimestampAsString := strconv.FormatInt(firstRecord.StartUnixTimestamp, 10)
	earliestDate := formattedDateFromUnixTimestamp(earliestStartUnixTimestampAsString, dateOnlyFormat)
	fmt.Println("earliestDate = ", earliestDate)

	// $natural: -1 => сортируем по убыванию
	opts = options.FindOne().SetSort(bson.M{"$natural": -1})
	var lastRecord signal.Signal
	if err := collection.FindOne(context.TODO(), bson.M{}, opts).Decode(&lastRecord); err != nil {
		log.Fatal(err)
	}

	latestStartUnixTimestampAsString := strconv.FormatInt(lastRecord.StartUnixTimestamp, 10)
	latestDate := formattedDateFromUnixTimestamp(latestStartUnixTimestampAsString, dateOnlyFormat)
	fmt.Println("latestDate = ", latestDate)

	fmt.Println("firstRecord.StartUnixTimestamp = ", firstRecord.StartUnixTimestamp)
	fmt.Println("lastRecord.StartUnixTimestamp = ", lastRecord.StartUnixTimestamp)

	return firstRecord.StartUnixTimestamp, lastRecord.StartUnixTimestamp
}

func daysBetweenEarliestAndLatestSignalsAtAll() float64 {
	earliestUnixTimestamp, latestUnixTimestamp := getEarliestAndLatestSignalsDates()

	earliest := time.Unix(earliestUnixTimestamp, 0)
	fmt.Println("earliest = ", earliest)

	latest := time.Unix(latestUnixTimestamp, 0)
	fmt.Println("latest = ", latest)

	diff := latest.Sub(earliest)
	fmt.Println("diff.Hours() = ", diff.Hours())

	hours := diff.Hours()
	days := math.Ceil(hours / 24)
	fmt.Println("days = ", days)

	return days
}

func getResultKoeff(stakesAtDays map[string]float64, totalAmountOfDays float64) float64 {
	/*
		stakesAtDays - это массив по типу:
			stakesAtDays = []
			stakesAtDays = append(set, 1499)
			stakesAtDays = append(set, 1499)
			stakesAtDays = append(set, 2)

		То есть в этом stakesAtDays у нас будет:
			Два дня, где было 1499 ставок;
			Один день, где было 2 ставки;

	*/

	// Считаем общее кол-во ставок
	var stakesAtAll float64
	for _, value := range stakesAtDays {
		stakesAtAll += value
	}

	// Определяем кол-во дней, когда у нас были ставки.
	amountOfDaysWithStakes := float64(len(stakesAtDays))
	// Считаем эталон.
	etalon := stakesAtAll / totalAmountOfDays

	/*
		pervayaSumma - это в примере Егора это 1350. Пример = сообщение 8 февраля в 22:47.
		Не знаю сейчас, как назвать переменную правильно.
	*/
	pervayaSumma := 0.0
	for _, stakesAtOneDay := range stakesAtDays {
		// Считаем отклонение
		otklonenie := stakesAtOneDay - etalon
		// Приплюсовывем отклонение к нашей pervayaSumma
		pervayaSumma = pervayaSumma + otklonenie
	}

	// Считаем кол-во дней, когда не было ставок.
	daysWithoutStakes := totalAmountOfDays - amountOfDaysWithStakes
	// Используем формулу Егора
	resultKoeff := (pervayaSumma + daysWithoutStakes*etalon) / stakesAtAll

	// Возвращаем итоговый коэффициент.
	return resultKoeff
}

func getKoeffForCombination(incomingCombination producerCode.Combination, totalAmountOfDays float64) float64 {
	collection := dbFunctions.ConnectToDB()

	filter := dbFunctions.MakeFilter(incomingCombination)

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	mapWithStakes := make(map[string]float64)

	for cursor.Next(context.TODO()) {
		var stake signal.Signal
		if err = cursor.Decode(&stake); err != nil {
			log.Fatal(err)
		}

		fmt.Println("stake = ", signal.SignalDataInOneStringWithComments(stake))
		//fmt.Println("stake.StartUnixTimestamp = ", stake.StartUnixTimestamp)
		startUnixTimestampInString := strconv.FormatInt(stake.StartUnixTimestamp, 10)
		//fmt.Println("startUnixTimestampInString = ", startUnixTimestampInString)

		formattedStartDateString := formattedDateFromUnixTimestamp(startUnixTimestampInString, dateOnlyFormat)
		fmt.Println("formattedStartDateString = ", formattedStartDateString)

		// +1 если дата уже существует. =1, если даты ещё не было.
		if _, ok := mapWithStakes[formattedStartDateString]; ok {
			mapWithStakes[formattedStartDateString] += 1
		} else {
			mapWithStakes[formattedStartDateString] = 1
		}
	}

	//fmt.Println("mapWithStakes = ", mapWithStakes)
	for date, amountOfStakes := range mapWithStakes {
		fmt.Println("[", date, "]: ", amountOfStakes)
	}

	resultKoeff := getResultKoeff(mapWithStakes, totalAmountOfDays)
	fmt.Println("resultKoeff = ", resultKoeff)

	return resultKoeff
}

func main() {
	// TODO: придумать другое название для переменной. Сейчас тут такое же, как и у функции.
	//  (эта переменная обозначает кол-во дней вообще. НЕ у конкретной комбинации, а вообще.)
	daysBetweenEarliestAndLatestSignalsAtAll := daysBetweenEarliestAndLatestSignalsAtAll()
	fmt.Println("daysBetweenEarliestAndLatestSignalsAtAll = ", daysBetweenEarliestAndLatestSignalsAtAll)

	{
		pairnameCondition := producerCode.Condition{"Pairname", "EUR/JPY"}
		timeframeCondition := producerCode.Condition{"Timeframe", "900"}
		maBuyCondition := producerCode.Condition{"MaBuy", "1"}
		tiSellCondition := producerCode.Condition{"TiSell", "4"}
		firstCombinationToTest := composeCombinationFromCondition(pairnameCondition, timeframeCondition, maBuyCondition, tiSellCondition)

		firstResultKoeff := getKoeffForCombination(firstCombinationToTest, daysBetweenEarliestAndLatestSignalsAtAll)
		fmt.Println("firstResultKoeff = ", firstResultKoeff)
	}

	{
		pairnameCondition := producerCode.Condition{"Pairname", "EUR/JPY"}
		timeframeCondition := producerCode.Condition{"Timeframe", "900"}
		maBuyCondition := producerCode.Condition{"MaBuy", "9"}
		tiSellCondition := producerCode.Condition{"TiSell", "4"}
		secondCombinationToTest := composeCombinationFromCondition(pairnameCondition, timeframeCondition, maBuyCondition, tiSellCondition)

		secondResultKoeff := getKoeffForCombination(secondCombinationToTest, daysBetweenEarliestAndLatestSignalsAtAll)
		fmt.Println("secondResultKoeff = ", secondResultKoeff)
	}
}
