package producerCode

import (
	"errors"
	"fmt"
	"math"
)

/*
	ВАЖНО!!!

	Этот файл служит для работы с powerset. На данный момент (7 февраля 2022) я НЕ планирую
	использовать powerset, так как слишком огромное число, а для генерации интересующих меня
	комбинаций можно использовать другой подход, чуть более понятный.
*/

// Комбинации состоят из условий. Эта структура представляет собой структуру условия.
type Condition struct {
	ColumnName string `json:"columnName"`
	Value      string `json:"value"`
}

// Комбинации состоят из условий. Эта структура представляет собой массив условий.
// TODO:плохое объяснение, переписать как-либо или убрать его, и сказать, чтобы меня спрашивали.
type Combination struct {
	Conditions []Condition `json:"conditions"`
}

/*
	Это массив со всеми возможными условиями, из которых будут составляться комбинации.
	Здесь у нас все возможные пары, все возможные таймфреймы, и все возможные значения индикаторов.

	TODO: может быть, в качестве проверки ставок в какое-либо применять start(и/или end)UnixTimestamp в диапазоне от
	 4 до 7 (условно, это для утра. Но числа могут быть и другие), а выяснять это при помощи
	 деления start(и/или end)UnixTimestamp с остатком (то есть оператор %) на 24?
*/

var conditions = []Condition{
	Condition{"Pairname", "AUD/USD"},
	Condition{"Pairname", "EUR/CHF"},
	Condition{"Pairname", "EUR/JPY"},
	Condition{"Pairname", "EUR/USD"},
	Condition{"Pairname", "GBP/USD"},
	Condition{"Pairname", "USD/CAD"},
	Condition{"Pairname", "USD/JPY"},

	Condition{"Timeframe", "300"},
	Condition{"Timeframe", "900"},
	Condition{"Timeframe", "1800"},
	Condition{"Timeframe", "3600"},
	Condition{"Timeframe", "7200"},
	Condition{"Timeframe", "18000"},
	Condition{"Timeframe", "86400"},

	Condition{"MaBuy", "1"},
	Condition{"MaBuy", "2"},
	Condition{"MaBuy", "3"},
	Condition{"MaBuy", "4"},
	Condition{"MaBuy", "5"},
	Condition{"MaBuy", "6"},
	Condition{"MaBuy", "7"},
	Condition{"MaBuy", "8"},
	Condition{"MaBuy", "9"},
	Condition{"MaBuy", "10"},
	Condition{"MaBuy", "11"},
	Condition{"MaBuy", "12"},

	Condition{"MaSell", "1"},
	Condition{"MaSell", "2"},
	Condition{"MaSell", "3"},
	Condition{"MaSell", "4"},
	Condition{"MaSell", "5"},
	Condition{"MaSell", "6"},
	Condition{"MaSell", "7"},
	Condition{"MaSell", "8"},
	Condition{"MaSell", "9"},
	Condition{"MaSell", "10"},
	Condition{"MaSell", "11"},
	Condition{"MaSell", "12"},

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

	Condition{"TiSell", "1"},
	Condition{"TiSell", "2"},
	Condition{"TiSell", "3"},
	Condition{"TiSell", "4"},
	Condition{"TiSell", "5"},
	Condition{"TiSell", "6"},
	Condition{"TiSell", "7"},
	Condition{"TiSell", "8"},
	Condition{"TiSell", "9"},
	Condition{"TiSell", "10"},
	Condition{"TiSell", "11"},

	// TODO: Добавить ROUNDTIME = 0 и ROUNDTIME = 1. Пока не знаю, как именно, но ВОЗМОЖНО просто
	//  как-либо просто делить unixTimestamp с остатком или что-то такое, не знаю. Хотелось бы, конечно, чтобы это
	//  работало как фильтр, на уровне find, НО можно ещё:
	//  просто проитерироваться по всем существующим ставкам (да, это будет долго, но это как бы единоразовый процесс)
	//  и каждой ставке (то есть сигналу) проставить его roundtime (0, 1, и или -1, или 2/3/4). И добавить функицонал,
	//  что при добавлении нового сигнала автоматически будет добавляться его roundtime, чтобы в будущем не приходилось
	//  делать вот так подолгу.
}

