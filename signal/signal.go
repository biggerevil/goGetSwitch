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

	Pairname      string  `bson:"Pairname" json:"Pairname"`
	CurrentPrice  float64 `bson:"CurrentPrice" json:"CurrentPrice"`
	Timeframe     int     `bson:"Timeframe" json:"Timeframe"`
	UnixTimestamp int64   `bson:"UnixTimestamp" json:"UnixTimestamp"`
}

func SignalDataInOneString(signal Signal) string {
	delimiter := ", "
	stringToReturn := signal.Pairname + delimiter +
		strconv.Itoa(signal.MaBuy) + delimiter + strconv.Itoa(signal.MaSell) +
		strconv.Itoa(signal.TiBuy) + delimiter + strconv.Itoa(signal.TiSell) +
		delimiter + fmt.Sprintf("%f", signal.CurrentPrice)

	return stringToReturn
}

func SignalDataInOneStringWithComments(signal Signal) string {
	stringToReturn := "Pairname: " + signal.Pairname +
		", MaBuy: " + strconv.Itoa(signal.MaBuy) + ", MaSell: " + strconv.Itoa(signal.MaSell) +
		", TiBuy: " + strconv.Itoa(signal.TiBuy) + ", TiSell: " + strconv.Itoa(signal.TiSell) +
		", CurrentPrice: " + fmt.Sprintf("%f", signal.CurrentPrice)

	return stringToReturn
}
