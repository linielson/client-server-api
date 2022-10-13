package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type QuotationResult struct {
	Bid string `json:"bid"`
}

type Quotation struct {
	Usdbrl QuotationResult `json:"usdbrl"`
}

func main() {
	http.HandleFunc("/quotation", quotationHandler)
	http.ListenAndServe(":8080", nil)
}

func quotationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/quotation" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	srcParam := r.URL.Query().Get("src")
	dstParam := r.URL.Query().Get("dst")
	if srcParam == "" || dstParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	quotation, err := requestQuotation(srcParam, dstParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quotation.Usdbrl)
}

func requestQuotation(currencySrc, currencyDst string) (*Quotation, error) {
	res, err := http.Get("https://economia.awesomeapi.com.br/json/last/" + currencySrc + "-" + currencyDst)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, error := ioutil.ReadAll(res.Body)
	if error != nil {
		return nil, error
	}
	var q Quotation
	error = json.Unmarshal(body, &q)
	if error != nil {
		return nil, error
	}
	return &q, nil
}
