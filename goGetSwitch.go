package main

import (
	"fmt"
	"goGetSwitch/dbFunctions"
	"goGetSwitch/getAndParseData"
	"goGetSwitch/signal"
	"log"
	"sync"
	"time"
)

// Можно использовать структуру для парсинга ответа от investing
type investingResponse struct {
	Page   int
	Fruits []string
}

func dataGetterAndParser(baseUrl string, timeframe string, unixTimestamp int64, wg *sync.WaitGroup, channelForSendingSignalsArrays chan []signal.Signal) {
	defer wg.Done()

	newSignalsForThisTimeframe := getAndParseData.GetAndParseData(baseUrl, timeframe, unixTimestamp)

	channelForSendingSignalsArrays <- newSignalsForThisTimeframe
}

func main() {
	start := time.Now()

	// Для ожидания завершения горутин
	var wg sync.WaitGroup

	baseUrl := "https://www.investing.com/common/technical_summary/api.php?action=TSB_updatePairs&pairs=1,2,3,5,7,9,10&timeframe="

	// Вариант с одним URL
	//respBody := getAndParseData.GetData(url)
	//maBuy := getAndParseData.ParseData(respBody)

	timeframes := []string{"300", "900", "1800", "3600", "18000", "86400"}

	// Определяю единый timestamp для всех сигналов, чтобы не было сигналов с timestamp, отличающихся на
	// несколько секунд
	unixTimestamp := time.Now().Unix()

	channelForGettingSignalsArray := make(chan []signal.Signal)
	var allNewSignals []signal.Signal

	for _, timeframe := range timeframes {
		wg.Add(1)
		go dataGetterAndParser(baseUrl, timeframe, unixTimestamp, &wg, channelForGettingSignalsArray)
	}

	// Получение данных из канала
	// TODO: этот код не рассчитывает, что данные откуда-либо могут не вернуться.
	//	Я думаю, это не совсем корректный способ получения данных
	for i := 0; i < len(timeframes); i++ {
		newSignals := <-channelForGettingSignalsArray
		allNewSignals = append(allNewSignals, newSignals...)
		fmt.Println("[While working] len(allNewSignals) = ", len(allNewSignals))
	}

	// Ждём окончания работы всех горутин (этот код написал ДО использования каналов). Возможно,
	// этот код уже не нужен
	wg.Wait()

	fmt.Println("len(allNewSignals) = ", len(allNewSignals))

	log.Println("Собираюсь позвать dbFunctions.WriteData")
	dbFunctions.WriteData(allNewSignals)
	log.Println("Закончил с вызовом dbFunctions.WriteData")

	elapsed := time.Since(start)
	log.Printf("Program took %s", elapsed)

	// TODO: для добавления финальной цены можно делать поиск по БД по полю EndUnixTimestamp, передавая в него
	// 	ТЕКУЩИЙ timestamp. И добавлять в таким сигналам finalPrice и finalPriceIsMoreThanInitial. Для более удобной
	//	 работы можно завести отдельный, не знаю, массив, в котором будут храниться цены всех пар (всего 7-10 штук,
	//	 сколько там пар), а браться этот массив будет из, не знаю, массива с сигналами, например. И для добавления к
	//	 сигналам финальной будут браться данные из этого массива (так, я думаю, будет удобнее, чем для каждого сигнала
	//	 искать финальную цену среди всего массива сигналов, так как будет много прогонов).
	//	И эту функцию можно запускать в районе конца программы, чтобы после работы GetWinrate и SwitchAlgorithm
}
