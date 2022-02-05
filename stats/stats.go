package stats

import (
	"fmt"
	"goGetSwitch/producerCode"
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

func PrintStats(incomingStats Stats) {
	fmt.Println("Combination: ", incomingStats.Combination, "\nStakes at all: ", incomingStats.StakesAtAll,
		"\nStakes where end price more than initial: ", incomingStats.StakesWhereEndPriceMoreThanInitialCount,
		"\nPercent of stakes where end price more than initial: ",
		incomingStats.PercentOfStakesWhereEndPriceMoreThanInitial)
}
