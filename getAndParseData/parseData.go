package getAndParseData

import (
	"encoding/json"
	"fmt"
	"goGetSwitch/signal"
	"strconv"
)

func getIndices() []string {
	return []string{"1", "2", "3", "5", "7", "9", "10"}
}

func parsePair(data interface{}) signal.Signal {
	// При попытке распарсить сразу в int выдаёт ошибку:
	// panic: interface conversion: interface {} is float64, not int.
	// Не знаю, что с этим делать, кроме как сделать int(...) на выражение
	maBuy := int(data.(map[string]interface{})["maBuy"].(float64))
	maSell := int(data.(map[string]interface{})["maSell"].(float64))
	tiBuy := int(data.(map[string]interface{})["tiBuy"].(float64))
	tiSell := int(data.(map[string]interface{})["tiSell"].(float64))

	pairname := data.(map[string]interface{})["summaryName"].(string)
	currentPriceInString := data.(map[string]interface{})["summaryLast"].(string)
	currentPrice, _ := strconv.ParseFloat(currentPriceInString, 64)

	//var parsedSignal signal.Signal
	parsedSignal := signal.Signal{Pairname: pairname, CurrentPrice: currentPrice, MaBuy: maBuy, MaSell: maSell, TiBuy: tiBuy, TiSell: tiSell}

	return parsedSignal
}

func ParseData(respBody []byte, timeframe int, unixTimestamp int64) []signal.Signal {
	/*
		Работа с JSON
	*/
	var dat map[string]interface{}

	if err := json.Unmarshal(respBody, &dat); err != nil {
		panic(err)
	}
	fmt.Println("\n\n dat:")
	fmt.Println(dat)

	var allNewSignalsForThisTimeframe []signal.Signal

	for _, indice := range getIndices() {
		newSignal := parsePair(dat[indice])
		newSignal.Timeframe = timeframe
		newSignal.UnixTimestamp = unixTimestamp
		allNewSignalsForThisTimeframe = append(allNewSignalsForThisTimeframe, newSignal)
		fmt.Println("newSignal data with comments: ", signal.SignalDataInOneStringWithComments(newSignal))
	}

	//// Получаем доступ к паре под индексом 1
	//firstindiceMaBuy := dat["1"].(map[string]interface{})["maBuy"].(float64)
	//firstindiceMaSell := dat["1"].(map[string]interface{})["maSell"].(float64)
	//firstindiceTiBuy := dat["1"].(map[string]interface{})["tiBuy"].(float64)
	//firstindiceTiSell := dat["1"].(map[string]interface{})["tiSell"].(float64)
	////str1 := firstIndice_maBuy[0].(string)
	//fmt.Println("firstindiceMaBuy = ", firstindiceMaBuy)
	//fmt.Println("firstindiceMaSell = ", firstindiceMaSell)
	//fmt.Println("firstindiceTiBuy = ", firstindiceTiBuy)
	//fmt.Println("firstindiceTiSell = ", firstindiceTiSell)
	//
	//// Получаем доступ к паре под индексом 2
	//secondindiceMaBuy := dat["2"].(map[string]interface{})["maBuy"].(float64)
	//secondindiceMaSell := dat["2"].(map[string]interface{})["maSell"].(float64)
	//secondindiceTiBuy := dat["2"].(map[string]interface{})["tiBuy"].(float64)
	//secondindiceTiSell := dat["2"].(map[string]interface{})["tiSell"].(float64)
	////str1 := secondIndice_maBuy[0].(string)
	//fmt.Println("secondindiceMaBuy = ", secondindiceMaBuy)
	//fmt.Println("secondindiceMaSell = ", secondindiceMaSell)
	//fmt.Println("secondindiceTiBuy = ", secondindiceTiBuy)
	//fmt.Println("secondindiceTiSell = ", secondindiceTiSell)
	//
	//// Получаем доступ к паре под индексом 3
	//parsePair(dat["3"])
	//
	//fmt.Println("dat[\"5\"]")
	//parsePair(dat["5"])
	//fmt.Println("dat[\"7\"]")
	//parsePair(dat["7"])
	//fmt.Println("dat[\"9\"]")
	//parsePair(dat["9"])
	//fmt.Println("dat[\"10\"]")
	//parsePair(dat["10"])

	// Какие индексы вообще есть
	// 1,2,3,5,7,9,10
	return allNewSignalsForThisTimeframe
}