/*
	Эта функция генерирует набор комбинаций, принимая на вход два десятичных числа.
	Каждая комбинация идёт под номером от 0 до *максимальное_число_комбинаций*.
	Всего ВОЗМОЖНЫХ комбинаций ОЧЕНЬ много, поэтому пытаться их обработать за раз бесполезно.

	Поэтому я написал эту функцию, чтобы можно было:
	1. Взял N комбинаций (например, N=50);
	2. Обработал N комбинаций;
	3. Взял следующие N комбинаций.
	То есть в плане кода это выглядело бы ПРИМЕРНО так:
	1. combinations := GeneratePowersetWithinBorders(0, 50);
	2. getStatsForCombinations(combinations);
	3. combinations := GeneratePowersetWithinBorders(50, 100);
*/
func GeneratePowersetWithinBorders(lowerBorder int64, upperBorder int64) ([]Combination, error) {
	/*
		Проверка, что переданный upperBorder не превышает допустимое значение
		Заметь, что здесь сначала считается степень, а уже ПОСЛЕ делается -1

		Максимально возможная комбинация считается по формуле:
		2 ** (количество_условий)
		(условие = Condition)

		Эта формула из информатики, сейчас, пока пишу этот комментарий, не могу вспомнить, почему именно
		эта формула, но попробуйте представить себе, что у вас комбинации состоят из
		двух условий - например, {Pairname: USD/JPY} и {Pairname: EUR/CHF}. Сколько у вас
		получится возможных комбинаций при наличии двух условий? (Ответ: 4, то есть 2 ** количество_условий.)
		И попробуйте проверить при наличии трёх условий - например, добавить {Timeframe: 300}. (Ответ: 8.)
	*/
	maxUpperBorder := int64(math.Pow(2, float64(len(conditions))) - 1)
	fmt.Println("float64(len(conditions)) - 1 = ", float64(len(conditions))-1)
	fmt.Println("math.Pow(2, float64(len(conditions))) - 1 = ", math.Pow(2, float64(len(conditions)))-1)
	fmt.Println("maxUpperBorder = ", maxUpperBorder)

	// Проверяем, больше ли переданная верхняя граница чем максимально возможная граница.
	if upperBorder > maxUpperBorder {
		// Возвращаем ошибку, если переданный upperBorder не превышает допустимое значение
		return nil, errors.New("error - passed upperBorder > than maxUpperBorder")
	}

	// Создаём результирующий массив комбинаций, который вернёт функция.
	// В процессе работы функция будет добавлять комбинации в этот массив
	var severalCombinations []Combination

	// Идея для алгоритма работы взята отсюда - https://stackoverflow.com/a/2779467/8604912
	// lowerBorder и upperBorder существуют для того, чтобы эта функция могла отработать НЕ ТОЛЬКО с начала до конца,
	// но и в несколько этапов (грубо говоря сначала числа с 0 до 50, потом с 50 до 100,
	// потом, например, с 100 до макс. число, но вряд ли по макс. число, так как оно слишком большое.)
	for i := lowerBorder; i <= upperBorder; i++ {
		// s - это двоичное представление числа i. Алгоритм основан на использовании двоичного представления.
		// (см. stackoverflow-ссылку, откуда взята идея для работы алгоритма)
		s := fmt.Sprintf("%b", i)

		// Создаём комбинацию для этой итерации. В эту комбинацию будем добавлять необходимые Condition
		var combination Combination

		// counter нужен для того, чтобы определить, брать ли condition. Я думаю, этот алгоритм можно как-либо улучшить
		// TODO: counter после коммита 4 февраля (2022 г.) по идее больше не нужен.
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
				combination.Conditions = append(combination.Conditions, conditions[bit])
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

// Эта функция просто возвращает максимально возможную границу. Сейчас я это не использую.
func GetMaxUpperBorder() int64 {
	return int64(int(math.Pow(2, float64(len(conditions))) - 1))
}
