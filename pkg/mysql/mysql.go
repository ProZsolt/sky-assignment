package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/ProZsolt/sky-assignment/pkg/metric"
)

type MySQL struct {
	database *sql.DB
}

func New(host string, username string, password string, database string) (MySQL, error) {
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+host+")/"+database)
	if err != nil {
		return MySQL{}, err
	}
	return MySQL{database: db}, nil
}

func (mySQL MySQL) LoadMetrics(from int, to int) ([]metric.Metric, error) {
	metrics := []metric.Metric{}
	query := "SELECT timestamp, cpuLoad, concurrency FROM metrics WHERE timestamp >= ? AND timestamp <= ?"
	stmt, err := mySQL.database.Prepare(query)
	if err != nil {
		return metrics, err
	}
	rows, err := stmt.Query(from, to)
	if err != nil {
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		var metric metric.Metric
		err := rows.Scan(&metric.Timestamp, &metric.CPULoad, &metric.Concurrency)
		if err != nil {
			return metrics, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (mySQL MySQL) SaveMetric(metric metric.Metric) error {
	query := "INSERT INTO metrics VALUES (?, ?, ?)"
	stmt, err := mySQL.database.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(metric.Timestamp, metric.CPULoad, metric.Concurrency)
	return err
}
