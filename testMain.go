package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"goGetSwitch/dbFunctions"
	"goGetSwitch/producerCode"
	"goGetSwitch/stats"
	"log"
	"os"
	"sort"
	"time"
)

/*
	Этот main-файл я добавил для быстрой проверки некоторых функций, по типу генерации
	powerset'а. То есть, чтобы запустить и посмотреть/показать рез-т работы каких-либо функций.
*/

// Для импорта переменной conditions из файла powerset.go надо всего лишь сделать первую букву заглавной, но
// я боюсь, что сломается что-либо в других файлах, так как у меня много что названо "conditions", а тестить я сейчас
// не хочу, так как не очень много времени и желания. Поэтому просто копирую эту переменную в testMain.go. Потом буду
// делать импорт из powerset.go или как-либо иначе сделаю.
var conditions = []producerCode.Condition{
	{"Pairname", "AUD/USD"},
	{"Pairname", "EUR/CHF"},
	{"Pairname", "EUR/JPY"},
	{"Pairname", "EUR/USD"},
	{"Pairname", "GBP/USD"},
	{"Pairname", "USD/CAD"},
	{"Pairname", "USD/JPY"},

	{"Timeframe", "300"},
	{"Timeframe", "900"},
	{"Timeframe", "1800"},
	{"Timeframe", "3600"},
	{"Timeframe", "7200"},
	{"Timeframe", "18000"},
	{"Timeframe", "86400"},

	{"MaBuy", "1"},
	{"MaBuy", "2"},
	{"MaBuy", "3"},
	{"MaBuy", "4"},
	{"MaBuy", "5"},
	{"MaBuy", "6"},
	{"MaBuy", "7"},
	{"MaBuy", "8"},
	{"MaBuy", "9"},
	{"MaBuy", "10"},
	{"MaBuy", "11"},
	{"MaBuy", "12"},

	{"MaSell", "1"},
	{"MaSell", "2"},
	{"MaSell", "3"},
	{"MaSell", "4"},
	{"MaSell", "5"},
	{"MaSell", "6"},
	{"MaSell", "7"},
	{"MaSell", "8"},
	{"MaSell", "9"},
	{"MaSell", "10"},
	{"MaSell", "11"},
	{"MaSell", "12"},

	{"TiBuy", "1"},
	{"TiBuy", "2"},
	{"TiBuy", "3"},
	{"TiBuy", "4"},
	{"TiBuy", "5"},
	{"TiBuy", "6"},
	{"TiBuy", "7"},
	{"TiBuy", "8"},
	{"TiBuy", "9"},
	{"TiBuy", "10"},
	{"TiBuy", "11"},

	{"TiSell", "1"},
	{"TiSell", "2"},
	{"TiSell", "3"},
	{"TiSell", "4"},
	{"TiSell", "5"},
	{"TiSell", "6"},
	{"TiSell", "7"},
	{"TiSell", "8"},
	{"TiSell", "9"},
	{"TiSell", "10"},
	{"TiSell", "11"},

	// TODO: Добавить ROUNDTIME = 0 и ROUNDTIME = 1. Пока не знаю, как именно, но ВОЗМОЖНО
	//	как-либо просто делить unixTimestamp с остатком или что-то такое, не знаю. Хотелось бы, конечно, чтобы это
	//	работало как фильтр, на уровне find, НО можно ещё:
	//	просто проитерироваться по всем существующим ставкам (да, это будет долго, но это как бы единоразовый процесс)
	//	и каждой ставке (то есть сигналу) проставить его roundtime (0, 1, и или -1, или 2/3/4). И добавить функицонал,
	//	что при добавлении нового сигнала автоматически будет добавляться его roundtime, чтобы в будущем не приходилось
	//	делать вот так подолгу.
}

