package getAndParseData

import "fmt"

func GetAndParseData(url string) string {
	respBody := GetData(url)
	maBuy := ParseData(respBody)
	fmt.Println("Gonna return maBuy = ", maBuy)
	return maBuy
}
