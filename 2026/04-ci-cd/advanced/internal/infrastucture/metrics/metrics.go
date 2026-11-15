package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
		},
		[]string{"method", "pattern", "status"},
	)

	HttpRequestsDurationHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds_histogram",
			Buckets: []float64{
				0.1,  // 100 ms
				0.2,  // 200 ms
				0.25, // 250 ms
				0.5,  // 500 ms
				1,    // 1 s
			},
		},
		[]string{"method", "pattern"},
	)
)
