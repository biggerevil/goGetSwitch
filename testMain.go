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

	// TODO: Возможно, алгоритм GeneratePowersetWi... не совсем корректно берёт значения в соответствии с
	//  полученным числом (то есть двоичная форма, вроде, правильная, но значения не те, что я ожидаю)
	/* TODO: варианты, что делать:
	    1. Попробовать передать число, соответствующее ТОЛЬКО TiBuy:1
	    2. Попробовать удалить пробел в GeneratePowersetWi... , если он там есть
	    3. Попробовать использовать такую схему:
			3.1. Указываешь condition'ы (то есть комбинацию), которую ты хочешь
			3.2. Находишь индексы этих condition'ов в conditions (переменная в powerset.go)
			3.3. Создаёшь строку из стольких нулей, сколько у нас значений в переменной conditions (то есть 60)
			3.4. Меняешь по индексам из 3.2. строке из 3.3. меняешь нули на единицы
			3.5. Обрезаешь строку по последней единице?
			3.6. Переводишь в десятичное, чтобы передать в GeneratePowersetWi... , так как она принимает десятичные,
				 так сделать или настроить двоичное? Надо ли мне будет там передавать двоичное число? Мне же нужны
				 по сути комбинации, которые я буду передавать для сбора статистики.
				 TODO: так ли мне нужен GeneratePowersetWi... , если я хочу передавать как бы типа конкретные
	 			  комбинации? То есть не все подряд, а вот по типу {Pairname: ..., Timeframe: ..., MaBuy: 1},
				  затем {Pairname: ..., Timeframe: ..., MaBuy: 1, TiBuy: 1},
				  затем {Pairname: ..., Timeframe: ..., MaBuy: 1, TiBuy: 2},
				  затем ...,
				  затем {Pairname: ..., Timeframe: ..., MaBuy: 1, TiSell: 1},
				  затем {Pairname: ..., Timeframe: ..., MaBuy: 1, TiSell: 2},
				  затем ...,
				  затем {Pairname: ..., Timeframe: ..., MaBuy: 2},
				  затем {Pairname: ..., Timeframe: ..., MaBuy: 2, TiBuy: 1},
				  затем ...
				  затем {Pairname: ..., Timeframe: ..., MaSell: 1},
				  и затем то же самое, только вместо MaBuy теперь MaSell

				 TODO: и всё вышеуказанное в цикле из пары, внутри которого цикл по timeframe
	*/
	powerset, _ := producerCode.GeneratePowersetWithinBorders(554084335616, 554084335618)
	for index, value := range powerset {
		fmt.Println("Combination #", index, " = ", value)
	}
	//fmt.Println("powerset = ", powerset)
}
