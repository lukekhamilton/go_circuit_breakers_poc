package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

var cb *gobreaker.CircuitBreaker

func init() {
	var st gobreaker.Settings
	st.Name = "HTTP GET"
	st.Timeout = time.Second * 10 // This doesn't look like this working as the net/http default IdleConnTimeout 30 * Time.Second is what my results are showing
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		fmt.Printf("counts.Requests = %d, failureRatio = %v\n", counts.Requests, failureRatio)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}
	st.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
		fmt.Printf("Name: %v, state from: %v, state to: %v", name, from, to)
	}

	cb = gobreaker.NewCircuitBreaker(st)
}

// Get wraps http.get in a cb
func Get(url string) ([]byte, error) {
	body, err := cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	})
	if err != nil {
		return nil, err
	}

	return body.([]byte), nil
}

func main() {
	fmt.Printf("%+v\n", "Hello Circuit Breaker POC")

	// Thanks uncle google for a timeout url :)
	url := "http://www.google.com:81/"

	start := time.Now()
	fmt.Println("CB wrapped Get")
	_, err := Get(url)
	if err != nil {
		log.Println(err)
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Timer: %v\n\n", elapsed)
	// fmt.Println(string(body))

	start1 := time.Now()
	fmt.Println("http.Get")
	_, err1 := http.Get(url)
	if err1 != nil {
		log.Println(err)
	}
	t1 := time.Now()
	elapsed1 := t1.Sub(start1)
	fmt.Printf("Timer: %v\n\n", elapsed1)

}
