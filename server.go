package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
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
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/"+currencySrc+"-"+currencyDst, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
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

	err = saveToDataBase(context.Background(), &q)
	if error != nil {
		return nil, error
	}

	return &q, nil
}

func saveToDataBase(ctx context.Context, quotation *Quotation) error {
	os.Remove("sqlite-database.db")
	file, err := os.Create("sqlite-database.db")
	if err != nil {
		panic(err)
	}
	file.Close()

	db, err := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createTable(db)

	select {
	case <-time.After(10 * time.Millisecond):
		err = insertQuotation(db, quotation.Usdbrl.Bid)
		if err != nil {
			return err
		}
		log.Println("Insert Quotation successfully executed")
	case <-ctx.Done():
		log.Println("Failed to insertation Quotation")
	}

	displayQuotation(db)
	return nil
}

func createTable(db *sql.DB) {
	createTable := `CREATE TABLE quotation (
		"id" TEXT NOT NULL PRIMARY KEY,
		"code" TEXT,
		"codein" TEXT,
		"bid" TEXT
	);`

	stmt, err := db.Prepare(createTable)
	if err != nil {
		panic(err)
	}
	stmt.Exec()
}

func insertQuotation(db *sql.DB, bid string) error {
	stmt, err := db.Prepare(`INSERT INTO quotation(id, code, codein, bid) values(?, ?, ?, ?)`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(uuid.New(), "USB", "BRL", bid)
	if err != nil {
		return err
	}
	return nil
}

func displayQuotation(db *sql.DB) {
	row, err := db.Query("SELECT * FROM quotation ORDER BY bid")
	if err != nil {
		panic(err)
	}
	defer row.Close()
	for row.Next() {
		var id string
		var code string
		var codein string
		var bid string
		row.Scan(&id, &code, &codein, &bid)
		log.Println(id, " ", code, " ", codein, " ", bid)
	}
}
