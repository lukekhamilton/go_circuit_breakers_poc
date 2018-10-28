package main

import (
	"fmt"
	"io/ioutil"
	"time"

	circuit "github.com/rubyist/circuitbreaker"
)

func main() {
	// url := "https://www.google.com/robots.txt"
	url := "http://www.google.com:81/"
	client := circuit.NewHTTPClient(time.Second*5, 10, nil)

	start := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	end := time.Now()
	elapsed := end.Sub(start)

	fmt.Printf("%+v\n", elapsed)
	// fmt.Printf("%+v\n", resp)

	var responseBytes []byte
	if resp != nil {
		responseBytes, _ = ioutil.ReadAll(resp.Body)
	}
	fmt.Println(string(responseBytes))

}
