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
//type investingResponse struct {
//	Page   int
//	Fruits []string
//}

func dataGetterAndParser(baseUrl string, timeframe string, unixTimestamp int64, wg *sync.WaitGroup, channelForSendingSignalsArrays chan []signal.Signal) {
	defer wg.Done()

	newSignalsForThisTimeframe := getAndParseData.GetAndParseData(baseUrl, timeframe, unixTimestamp)

	// TODO: при наличии ошибки я могу возвращать эту ошибку (и впоследствии выводить её и продолжать работу)
	// 	вместо сигналов. Тогда мой код станет более надёжным перед лицом ошибок
	channelForSendingSignalsArrays <- newSignalsForThisTimeframe
}

func main() {
	// Замеряем время работы программы.
	start := time.Now()

	// Для ожидания завершения горутин
	// TODO: с добавлением каналов это по идее уже не очень нужно.
	//  Наверное, стоит делать разработку в отдельной ветке, и написать там много тестов, чтобы быть более спокойными,
	//  что мы ничего не сломали.
	var wg sync.WaitGroup

	// Это базовая ссылка, с которой мы получаем данные (эти данные потом записываем в БД).
	// К этой ссылке нужно только добавить время, на которое мы хотим получить
	// данные (на 5 минут (300 секунд), 15 минут (900 секунд) и так далее).
	// Данные мы получаем на все пары (номера пар (я не знаю, почему номера именно такие, просто вот
	// так вот сайт работает) передаются в параметре pairs, см. ссылку)
	baseUrl := "https://www.investing.com/common/technical_summary/api.php?action=TSB_updatePairs&pairs=1,2,3,5,7,9,10&timeframe="

	// Вариант с одним URL
	//respBody := getAndParseData.GetData(url)
	//maBuy := getAndParseData.ParseData(respBody)

	// TODO: добавить timeframe на 2 часа (то есть 7200 по идее) и проверить, что всё
	// 	корректно работает
	timeframes := []string{"300", "900", "1800", "3600", "7200", "18000", "86400"}

	// Определяю единый timestamp для всех сигналов, чтобы не было сигналов с timestamp, отличающихся на
	// несколько секунд
	currentUnixTimestamp := time.Now().Unix()

	// Канал, из которого мы будем получать сигналы из горутин
	channelForGettingSignalsArray := make(chan []signal.Signal)
	// Массив со всеми новыми сигналами
	var allNewSignals []signal.Signal

	// Запускаем горутины со всеми timeframe (время ставки)
	for _, timeframe := range timeframes {
		wg.Add(1)
		go dataGetterAndParser(baseUrl, timeframe, currentUnixTimestamp, &wg, channelForGettingSignalsArray)
	}

	// Получение данных из канала и добавление их в массив всех новых сигналов
	// TODO: этот код не рассчитывает, что данные откуда-либо могут не вернуться.
	//	Я думаю, это не совсем корректный способ получения данных
	for i := 0; i < len(timeframes); i++ {
		// Получение данных из канала
		newSignals := <-channelForGettingSignalsArray
		// Добавление новых сигналов в массив со всеми новыми сигналами
		allNewSignals = append(allNewSignals, newSignals...)
		fmt.Println("[While working] len(allNewSignals) = ", len(allNewSignals))
	}

	// Ждём окончания работы всех горутин (этот код написал ДО использования каналов). Возможно,
	// этот код уже не нужен
	wg.Wait()

	fmt.Println("len(allNewSignals) = ", len(allNewSignals))

	log.Println("Собираюсь позвать dbFunctions.WriteNewSignalsToDB")
	// Записываем данные в БД.
	dbFunctions.WriteNewSignalsToDB(allNewSignals)
	log.Println("Закончил с вызовом dbFunctions.WriteNewSignalsToDB")

	// Добавляем к "заканчивающимся" ставкам конечную цену
	// и информацию, больше ли конечная цена, чем начальная
	dbFunctions.UpdateEndingStakes(currentUnixTimestamp, allNewSignals)

	// Заканчиваем замер времени работы программы и выводим эту информацию.
	elapsed := time.Since(start)
	log.Printf("Program took %s", elapsed)
}
