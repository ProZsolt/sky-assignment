package metric

type Metric struct {
	Timestamp   int64
	CPULoad     float32
	Concurrency int
}

type MetricLoader interface {
	LoadMetrics(from int, to int) ([]Metric, error)
}

type MetricSaver interface {
	SaveMetric(metric Metric) error
}
