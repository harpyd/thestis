package app

type MetricsService interface {
	IncRequestsCount(status, method, path string)
}
