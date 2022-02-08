package stats

import (
	"fmt"
	"goGetSwitch/producerCode"
	"strconv"
)

/*
	Эта структура служит для хранения данных статистики по комбинации.
	При помощи этой структуры мы можем сравнивать различные комбинации друг с другом на
	уровне кода (или просто записывать комбинации в Google Таблицы или куда-либо ещё, и там уже анализировать).
*/
type Stats struct {
	Combination producerCode.Combination
	StakesAtAll int64
	// Кол-во ставок и винрейт относительно всех абсолютно всех ставок
	AllRelativeStakesWhereEndPriceMoreThanInitialCount     int64
	AllRelativePercentOfStakesWhereEndPriceMoreThanInitial float64
	// Кол-во ставок и винрейт относительно всех ставок ЗА ИСКЛЮЧЕНИЕМ ставок, где
	// endPriceMoreThanInitial = 0
	WithoutZeroRelativeStakesWhereEndPriceMoreThanInitialCount     int64
	WithoutZeroRelativePercentOfStakesWhereEndPriceMoreThanInitial float64
}

// ByAge implements sort.Interface based on the Age field.
//type ByPercent []Stats
//
//func (statsArray ByPercent) Len() int           { return len(statsArray) }
//func (statsArray ByPercent) Less(i, j int) bool { return statsArray[i].AllRelativePercentOfStakesWhereEndPriceMoreThanInitial < statsArray[j].AllRelativePercentOfStakesWhereEndPriceMoreThanInitial }
//func (statsArray ByPercent) Swap(i, j int)      { statsArray[i], statsArray[j] = statsArray[j], statsArray[i] }

/*
	Выводим условия комбинации в виде строки. На данный момент эта функция нужна только для создания
	строки из объекта Stats.
*/
func ConditionsAsString(incomingStats Stats) string {
	stringWithConditions := ""
	for _, condition := range incomingStats.Combination.Conditions {
		stringWithConditions += condition.ColumnName + ": " + condition.Value + ", "
	}
	return stringWithConditions
}

/*
	Создаём строку из объекта Stats.
	Ранее тут была функция, которая просто выводила в консоль значения полей Stats, но я решил, что
	создавать строку, так как строку я смогу как вывести при помощи того же fmt.Println, так и при
	необходимости записать, например, в текстовый файл.
*/
func StatsAsPrettyString(incomingStats Stats) string {
	stringToReturn := "Combination: " + ConditionsAsString(incomingStats) + "\nStakes at all: " + strconv.FormatInt(incomingStats.StakesAtAll, 10) +
		"\nStakes where end price more than initial: " + strconv.FormatInt(incomingStats.AllRelativeStakesWhereEndPriceMoreThanInitialCount, 10) +
		"\nPercent of stakes where end price more than initial: " +
		fmt.Sprintf("%f", incomingStats.AllRelativePercentOfStakesWhereEndPriceMoreThanInitial)

	return stringToReturn
}
