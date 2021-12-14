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

func WriteData(allNewSignals []signal.Signal) {
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping the primary
	log.Println("Собираюсь сделать Ping")
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("Закончил с вызовом Ping")

	fmt.Println("Successfully connected and pinged.")

	//
	// Inserting
	//
	coll := client.Database("history_stakes").Collection("stakes")
	// Example of inserting
	//docs := []interface{}{
	//	bson.D{{"title", "My Brilliant Friend"}, {"author", "Elena Ferrante"}, {"year_published", 2012}},
	//	bson.D{{"title", "Lucy"}, {"author", "Jamaica Kincaid"}, {"year_published", 2002}},
	//	bson.D{{"title", "Cat's Cradle"}, {"author", "Kurt Vonnegut Jr."}, {"year_published", 1998}},
	//}

	var docs []interface{}
	for _, newSignal := range allNewSignals {
		//signalInJson, _ := json.Marshal(signal)
		//signalInBson := bson.D(signal)
		docs = append(docs, newSignal)
		//bson.UnmarshalJSON([]byte(`{"id": 1,"name": "A green door","price": 12.50,"tags": ["home", "green"]}`),&bdoc)
	}

	log.Println("[WriteData] Собираюсь позвать InsertMany")
	result, err := coll.InsertMany(context.TODO(), docs)

	if err != nil {
		fmt.Println("Error - InsertMany returned error - ", err)
	}

	list_ids := result.InsertedIDs
	fmt.Printf("Documents inserted: %v\n", len(list_ids))
	for _, id := range list_ids {
		fmt.Printf("Inserted document with _id: %v\n", id)
	}

	// Закрытие соединения (такой ручной вызов вызывает ошибку. Похоже,
	// это делается автоматически)
	//log.Println("Собираюсь вызвать Disconnect")
	//err = client.Disconnect(context.TODO())
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Connection to MongoDB closed.")
}
