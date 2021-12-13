package getAndParseData

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetData(url string) []byte {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", url, nil)

	// Headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36")
	req.Header.Add("Sec-Ch-Ua-Platform", "\"macOS\"")

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("string(respBody) = ", string(respBody))

	return respBody
}
