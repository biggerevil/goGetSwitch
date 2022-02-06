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

/*
	Эта функция принимает JSON и достаёт из него:
	1. Значения индикаторов (MaBuy, MaSell, TiBuy, TiSell);
	2. Название пары;
	3. Текущую цену.
	И после этого эта ф-я создаёт структуру signal.Signal, заполняем её полученным данными и
	возвращает эту структуру.
	TODO: Сменить названые на parseJSON или что-либо такое.
*/
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

/*
	Это наша "входная" функция, при помощи которой мы парсим данные.
	Эта функция принимает тело ответа в JSON (в []byte, если быть точнее, но внутри там json), и
	возвращает массив сигналов типа signal.Signal.
*/
func ParseData(respBody []byte, timeframe int, unixTimestamp int64) []signal.Signal {
	/*
		Работа с JSON
	*/
	// Создаём словарь, в котором мы будем хранить распаршенные из JSON данные.
	var dat map[string]interface{}

	// При помощи встроенной функции json.Unmarshal() парсим JSON, чтобы получить структуру, с которой можем уже работать
	// "на уровне языка программирования". И сразу же проверяем, были ли ошибки при парсинге JSON.
	if err := json.Unmarshal(respBody, &dat); err != nil {
		panic(err)
	}
	// На всякий случай (чтобы в случае ошибок можно было посмотреть логи) выводим распаршенные данные.
	// TODO: сменить fmt на log
	fmt.Println("\n\n dat:")
	fmt.Println(dat)

	// Создаём массив, в котором будем хранить распаршенные сигналы.
	var allNewSignalsForThisTimeframe []signal.Signal

	// У нас есть массив индексов пар. В JSON данные по парам лежат под этими индексами. Поэтому
	// мы итерируемся по этим индексам и парсим данные.
	for _, indice := range getIndices() {
		// Достаём из JSON название пары, текущую цену и индикаторы. Из этой же функции получаем
		// почти готовый объект signal.Signal.
		newSignal := parsePair(dat[indice])
		// Добавляем unixTimestamp, который мы передавали в ParseData()
		newSignal.Timeframe = timeframe
		// Добавляем unixTimestamp, который мы передавали в ParseData(). (У этого отдельная логика, спросите меня).
		newSignal.StartUnixTimestamp = unixTimestamp
		// Вычисляем EndUnixTimestamp, то есть время, когда ставка должна закончиться.
		newSignal.EndUnixTimestamp = unixTimestamp + int64(timeframe)
		// Добавляем этот сигнал в массив сигналов.
		allNewSignalsForThisTimeframe = append(allNewSignalsForThisTimeframe, newSignal)
		fmt.Println("newSignal data with comments: ", signal.SignalDataInOneStringWithComments(newSignal))
	}

	// TODO: удалить эти комментарии. Это был тестовый код.
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
