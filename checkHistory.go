package main

import (
	"fmt"
	"goGetSwitch/dbFunctions"
	"goGetSwitch/stats"
	"log"
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

func main() {
	start := time.Now()

	// Генерируем комбинации
	severalCombinations, err := producerCode.GeneratePowersetWithinBorders(0, 16382)
	if err != nil {
		log.Fatalln("err = ", err)
	}

	collection := dbFunctions.ConnectToDB()

	//allCombinationsStatsArray := mapset.NewSet()
	var allCombinationsStatsArray []stats.Stats

	// Просто выводим сгенерированные комбинации
	for _, combination := range severalCombinations {
		fmt.Println("combination.Conditions() = ", combination.Conditions)
		stats := dbFunctions.GetCombinationStats(combination, collection)
		fmt.Println("stats: ", stats)
		if stats.PercentOfStakesWhereEndPriceMoreThanInitial > 57 || stats.PercentOfStakesWhereEndPriceMoreThanInitial < 43 {
			allCombinationsStatsArray = append(allCombinationsStatsArray, stats)

			// TODO: Из-за такой сортировки все значения с топовым "ПОЛОЖИТЕЛЬНЫМ" винрейтом (то есть когда больше 57%,
			// 	а не меньше 43) не входят в финальный список
			sort.SliceStable(allCombinationsStatsArray, func(i, j int) bool {
				return allCombinationsStatsArray[i].PercentOfStakesWhereEndPriceMoreThanInitial < allCombinationsStatsArray[j].PercentOfStakesWhereEndPriceMoreThanInitial
			})

			if len(allCombinationsStatsArray) > 20 {
				allCombinationsStatsArray = allCombinationsStatsArray[:20]
			}
		}
	}

	fmt.Println("allCombinationsStatsArray = ", allCombinationsStatsArray)
	for index, stats := range allCombinationsStatsArray {
		fmt.Println("#", index, " stats: ", stats)
	}

	elapsed := time.Since(start)
	log.Printf("Program took %s", elapsed)
}
