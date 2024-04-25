package utils

import (
	"fmt"
	"net/http"
	"time"
)

type RetryFunc func(string) (*http.Response, error)

func HTTPWithRetry(f RetryFunc, url string) (*http.Response, error) {
	count := 10
	var resp *http.Response
	var err error
	for i := 0; i < count; i++ {
		resp, err = f(url)
		if err != nil {
			fmt.Printf("Error calling url %v\n", url)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	return resp, err
}
