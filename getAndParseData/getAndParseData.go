package getAndParseData

import (
	"goGetSwitch/signal"
	"strconv"
)

func GetAndParseData(baseUrl string, timeframe string, unixTimestamp int64) []signal.Signal {
	respBody := GetData(baseUrl + timeframe)

	timeframeInInt, _ := strconv.Atoi(timeframe)
	newSignalsForTimeframe := ParseData(respBody, timeframeInInt, unixTimestamp)

	return newSignalsForTimeframe
}
