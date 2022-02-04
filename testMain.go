package main

import (
	"fmt"
	"goGetSwitch/producerCode"
)

/*
	Этот main-файл я добавил для быстрой проверки некоторых функций, по типу генерации
	powerset'а. То есть, чтобы запустить и посмотреть/показать рез-т работы каких-либо функций.
*/

func main() {
	powerset, _ := producerCode.GeneratePowersetWithinBorders(554084335616, 554084335618)
	for index, value := range powerset {
		fmt.Println("Combination #", index, " = ", value)
	}
	//fmt.Println("powerset = ", powerset)
}
