package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/quotation", nil)
	if err != nil {
		panic(err)
	}

	// query args
	dollar := "USD"
	real := "BRL"
	args := req.URL.Query()
	args.Add("src", dollar)
	args.Add("dst", real)

	// assign encoded query string to http request
	req.URL.RawQuery = args.Encode()

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Status)
	fmt.Println(string(body))
}
