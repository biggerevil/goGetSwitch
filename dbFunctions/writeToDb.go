package dbFunctions

import (
	"context"
	"fmt"
	"goGetSwitch/signal"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connection URI
const uri = "mongodb://localhost"

func WriteNewSignalsToDB(allNewSignals []signal.Signal) {
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	// Проверяем наличие ошибки
	if err != nil {
		panic(err)
	}
	// Пишем, что произойдёт, когда ф-я WriteNewSignalsToDB закончится
	defer func() {
		// Отключаемся от mongo
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping the primary
	// Проверяем, что нам удалось подключиться
	log.Println("Собираюсь сделать Ping")
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("Закончил с вызовом Ping")

	fmt.Println("Successfully connected and pinged.")

	//
	// Inserting
	//
	// Выбираем нашу базу данных и коллекцию в ней. Я назвал их history_stakes и stakes.
	coll := client.Database("history_stakes").Collection("stakes")
	// Example of inserting
	//docs := []interface{}{
	//	bson.D{{"title", "My Brilliant Friend"}, {"author", "Elena Ferrante"}, {"year_published", 2012}},
	//	bson.D{{"title", "Lucy"}, {"author", "Jamaica Kincaid"}, {"year_published", 2002}},
	//	bson.D{{"title", "Cat's Cradle"}, {"author", "Kurt Vonnegut Jr."}, {"year_published", 1998}},
	//}

	// В массиве docs у нас будут храниться все новые сигналы.
	// Просто они будут храниться по идее в хорошем формате, как нам нужно.
	var docs []interface{}
	// По очереди берём каждый сигнал из переданного в ф-ю массива и добавляем в docs.
	for _, newSignal := range allNewSignals {
		//signalInJson, _ := json.Marshal(signal)
		//signalInBson := bson.D(signal)
		docs = append(docs, newSignal)
		//bson.UnmarshalJSON([]byte(`{"id": 1,"name": "A green door","price": 12.50,"tags": ["home", "green"]}`),&bdoc)
	}

	// Вызываем саму функцию InsertMany, то есть добавляем все новые сигналы в БД.
	// (Эту функцию писал не я, она идёт от разработчиков библиотеки для работы с mongo)
	log.Println("[WriteNewSignalsToDB] Собираюсь позвать InsertMany")
	result, err := coll.InsertMany(context.TODO(), docs)
	// (Заметь, что результат работы InsertMany сохраняется в переменной result)

	// Проверяем, не вылетела ли ошибка при InsertMany
	if err != nil {
		fmt.Println("Error - InsertMany returned error - ", err)
	}

	// Просто для дополнительной ручной проверки ошибок я вывожу все ID добавленных сигналов.
	// От этого можно в принципе избавиться, мне это уже не нужно особо.
	list_ids := result.InsertedIDs
	fmt.Printf("Documents inserted: %v\n", len(list_ids))
	for _, id := range list_ids {
		fmt.Printf("Inserted document with _id: %v\n", id)
	}

	// Закрытие соединения (такой ручной вызов вызывает ошибку. Похоже,
	// это делается автоматически) ((Это давний комментарий.))
	//log.Println("Собираюсь вызвать Disconnect")
	//err = client.Disconnect(context.TODO())
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Connection to MongoDB closed.")
}
