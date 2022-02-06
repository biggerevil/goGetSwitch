package stats

import (
	"fmt"
	"goGetSwitch/producerCode"
	"strconv"
)

type Stats struct {
	Combination                                 producerCode.Combination
	StakesAtAll                                 int64
	StakesWhereEndPriceMoreThanInitialCount     int64
	PercentOfStakesWhereEndPriceMoreThanInitial float64
}

// ByAge implements sort.Interface based on the Age field.
//type ByPercent []Stats
//
//func (statsArray ByPercent) Len() int           { return len(statsArray) }
//func (statsArray ByPercent) Less(i, j int) bool { return statsArray[i].PercentOfStakesWhereEndPriceMoreThanInitial < statsArray[j].PercentOfStakesWhereEndPriceMoreThanInitial }
//func (statsArray ByPercent) Swap(i, j int)      { statsArray[i], statsArray[j] = statsArray[j], statsArray[i] }

func ConditionsAsString(incomingStats Stats) string {
	stringWithConditions := ""
	for _, condition := range incomingStats.Combination.Conditions {
		stringWithConditions += condition.ColumnName + ": " + condition.Value + ", "
	}
	return stringWithConditions
}

func StatsAsPrettyString(incomingStats Stats) string {
	stringToReturn := "Combination: " + ConditionsAsString(incomingStats) + "\nStakes at all: " + strconv.FormatInt(incomingStats.StakesAtAll, 10) +
		"\nStakes where end price more than initial: " + strconv.FormatInt(incomingStats.StakesWhereEndPriceMoreThanInitialCount, 10) +
		"\nPercent of stakes where end price more than initial: " +
		fmt.Sprintf("%f", incomingStats.PercentOfStakesWhereEndPriceMoreThanInitial)

	return stringToReturn
}
