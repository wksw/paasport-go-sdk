package main

import (
	"log"
	"net/http"

	paasport "gitee.com/paasport/go-sdk"
)

func main() {
	client, err := paasport.NewClient("123", "123", "http://127.0.0.1:9091",
		paasport.WithTimeout(paasport.HTTPTimeout{}))
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := client.Do(http.MethodGet, "/", nil, nil); err != nil {
		log.Fatal(err.Message)
	}
}
