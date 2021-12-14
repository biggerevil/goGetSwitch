package getAndParseData

import (
	"goGetSwitch/signal"
	"strconv"
)

func GetAndParseData(baseUrl string, timeframe string) []signal.Signal {
	respBody, unixTimestamp := GetData(baseUrl + timeframe)

	timeframeInInt, _ := strconv.Atoi(timeframe)
	newSignalsForTimeframe := ParseData(respBody, timeframeInInt, unixTimestamp)

	return newSignalsForTimeframe
}
