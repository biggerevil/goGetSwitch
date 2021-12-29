package signal

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	// По поводу omitempty см. комментарий чуть ниже
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	/*
		По поводу omitempty:

		The "omitempty" option specifies that the field should be omitted from the encoding if the field has an empty
		value, defined as false, 0, a nil pointer, a nil interface value, and any empty array, slice, map, or string.
		Чуть подробнее - https://stackoverflow.com/a/49043598/8604912
	*/
}

func SignalDataInOneStringWithComments(signal Signal) string {
	stringToReturn := "Pairname: " + signal.Pairname +
		//stringToReturn := "ID: " + signal.ID.String() +
		//	" , Pairname: " + signal.Pairname +
		" , MaBuy: " + strconv.Itoa(signal.MaBuy) + ", MaSell: " + strconv.Itoa(signal.MaSell) +
		" , TiBuy: " + strconv.Itoa(signal.TiBuy) + ", TiSell: " + strconv.Itoa(signal.TiSell) +
		" , CurrentPrice: " + fmt.Sprintf("%f", signal.CurrentPrice) +
		" , StartUnixTimestamp: " + fmt.Sprintf("%v", signal.StartUnixTimestamp) +
		" , EndUnixTimestamp: " + fmt.Sprintf("%v", signal.EndUnixTimestamp) +
		" , Timeframe: " + fmt.Sprintf("%v", signal.Timeframe)

	return stringToReturn
}