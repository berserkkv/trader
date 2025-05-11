package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Trade represents a trade entry
type Trade struct {
	ID     int     `json:"id"`
	Symbol string  `json:"symbol"`
	Action string  `json:"action"`
	Price  float64 `json:"price"`
	Amount int     `json:"amount"`
}

var db *sql.DB

// Initialize the database and create a table
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./trades.db")
	if err != nil {
		log.Fatal(err)
	}
	createTable := `
    CREATE TABLE IF NOT EXISTS trades (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        symbol TEXT,
        action TEXT,
        price REAL,
        amount INTEGER
    );
    `
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

// Handle fetching trades from the database
func getTrades(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, symbol, action, price, amount FROM trades")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var trades []Trade
	for rows.Next() {
		var trade Trade
		err := rows.Scan(&trade.ID, &trade.Symbol, &trade.Action, &trade.Price, &trade.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		trades = append(trades, trade)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trades)
}

// Handle adding a new trade
func addTrade(w http.ResponseWriter, r *http.Request) {
	var trade Trade
	err := json.NewDecoder(r.Body).Decode(&trade)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statement, err := db.Prepare("INSERT INTO trades (symbol, action, price, amount) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = statement.Exec(trade.Symbol, trade.Action, trade.Price, trade.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Trade added successfully")
}

// Serve static files
func serveFiles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "fr/index.html")
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/", serveFiles)
	http.HandleFunc("/api/trades", getTrades)
	http.HandleFunc("/api/add_trade", addTrade)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
