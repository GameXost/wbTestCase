package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	CacheHits = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_hits_total",
		Help: "total number of cache hits",
	})
	CacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_misses_total",
		Help: "total number of cache misses",
	})
	RequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "total number of http requests",
		},
	)

	RequestsSuccess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_success",
			Help: "total number of http requests success returned",
		},
	)

	RequestsNotFound = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_NotFound",
			Help: "total number of responses not found",
		},
	)

	RequestsBadRequest = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_bad_requests",
			Help: "total number of bad requests",
		},
	)

	RequestsServerError = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_serv_err",
			Help: "total number of http requests failed due to server error",
		},
	)
)
