package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Bid struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/quotation", nil)
	if err != nil {
		panic(err)
	}

	dollar := "USD"
	real := "BRL"
	// query args
	args := req.URL.Query()
	args.Add("src", dollar)
	args.Add("dst", real)
	// assign encoded query string to http request
	req.URL.RawQuery = args.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	saveToFile(body)
}

func saveToFile(data []byte) {
	f, err := os.Create("quotation.txt")
	if err != nil {
		panic(err)
	}

	var b Bid
	err = json.Unmarshal(data, &b)
	if err != nil {
		panic(err)
	}

	_, err = f.WriteString("Dollar: " + b.Bid)
	if err != nil {
		panic(err)
	}
	f.Close()
}
