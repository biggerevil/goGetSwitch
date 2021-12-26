package signal

import (
	"fmt"
	"strconv"
)

type Signal struct {
	// Объявляю на новых строчках (а не в одну) для наглядности
	MaBuy  int `bson:"MaBuy" json:"MaBuy"`
	MaSell int `bson:"MaSell" json:"MaSell"`
	TiBuy  int `bson:"TiBuy" json:"TiBuy"`
	TiSell int `bson:"TiSell" json:"TiSell"`

	Pairname           string  `bson:"Pairname" json:"Pairname"`
	CurrentPrice       float64 `bson:"CurrentPrice" json:"CurrentPrice"`
	Timeframe          int     `bson:"Timeframe" json:"Timeframe"`
	StartUnixTimestamp int64   `bson:"StartUnixTimestamp" json:"StartUnixTimestamp"`
	EndUnixTimestamp   int64   `bson:"EndUnixTimestamp" json:"EndUnixTimestamp"`

	// Служебные поля
	ID string `bson:"_id" json:"_id"`
}

func SignalDataInOneStringWithComments(signal Signal) string {
	stringToReturn := "ID: " + signal.ID +
		" , Pairname: " + signal.Pairname +
		" , MaBuy: " + strconv.Itoa(signal.MaBuy) + ", MaSell: " + strconv.Itoa(signal.MaSell) +
		" , TiBuy: " + strconv.Itoa(signal.TiBuy) + ", TiSell: " + strconv.Itoa(signal.TiSell) +
		" , CurrentPrice: " + fmt.Sprintf("%f", signal.CurrentPrice) +
		" , StartUnixTimestamp: " + fmt.Sprintf("%v", signal.StartUnixTimestamp) +
		" , EndUnixTimestamp: " + fmt.Sprintf("%v", signal.EndUnixTimestamp) +
		" , Timeframe: " + fmt.Sprintf("%v", signal.Timeframe)

	return stringToReturn
}
