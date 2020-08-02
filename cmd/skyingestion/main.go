package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ProZsolt/sky-assignment/pkg/metric"
	"github.com/ProZsolt/sky-assignment/pkg/mysql"
)

func main() {
	host := os.Getenv("SKY_DB_HOST")
	username := os.Getenv("SKY_DB_USERNAME")
	password := os.Getenv("SKY_DB_PASSWORD")
	database := os.Getenv("SKY_DB_DATABASE")

	db, err := mysql.New(host, username, password, database)
	if err != nil {
		fmt.Printf("Can't open database: %v\n", err)
		os.Exit(1)
	}

	now := time.Now()
	for i := 5; i >= 0; i-- {
		timestamp := now.Add(-(time.Duration(i) * time.Minute))
		m := metric.Metric{
			Timestamp:   timestamp.Unix(),
			CPULoad:     rand.Float32() * 100,
			Concurrency: rand.Intn(500001),
		}
		err = db.SaveMetric(m)
		if err != nil {
			fmt.Printf("Can't save %v: %v\n", m, err)
			os.Exit(1)
		}
	}
	fmt.Println("6 metric saved to database")
}
