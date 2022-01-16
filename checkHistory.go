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

func topCombinationsWithinBorders(collection *mongo.Collection, lowerBorder int, upperBorder int, channelForSendingTopCombinationStats chan []stats.Stats) {
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
			// 	а не меньше 43) не входят в финальный список
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
	go topCombinationsWithinBorders(collection, 1, 20, channelForGettingTopStats)

	fmt.Println("[main] Перед получением данных из канала")
	//topStats := <-channelForGettingTopStats
	topStats := <-channelForGettingTopStats
	fmt.Println("[main] После получения данных из канала")
	fmt.Println("len(topStats) = ", len(topStats))

	elapsed := time.Since(start)
	log.Printf("Program took %s", elapsed)
}
