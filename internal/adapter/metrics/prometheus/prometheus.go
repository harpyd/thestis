package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsService struct {
	requestsTotal prometheus.Counter
	requests      *prometheus.CounterVec
}

func NewMetricsService() (*MetricsService, error) {
	requestsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "thestis_requests_total",
		Help: "Total number of requests",
	})

	if err := prometheus.Register(requestsTotal); err != nil {
		return nil, err
	}

	requests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "thestis_specified_requests_total",
			Help: "Total number of requests with labels",
		},
		[]string{"status", "method", "path"},
	)

	if err := prometheus.Register(requests); err != nil {
		return nil, err
	}

	return &MetricsService{
		requestsTotal: requestsTotal,
		requests:      requests,
	}, nil
}

func (m *MetricsService) IncRequestsCount(status, method, path string) {
	m.requestsTotal.Inc()
	m.requests.WithLabelValues(status, method, path).Inc()
}
