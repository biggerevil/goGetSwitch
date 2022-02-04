package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"goGetSwitch/dbFunctions"
	"goGetSwitch/stats"
	"log"
	"math"
	"sort"
	"time"
)
import "goGetSwitch/producerCode"

/*
	Это второй main-file, в котором будет код для "проверки ставок по истории"/"сбора статистики"
*/

/*
	TODO: как сделать проверку по истории многопоточной:
	(в итоге мне кажется, что лучше делать ДРУГОЙ вариант, см. TODO ниже)
	1. Сделать переменную здесь, в которой будет храниться startCombinationNumber
	2. Запустить отдельную ф-ю для горутин, которая будет:
		2.1. Пока while не выше maxCombinationNumber (cтоит сделать отдельную функцию для определения этого. И мб
			 вообще убрать проверку из начала той функции, где она сейчас есть. Чтобы не тратить каждый раз так ресурсы)
		2.2. Поднять текущий startCombinationNumber на N (5 или 50, или какое-либо ещё число, сам решу) (и поскольку
			 поднятие происходит перед запуском checkHistory, наверное, стоит делать его -5 или -50 изначально...
			 Но как решать проблемы с совместным доступом)
		2.3. Запустить checkHistory с текущим startCombinationNumber
		2.4. Тут прервался и расписываю следующий вариант

	TODO: другой вариант, как сделать проверку по истории многопоточной:
	1. Поделить все ставки на кол-во горутин
	2. Запустить каждую горутину с необходимым оффсетом (то есть запускать отдельную функцию)
	3. В каждой горутине отдельно принимать массив с топовыми результатами
	4. В конце каждой горутины возвращать массив с топовыми результатами В КАНАЛ
	5. В main-функции ждать получения всех массивов из каналов.
	6. В конце main-функции отсортировывать лучшие результаты
	7. Выводить лучшие результаты (или записывать в mongodb, или ещё что делать, не суть. На данный момент просто
	   выводить)
*/

func topCombinationsWithinBorders(collection *mongo.Collection, lowerBorder int64, upperBorder int64, channelForSendingTopCombinationStats chan []stats.Stats) {
	fmt.Println("[topCombinationsWithinBorders] Начинаю со значениями: lowerBorder = ", lowerBorder,
		" и upperBorder = ", upperBorder)

	// Генерируем комбинации
	severalCombinations, err := producerCode.GeneratePowersetWithinBorders(lowerBorder, upperBorder)
	if err != nil {
		log.Fatalln("err = ", err)
	}

	var topCombinationsStatsWithinPassedBorders []stats.Stats

	// Просто выводим сгенерированные комбинации
	for _, combination := range severalCombinations {
		fmt.Println("combination.Conditions() = ", combination.Conditions)
		stats := dbFunctions.GetCombinationStats(combination, collection)
		fmt.Println("stats: ", stats)
		if stats.PercentOfStakesWhereEndPriceMoreThanInitial > 57 || stats.PercentOfStakesWhereEndPriceMoreThanInitial < 43 {
			topCombinationsStatsWithinPassedBorders = append(topCombinationsStatsWithinPassedBorders, stats)

			if stats.PercentOfStakesWhereEndPriceMoreThanInitial > 57 {
				log.Fatalln("stats с винрейтом больше 57. stats = ", stats)
				panic("паника-паника")
			}

			// TODO: Из-за такой сортировки все значения с топовым "ПОЛОЖИТЕЛЬНЫМ" винрейтом (то есть когда больше 57%,
			// 	а не меньше 43) не входят в финальный список. Или я это уже исправил?
			sort.SliceStable(topCombinationsStatsWithinPassedBorders, func(i, j int) bool {
				firstValue := topCombinationsStatsWithinPassedBorders[i].PercentOfStakesWhereEndPriceMoreThanInitial
				secondValue := topCombinationsStatsWithinPassedBorders[j].PercentOfStakesWhereEndPriceMoreThanInitial
				// Просто значение, чтобы в выражении в 2 местах стояла переменная.
				fifty := 50.0
				/*
					Смысл выражения ниже в том, чтобы можно было понять, что, например, 10% винрейта лучше, чем 60%
					(так как если переворачивать ставки с 10% винрейтом, то мы получаем винрейт 90%)
				*/
				return math.Abs(fifty-firstValue) > math.Abs(fifty-secondValue)
			})

			if len(topCombinationsStatsWithinPassedBorders) > 20 {
				topCombinationsStatsWithinPassedBorders = topCombinationsStatsWithinPassedBorders[:20]
			}
		}
	}

	fmt.Println("topCombinationsStatsWithinPassedBorders = ", topCombinationsStatsWithinPassedBorders)
	for index, stats := range topCombinationsStatsWithinPassedBorders {
		fmt.Println("#", index, " stats: ", stats)
	}

	fmt.Println("[topCombinationsWithinBorders] Перед отправкой в канал")
	channelForSendingTopCombinationStats <- topCombinationsStatsWithinPassedBorders
	fmt.Println("[topCombinationsWithinBorders] После отправки в канал")
}

