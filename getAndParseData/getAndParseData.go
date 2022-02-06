package getAndParseData

import (
	"goGetSwitch/signal"
	"strconv"
)

/*
	Это начальная функция, которую мы по идее вызываем из main().
	Таким образом мы можем (грубо говоря) как угодно менять получение данных, так как для main() важно только
	что данные возвращаются массивом signal.Signal. А откуда данные берутся, как они парсятся - это
	для main() неважно.
	TODO: (Хотя мы передаём из main() baseUrl, так что main() немного в курсе, что и как происходит, но после написания
	 тестов этот baseUrl можно перенести сюда, так как в main() он больше нигде не используется на данный момент.)
*/
func GetAndParseData(baseUrl string, timeframe string, unixTimestamp int64) []signal.Signal {
	// Сначала мы получаем данные в json.
	respBody := GetData(baseUrl + timeframe)

	// Затем мы понимаем, с каким timeframe мы работаем.
	timeframeInInt, _ := strconv.Atoi(timeframe)
	// И затем мы парсим из json данные.
	newSignalsForTimeframe := ParseData(respBody, timeframeInInt, unixTimestamp)

	// И возвращаем массив сигналов для нашего timeframe.
	return newSignalsForTimeframe
}
