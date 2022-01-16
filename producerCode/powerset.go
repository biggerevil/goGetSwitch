package producerCode

import (
	"errors"
	"fmt"
	"math"
)

type Condition struct {
	ColumnName string `json:"columnName"`
	Value      string `json:"value"`
}

type Combination struct {
	Conditions []Condition `json:"conditions"`
}

// TODO: может быть, в качестве проверки ставок в какое-либо применять start(и/или end)UnixTimestamp в диапазоне от
// 	4 до 7 (условно, это для утра. Но числа могут быть и другие), а выяснять это при помощи
// 	деления start(и/или end)UnixTimestamp с остатком (то есть оператор %) на 24?
var conditions = []Condition{
	//Condition{"Pairname", "AUD/USD"},
	//Condition{"Pairname", "EUR/CHF"},
	//Condition{"Pairname", "EUR/JPY"},
	//Condition{"Pairname", "EUR/USD"},
	//Condition{"Pairname", "GBP/USD"},
	//Condition{"Pairname", "USD/CAD"},
	//Condition{"Pairname", "USD/JPY"},
	//Condition{"Timeframe", "300"},
	//Condition{"Timeframe", "900"},
	//Condition{"Timeframe", "1800"},
	//Condition{"Timeframe", "3600"},
	//Condition{"Timeframe", "7200"},
	//Condition{"Timeframe", "18000"},
	//Condition{"Timeframe", "86400"},
	Condition{"TiBuy", "1"},
	Condition{"TiBuy", "2"},
	Condition{"TiBuy", "3"},
	Condition{"TiBuy", "4"},
	Condition{"TiBuy", "5"},
	Condition{"TiBuy", "6"},
	Condition{"TiBuy", "7"},
	Condition{"TiBuy", "8"},
	Condition{"TiBuy", "9"},
	Condition{"TiBuy", "10"},
	Condition{"TiBuy", "11"},
}

func GeneratePowersetWithinBorders(lowerBorder int, upperBorder int) ([]Combination, error) {
	// Проверка, что переданный upperBorder не превышает допустимое значение
	// Заметь,что здесь сначала считается степень, а уже ПОСЛЕ делается -1
	maxUpperBorder := int(math.Pow(2, float64(len(conditions))) - 1)
	fmt.Println("float64(len(conditions)) - 1 = ", float64(len(conditions))-1)
	fmt.Println("math.Pow(2, float64(len(conditions))) - 1 = ", math.Pow(2, float64(len(conditions)))-1)
	fmt.Println("maxUpperBorder = ", maxUpperBorder)

	if upperBorder > maxUpperBorder {
		// Возвращаем ошибку, если переданный upperBorder не превышает допустимое значение
		return nil, errors.New("error - passed upperBorder > than maxUpperBorder")
	}

	// Создаём результирующий массив комбинаций, который вернёт функция.
	// В процессе работы функция будет добавлять комбинации в этот массив
	var severalCombinations []Combination

	// Идея для алгоритма работы взята отсюда - https://stackoverflow.com/a/2779467/8604912
	// lowerBorder и upperBorder существуют для того, чтобы эта функция могла отработать НЕ ТОЛЬКО с начала до конца,
	// но и в несколько этапов (грубо говоря сначала числа с 0 до 50, потом с 50 до 100, потом с 100 до макс. число)
	for i := lowerBorder; i <= upperBorder; i++ {
		// s - это двоичное представление числа i. Алгоритм основан на использовании двоичного представления.
		// (см. stackoverflow-ссылку, откуда взята идея для работы алгоритма)
		s := fmt.Sprintf("%b", i)

		// Создаём комбинацию для этой итерации. В эту комбинацию будем добавлять необходимые Condition
		var combination Combination

		// counter нужен для того, чтобы определить, брать ли condition. Я думаю, этот алгоритм можно как-либо улучшить
		counter := 0
		// Итерируемся по двоичному представлению числа (по переменной s).
		// По очереди будем брать каждую цифру из двоичного представления
		for bit := len(s) - 1; bit >= 0; bit-- {
			// Берём цифру
			value := string(s[bit])

			// Проверяем цифру из двоичного представления на значение 1.
			// Если 1, то надо взять соответствующее значение из массива conditions (см. stackoverflow-ссылку, откуда
			// взята идея для работы алгоритма)
			if value == "1" {
				combination.Conditions = append(combination.Conditions, conditions[counter])
			}
			counter++
		}

		// Добавляем в массив комбинаций полученную комбинацию и переходим на следующую итерацию цикла с lowerBorder и
		// upperBorder
		severalCombinations = append(severalCombinations, combination)
	}

	// Возвращаем сгенерированные комбинации и "пустую" ошибку
	return severalCombinations, nil
}