func main() {
	start := time.Now()

	collection := dbFunctions.ConnectToDB()
	channelForGettingTopStats := make(chan []stats.Stats)

	var numberOfGoroutines int64
	numberOfGoroutines = 2
	//maxUpperBorder := producerCode.GetMaxUpperBorder()
	var startBorder int64
	//startBorder = 0
	startBorder = 1000000
	var maxUpperBorder int64
	//maxUpperBorder = 20
	maxUpperBorder = 1000020
	step := (maxUpperBorder - startBorder) / numberOfGoroutines
	fmt.Println("step = ", step)

	var topCombinationsStats []stats.Stats

	//for i := 0; i < numberOfGoroutines; i++ {
	//	go topCombinationsWithinBorders(collection, 1, 20, channelForGettingTopStats)
	//}

	var iInInt64 int64
	// TODO: При таких значения startBorder, maxUpperBorder, numberOfGoroutines и step, как у меня сейчас, условие
	//	цикла выполняется СРАЗУ, то есть цикл не делает НИ ОДНОЙ итерации. В этом баг
	for iInInt64 = startBorder; iInInt64 <= numberOfGoroutines*step; iInInt64 += step {
		fmt.Println("iInInt64 = ", iInInt64)
		go topCombinationsWithinBorders(collection, iInInt64, iInInt64+step, channelForGettingTopStats)
	}

	for iInInt64 = 0; iInInt64 < numberOfGoroutines; iInInt64++ {
		fmt.Println("[main] Перед получением данных из канала")
		topStats := <-channelForGettingTopStats
		fmt.Println("[main] После получения данных из канала")

		// #1. Сначала оставляем в topStats только лучшие 20 значений
		// TODO: Дублирование кода
		sort.SliceStable(topStats, func(i, j int) bool {
			firstValue := topStats[i].PercentOfStakesWhereEndPriceMoreThanInitial
			secondValue := topStats[j].PercentOfStakesWhereEndPriceMoreThanInitial
			// Просто значение, чтобы в выражении в 2 местах стояла переменная.
			fifty := 50.0
			/*
				Смысл выражения ниже в том, чтобы можно было понять, что, например, 10% винрейта лучше, чем 60%
				(так как если переворачивать ставки с 10% винрейтом, то мы получаем винрейт 90%)
			*/
			return math.Abs(fifty-firstValue) > math.Abs(fifty-secondValue)
		})

		if len(topStats) > 20 {
			topStats = topStats[:20]
		}

		// #2. Затем добавляем все значения из topStats в topCombinationsStats
		for _, stat := range topStats {
			topCombinationsStats = append(topCombinationsStats, stat)
		}
		fmt.Println("len(topCombinationsStats) После добавления новых топовых комбинаций = ", len(topCombinationsStats))

		// #3. Затем сортируем topCombinationsStats и оставляем лучшие 20 значений
		// TODO: дублирование кода
		sort.SliceStable(topCombinationsStats, func(i, j int) bool {
			firstValue := topCombinationsStats[i].PercentOfStakesWhereEndPriceMoreThanInitial
			secondValue := topCombinationsStats[j].PercentOfStakesWhereEndPriceMoreThanInitial
			// Просто значение, чтобы в выражении в 2 местах стояла переменная.
			fifty := 50.0
			/*
				Смысл выражения ниже в том, чтобы можно было понять, что, например, 10% винрейта лучше, чем 60%
				(так как если переворачивать ставки с 10% винрейтом, то мы получаем винрейт 90%)
			*/
			return math.Abs(fifty-firstValue) > math.Abs(fifty-secondValue)
		})

		if len(topCombinationsStats) > 20 {
			topCombinationsStats = topCombinationsStats[:20]
		}
	}

	fmt.Println("[main] topCombinationsStats:")
	for index, stats := range topCombinationsStats {
		fmt.Println("#", index, " stats: ", stats)
	}

	elapsed := time.Since(start)
	log.Printf("Program took %s", elapsed)
}
