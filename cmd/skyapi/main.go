package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ProZsolt/sky-assignment/pkg/metric"
	"github.com/ProZsolt/sky-assignment/pkg/mysql"
)

type app struct {
	db metric.MetricLoader
}

func (a app) parseRequest(r *http.Request) (int, int, error) {
	fromStr := r.URL.Query().Get("from")
	if fromStr == "" {
		return 0, 0, fmt.Errorf("missing or empty 'from' query parameter")
	}
	toStr := r.URL.Query().Get("to")
	if toStr == "" {
		return 0, 0, fmt.Errorf("missing or empty 'to' query parameter")
	}
	from, err := strconv.Atoi(fromStr)
	if err != nil {
		return 0, 0, fmt.Errorf("query parameter 'from' is not an integer")
	}
	to, err := strconv.Atoi(toStr)
	if err != nil {
		return 0, 0, fmt.Errorf("query parameter 'to' is not an integer")
	}

	return from, to, nil
}

func (a app) api(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling request: %v", r)
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	from, to, err := a.parseRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err), http.StatusBadRequest)
		return
	}
	metrics, err := a.db.LoadMetrics(from, to)
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\": \"internal server error\"}"), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(metrics)
}

func main() {
	host := os.Getenv("SKY_DB_HOST")
	username := os.Getenv("SKY_DB_USERNAME")
	password := os.Getenv("SKY_DB_PASSWORD")
	database := os.Getenv("SKY_DB_DATABASE")

	db, err := mysql.New(host, username, password, database)
	if err != nil {
		log.Fatalf("can't open database: %v", err)
	}
	a := app{db: db}

	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/api", a.api)
	http.ListenAndServe(":8080", nil)
}