/*
	Эта функция служит для "доставания" из массива всех возможных условий (Condition) всех возможных
	значений.
	Это на данный момент (7 февраля 2022) нужно для последующей итерации по этим значениям и составления комбинаций.
*/
func getColumnNameFromConditions(requiredColumnName string, conditionsToLookFrom []producerCode.Condition) []producerCode.Condition {
	var conditionsWithRequiredColumnName []producerCode.Condition
	for _, condition := range conditionsToLookFrom {
		if condition.ColumnName == requiredColumnName {
			conditionsWithRequiredColumnName = append(conditionsWithRequiredColumnName, condition)
		}
	}
	return conditionsWithRequiredColumnName
}

/*
	Эта ф-я служит для создания комбинации из переданных условий (incomingConditions) и последующего
	вызова ф-и для подсчёта статистики по комбинации (вызова ф-и GetCombinationStats)
*/
func getStatsFor(collection *mongo.Collection, incomingConditions ...producerCode.Condition) stats.Stats {
	var combination producerCode.Combination
	for _, condition := range incomingConditions {
		fmt.Println("Adding condition = ", condition)
		combination.Conditions = append(combination.Conditions, producerCode.Condition{ColumnName: condition.ColumnName, Value: condition.Value})
	}
	statsOfCombination := dbFunctions.GetCombinationStats(combination, collection)
	fmt.Println(stats.StatsAsPrettyString(statsOfCombination))

	return statsOfCombination
}

