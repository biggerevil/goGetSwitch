package main

import (
	"goGetSwitch/getAndParseData"
	"log"
	"sync"
	"time"
)

type investingResponse struct {
	Page   int
	Fruits []string
}

func dataGetterAndParser(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	getAndParseData.GetAndParseData(url)
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

	for _, timeframe := range timeframes {
		wg.Add(1)
		go dataGetterAndParser(baseUrl + timeframe, &wg)
	}

	// Ждём окончания работы всех горутин
	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Program took %s", elapsed)
}
