package main

import (
	"fmt"
	"goGetSwitch/dbFunctions"
	"goGetSwitch/producerCode"
	"goGetSwitch/stats"
)

/*
	Этот main-файл я добавил для быстрой проверки некоторых функций, по типу генерации
	powerset'а. То есть, чтобы запустить и посмотреть/показать рез-т работы каких-либо функций.
*/

func main() {
	// 277042167809
	// 554084335617
	//powerset, _ := producerCode.GeneratePowersetWithinBorders(277042167809, 277042167815)
	//for index, value := range powerset {
	//	fmt.Println("Combination #", index, " = ", value)
	//}

	var combination producerCode.Combination

	combination.Conditions = append(combination.Conditions, producerCode.Condition{"Pairname", "AUD/USD"})
	combination.Conditions = append(combination.Conditions, producerCode.Condition{"Timeframe", "300"})

	fmt.Println("combination = ", combination)

	collection := dbFunctions.ConnectToDB()
	statsOfCombination := dbFunctions.GetCombinationStats(combination, collection)

	//fmt.Println("stats = ", stats)
	stats.PrintStats(statsOfCombination)
}
