package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

	err = saveToDataBase(&q)
	if error != nil {
		return nil, error
	}

	return &q, nil
}

func saveToDataBase(quotation *Quotation) error {
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
	err = insertQuotation(db, quotation.Usdbrl.Bid)
	if err != nil {
		return err
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
