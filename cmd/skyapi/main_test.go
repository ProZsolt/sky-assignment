package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ProZsolt/sky-assignment/pkg/metric"
)

func checkCodeAndBody(
	t *testing.T,
	handlerfunc func(http.ResponseWriter, *http.Request),
	request *http.Request,
	code int,
	body string,
) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerfunc)
	handler.ServeHTTP(rr, request)
	if rr.Code != code {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, code)
	}
	if rr.Body.String() != body {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), body)
	}
}

type metricLoaderMock struct {
	t    *testing.T
	from int
	to   int
	met  []metric.Metric
	err  error
}

func (mlm metricLoaderMock) LoadMetrics(from int, to int) ([]metric.Metric, error) {
	if mlm.from != from {
		mlm.t.Errorf("wrong query parameter 'from': got %v want %v", from, mlm.from)
	}
	if mlm.to != to {
		mlm.t.Errorf("wrong query parameter 'to': got %v want %v", to, mlm.to)
	}
	return mlm.met, mlm.err
}

func appWithDBMock(t *testing.T, from int, to int, met []metric.Metric, err error) app {
	return app{
		db: metricLoaderMock{
			t:    t,
			from: from,
			to:   to,
			met:  met,
			err:  err,
		},
	}
}

func TestAPI(t *testing.T) {
	metrics := []metric.Metric{
		{
			Timestamp:   1500000000,
			CPULoad:     50.1,
			Concurrency: 100000,
		},
		{
			Timestamp:   1500000060,
			CPULoad:     51.2,
			Concurrency: 200000,
		},
	}
	a := appWithDBMock(t, 1500000000, 1500000060, metrics, nil)
	req, err := http.NewRequest(http.MethodGet, "/api?from=1500000000&to=1500000060", nil)
	if err != nil {
		t.Fatal(err)
	}
	checkCodeAndBody(t, a.api, req, http.StatusOK, "[{\"timestamp\":1500000000,\"cpuLoad\":50.1,\"concurrency\":100000},{\"timestamp\":1500000060,\"cpuLoad\":51.2,\"concurrency\":200000}]\n")
}

func TestMethodNotAllowed(t *testing.T) {
	a := appWithDBMock(t, 10, 100, []metric.Metric{}, nil)
	req, err := http.NewRequest(http.MethodPost, "/api?from=10&to=100", nil)
	if err != nil {
		t.Fatal(err)
	}
	checkCodeAndBody(t, a.api, req, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed)+"\n")
}

var parametertests = []struct {
	url  string
	text string
}{
	{"/api?to=100", "{\"error\": \"missing or empty 'from' query parameter\"}\n"},
	{"/api?from=&to=100", "{\"error\": \"missing or empty 'from' query parameter\"}\n"},
	{"/api?from=10", "{\"error\": \"missing or empty 'to' query parameter\"}\n"},
	{"/api?from=10&to=", "{\"error\": \"missing or empty 'to' query parameter\"}\n"},
	{"/api?from=asdf&to=100", "{\"error\": \"query parameter 'from' is not an integer\"}\n"},
	{"/api?from=10&to=asdf", "{\"error\": \"query parameter 'to' is not an integer\"}\n"},
}

func TestParameters(t *testing.T) {
	for _, tt := range parametertests {
		t.Run(tt.url, func(t *testing.T) {
			a := appWithDBMock(t, 10, 100, []metric.Metric{}, nil)
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			checkCodeAndBody(t, a.api, req, http.StatusBadRequest, tt.text)
		})
	}
}

func TestDBError(t *testing.T) {
	a := appWithDBMock(t, 1500000000, 1500000060, []metric.Metric{}, fmt.Errorf("DB error"))
	req, err := http.NewRequest(http.MethodGet, "/api?from=1500000000&to=1500000060", nil)
	if err != nil {
		t.Fatal(err)
	}
	checkCodeAndBody(t, a.api, req, http.StatusInternalServerError, "{\"error\": \"internal server error\"}\n")
}
