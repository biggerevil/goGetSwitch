package main

import (
	"fmt"
	"goGetSwitch/dbFunctions"
	"log"
)
import "goGetSwitch/producerCode"

/*
	Это второй main-file, в котором будет код для "проверки ставок по истории"/"сбора статистики"
*/

func main() {
	fmt.Println("Hello there")

	// Генерируем комбинации
	severalCombinations, err := producerCode.GeneratePowersetWithinBorders(0, 5)
	if err != nil {
		log.Fatalln("err = ", err)
	}

	collection := dbFunctions.ConnectToDB()

	// Просто выводим сгенерированные комбинации
	for _, combination := range severalCombinations {
		fmt.Println("combination.Conditions() = ", combination.Conditions)
		stats := dbFunctions.GetCombinationStats(combination, collection)
		fmt.Println("stats: ", stats)
	}

}