func main() {
	// Замеряем время работы программы.
	start := time.Now()

	// 277042167809
	// 554084335617
	//powerset, _ := producerCode.GeneratePowersetWithinBorders(277042167809, 277042167815)
	//for index, value := range powerset {
	//	fmt.Println("Combination #", index, " = ", value)
	//}

	//var combination producerCode.Combination
	//
	//combination.Conditions = append(combination.Conditions, producerCode.Condition{"Pairname", "AUD/USD"})
	//combination.Conditions = append(combination.Conditions, producerCode.Condition{"Timeframe", "300"})
	//
	//fmt.Println("combination = ", combination)
	//
	//collection := dbFunctions.ConnectToDB()
	//statsOfCombination := dbFunctions.GetCombinationStats(combination, collection)
	//
	////fmt.Println("stats = ", stats)
	//stats.PrintStats(statsOfCombination)

	/*
		По очереди достаём все возможные значения всех возможный полей и сохраняем в отдельных массивах.
		Возможно, стоит делать это как-либо иначе, а не вызывать почти один и тот же код несколько раз.
	*/
	pairnamesFromConditions := getColumnNameFromConditions("Pairname", conditions)
	fmt.Println("pairnamesFromConditions = ", pairnamesFromConditions)

	timeframesFromConditions := getColumnNameFromConditions("Timeframe", conditions)
	fmt.Println("timeframesFromConditions = ", timeframesFromConditions)

	maBuysFromConditions := getColumnNameFromConditions("MaBuy", conditions)
	fmt.Println("maBuysFromConditions = ", maBuysFromConditions)

	maSellsFromConditions := getColumnNameFromConditions("MaSell", conditions)
	fmt.Println("maSellsFromConditions = ", maSellsFromConditions)

	tiBuysFromConditions := getColumnNameFromConditions("TiBuy", conditions)
	fmt.Println("tiBuysFromConditions = ", tiBuysFromConditions)

	tiSellsFromConditions := getColumnNameFromConditions("TiSell", conditions)
	fmt.Println("tiSellsFromConditions = ", tiSellsFromConditions)

	fmt.Println("\n\n\n")

	// Создаём массив, в котором будем хранить ВСЮ статистику по нашим комбинациям.
	var resultStats []stats.Stats
	// Подключаемся к БД и сохраняем объект подключения в переменной collection (некоторые ф-и этого
	// проекта для работы с БД принимают на вход объект подключения, а точнее коллекции).
	collection := dbFunctions.ConnectToDB()

	/*
		В этих циклах мы по очереди берём:
		1. Сначала каждое название пары;
		2. Затем каждое значение таймфрейма;
		3. Затем каждое значение MaBuy;
		4. И затем каждое значение TiBuy и TiSell. TiBuy и TiSell берутся отдельно друг от друга. То есть
		для каждой комбинации (пара; таймфрейм; MaBuy) мы проверим (пара; таймфрейм; MaBuy; TiBuy), а затем
		проверим (пара; таймфрейм; MaBuy; TiSell).
	*/
	/*
		TODO: Добавить работу с крайними значениями TiBuy и TiSell.
		 Вот как здесь - TiBuy всегда равен 1, а TiSell от 7 до 9.
		 И мы хотим проверить все комбинации. Возможно, мы хотим для КАЖДОГО TiBuy и
		 затем для КАЖДОГО (тоже для каждого или всё-таки нет?) TiSell проверять
		 все значения соответственно TiSell и TiBuy.
		 > db.stakes.find({ "TiBuy": 1, "TiSell": 8 }).count()
		 59886
		 > db.stakes.find({ "TiBuy": 1, "TiSell": 7 }).count()
		 60587
		 > db.stakes.find({ "TiBuy": 1, "TiSell": 9 }).count()
		 30753
	*/
	for _, pairnameDesiredCondition := range pairnamesFromConditions {
		for _, timeframeDesiredCondition := range timeframesFromConditions {
			for _, maBuyDesiredCondition := range maBuysFromConditions {
				// На данный момент имеем pairname, timeframe и maBuy.

				// Проверяем в двух случаях:
				// 1. С tiBuy
				for _, tiBuyDesiredCondition := range tiBuysFromConditions {
					returnedStats := getStatsFor(collection, pairnameDesiredCondition, timeframeDesiredCondition, maBuyDesiredCondition, tiBuyDesiredCondition)
					resultStats = append(resultStats, returnedStats)
				}
				// 2. С tiSell
				for _, tiSellDesiredCondition := range tiSellsFromConditions {
					returnedStats := getStatsFor(collection, pairnameDesiredCondition, timeframeDesiredCondition, maBuyDesiredCondition, tiSellDesiredCondition)
					resultStats = append(resultStats, returnedStats)
				}
			}

			for _, maSellDesiredCondition := range maSellsFromConditions {
				// На данный момент имеем pairname, timeframe и maSell.

				// Проверяем в двух случаях:
				// 1. С tiBuy
				for _, tiBuyDesiredCondition := range tiBuysFromConditions {
					returnedStats := getStatsFor(collection, pairnameDesiredCondition, timeframeDesiredCondition, maSellDesiredCondition, tiBuyDesiredCondition)
					resultStats = append(resultStats, returnedStats)
				}
				// 2. С tiSell
				for _, tiSellDesiredCondition := range tiSellsFromConditions {
					returnedStats := getStatsFor(collection, pairnameDesiredCondition, timeframeDesiredCondition, maSellDesiredCondition, tiSellDesiredCondition)
					resultStats = append(resultStats, returnedStats)
				}
			}
		}
	}

	fmt.Println("resultStats = ", resultStats)
	fmt.Println("len(resultStats) = ", len(resultStats))

	// Сортируем массив со всей статистикой по кол-ву ставок комбинации.
	sort.Slice(resultStats, func(i, j int) bool {
		return resultStats[i].StakesAtAll > resultStats[j].StakesAtAll
	})

	// Создаём файл, чтобы записать в него данные.
	f, err := os.Create("resultingStats.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Указываем, что при завершении main() мы хотим закрыть этот файл.
	defer f.Close()

	// Итерируемся по всей статистике и записываем в файл.
	for _, resultingStats := range resultStats {
		// Получаем string со статистикой.
		stringToWriteDown := stats.StatsAsPrettyString(resultingStats)
		// Записываем в файл.
		_, err2 := f.WriteString(stringToWriteDown)
		if err2 != nil {
			log.Fatal(err2)
		}
	}

	// Заканчиваем замер времени работы программы и выводим эту информацию.
	elapsed := time.Since(start)
	fmt.Printf("Done in %s\n", elapsed)

}
