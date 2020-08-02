package mysql

import (
	"os"
	"testing"

	"github.com/ProZsolt/sky-assignment/pkg/metric"
)

func TestSaveAndLoadMetrics(t *testing.T) {
	host := os.Getenv("SKY_TEST_DB_HOST")
	username := os.Getenv("SKY_TEST_DB_USERNAME")
	password := os.Getenv("SKY_TEST_DB_PASSWORD")
	database := os.Getenv("SKY_TEST_DB_DATABASE")
	db, err := New(host, username, password, database)
	if err != nil {
		t.Fatalf("Can't open database: %v", err)
	}

	metric1 := metric.Metric{
		Timestamp:   1500000000,
		CPULoad:     50.1,
		Concurrency: 100000,
	}
	metric2 := metric.Metric{
		Timestamp:   1500000060,
		CPULoad:     51.2,
		Concurrency: 200000,
	}
	metric3 := metric.Metric{
		Timestamp:   1500000120,
		CPULoad:     52.3,
		Concurrency: 300000,
	}
	err = db.SaveMetric(metric1)
	if err != nil {
		t.Fatalf("Can't save metric1: %v", err)
	}
	err = db.SaveMetric(metric2)
	if err != nil {
		t.Fatalf("Can't save metric2: %v", err)
	}
	err = db.SaveMetric(metric3)
	if err != nil {
		t.Fatalf("Can't save metric3: %v", err)
	}
	metrics, err := db.LoadMetrics(1500000060, 1500000120)
	if err != nil {
		t.Fatalf("Can't load metrics: %v", err)
	}
	if len(metrics) != 2 || metric2 != metrics[0] || metric3 != metrics[1] {
		t.Fatalf("Expected 2 metrics: %v %v, but got these: %v", metric2, metric3, metrics)
	}
}
