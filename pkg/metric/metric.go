package metric

type Metric struct {
	Timestamp   int64   `json:"timestamp"`
	CPULoad     float32 `json:"cpuLoad"`
	Concurrency int     `json:"concurrency"`
}

type MetricLoader interface {
	LoadMetrics(from int, to int) ([]Metric, error)
}

type MetricSaver interface {
	SaveMetric(metric Metric) error
}
