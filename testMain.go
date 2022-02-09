package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"goGetSwitch/dbFunctions"
	"goGetSwitch/signal"
	"log"
	"strconv"
	"time"
)

func addStartDateStringToStake(stakesCollection *mongo.Collection, idOfsignaltoupdate primitive.ObjectID, formattedDateString string) {
	result, err := stakesCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": idOfsignaltoupdate},
		bson.D{
			{"$set", bson.D{{"StartDateString", formattedDateString}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
}

const dateFormat = "2006-01-02T15:04:05"

func formattedDateFromUnixTimestamp(incomingUnixTimestamp string) string {
	fmt.Println("incomingUnixTimestamp =", incomingUnixTimestamp)

	i, err := strconv.ParseInt(incomingUnixTimestamp, 10, 64)
	if err != nil {
		panic(err)
	}

	tm := time.Unix(i, 0)
	// Превращаем строку в UTC
	tm = tm.UTC()
	formattedDate := tm.Format(dateFormat)
	fmt.Println("formattedDate = ", formattedDate)

	return formattedDate
}

func main() {
	start := time.Now()

	// Для тестирования
	//unixTimestamp := "1641020400"
	//fmt.Println("unixTimestamp = ", unixTimestamp)
	//formattedDateFromUnixTimestamp(unixTimestamp)

	collection := dbFunctions.ConnectToDB()

	filter := bson.D{{"StartDateString", nil}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var stake signal.Signal
		if err = cursor.Decode(&stake); err != nil {
			log.Fatal(err)
		}

		fmt.Println("stake = ", signal.SignalDataInOneStringWithComments(stake))
		//fmt.Println("stake.StartUnixTimestamp = ", stake.StartUnixTimestamp)
		startUnixTimestampInString := strconv.FormatInt(stake.StartUnixTimestamp, 10)
		//fmt.Println("startUnixTimestampInString = ", startUnixTimestampInString)

		formattedStartDateString := formattedDateFromUnixTimestamp(startUnixTimestampInString)
		addStartDateStringToStake(collection, stake.ID, formattedStartDateString)
	}

	elapsed := time.Since(start)
	log.Printf("Done in %s", elapsed)
}
