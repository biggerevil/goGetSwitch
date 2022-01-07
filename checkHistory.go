package main

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"go.mongodb.org/mongo-driver/mongo"
	"goGetSwitch/dbFunctions"
	"log"
)
import "goGetSwitch/producerCode"

/*
	Это второй main-file, в котором будет код для "проверки ставок по истории"/"сбора статистики"
*/

func getAndCountCombination(startCombinationNumber int, howMuchCombinations int, collection *mongo.Collection, channelForGettingArraysOfStats chan []map[string]interface{}) {
	// Генерируем комбинации
	severalCombinations, err := producerCode.GeneratePowersetWithinBorders(startCombinationNumber, startCombinationNumber+howMuchCombinations)
	if err != nil {
		log.Fatalln("err = ", err)
	}

	var arrayOfStats []map[string]interface{}

	// Просто выводим сгенерированные комбинации
	for _, combination := range severalCombinations {
		fmt.Println("combination.Conditions() = ", combination.Conditions)
		stats := dbFunctions.GetCombinationStats(combination, collection)
		fmt.Println("Stats for current combination:")
		for key, value := range stats {
			fmt.Println(key, " = ", value)
		}

		arrayOfStats = append(arrayOfStats, stats)
	}

	channelForGettingArraysOfStats <- arrayOfStats
}

//func addToArrayOfAllStats(arrayWithStatsFromSeveralCombinations []map[string]interface{}, allCombinationsStatsSet mapset.Set) {
//	for _, statsFromCombination := range arrayWithStatsFromSeveralCombinations {
//		allCombinationsStatsSet.Add(statsFromCombination)
//		if allCombinationsStatsSet.Cardinality() > 20 {
//			allCombinationsStatsSet.Pop()
//		}
//	}
//
//	return allCombinationsStatsSet
//}

func main() {
	fmt.Println("Hello there")

	//maxNumber := 16382
	maxNumber := 20

	collection := dbFunctions.ConnectToDB()

	numberOfGoroutines := 20
	// Возможно неточное значение. Потерять парочку значений в конце, как я думаю сейчас - не страшно
	howMuchCombinationsForRoutine := maxNumber / numberOfGoroutines
	fmt.Println("howMuchCombinationsForRoutine = ", howMuchCombinationsForRoutine)

	channelForGettingArraysOfStats := make(chan []map[string]interface{})

	startCombinationNumber := 0

	// -1 потому что делаем continue при последней итерации
	for i := 0; i < numberOfGoroutines-1; i++ {
		startCombinationNumber += howMuchCombinationsForRoutine
		if maxNumber-startCombinationNumber <= howMuchCombinationsForRoutine {
			continue
		}
		go getAndCountCombination(startCombinationNumber, howMuchCombinationsForRoutine, collection, channelForGettingArraysOfStats)
	}

	//var allCombinationsStatsArray []map[string]interface{}
	allCombinationsStatsSet := mapset.NewSet()

	// -1 потому что делаем continue при последней итерации (в предыдущем цикле)
	for i := 0; i < numberOfGoroutines; i++ {
		arrayWithStatsFromSeveralCombinations := <-channelForGettingArraysOfStats

		for _, statsFromCombination := range arrayWithStatsFromSeveralCombinations {
			allCombinationsStatsSet.Add(statsFromCombination)
			if allCombinationsStatsSet.Cardinality() > 20 {
				allCombinationsStatsSet.Pop()
			}
		}
	}

	fmt.Println("allCombinationsStatsSet = ", allCombinationsStatsSet)

	fmt.Println("End of program")
}
