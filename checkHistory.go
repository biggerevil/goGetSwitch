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
