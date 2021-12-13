package signal

import (
	"fmt"
	"strconv"
)

type Signal struct {
	// Объявляю на новых строчках (а не в одну) для наглядности
	MaBuy int
	MaSell int
	TiBuy int
	TiSell int

	Pairname string
	CurrentPrice float64
}

func SignalDataInOneString(signal Signal) (string) {
	delimiter := ", "
	stringToReturn := signal.Pairname + delimiter +
		strconv.Itoa(signal.MaBuy) + delimiter + strconv.Itoa(signal.MaSell) +
		strconv.Itoa(signal.TiBuy) + delimiter + strconv.Itoa(signal.TiSell) +
		delimiter + fmt.Sprintf("%f", signal.CurrentPrice)

	return stringToReturn
}

func SignalDataInOneStringWithComments(signal Signal) (string) {
	stringToReturn := "Pairname: " + signal.Pairname +
		", MaBuy: " + strconv.Itoa(signal.MaBuy) + ", MaSell: " + strconv.Itoa(signal.MaSell) +
		", TiBuy: " + strconv.Itoa(signal.TiBuy) + ", TiSell: " + strconv.Itoa(signal.TiSell) +
		", CurrentPrice: " + fmt.Sprintf("%f", signal.CurrentPrice)

	return stringToReturn
}
