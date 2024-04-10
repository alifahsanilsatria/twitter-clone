package common

import (
	"github.com/prometheus/client_golang/prometheus"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "twitter_backend_http_requests_total",
		Help: "Number of requests received",
	},
	[]string{"path"},
)

var totalSuccessfulRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "twitter_backend_http_successful_requests_total",
		Help: "Number of successful requests",
	},
	[]string{"path"},
)

var totalFailedRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "twitter_backend_http_failed_requests_total",
		Help: "Number of failed requests",
	},
	[]string{"path"},
)

type PushMetricsToPrometheusParam struct {
	Path                string
	IsRequestSuccessful bool
}

func PushMetricsToPrometheus(param PushMetricsToPrometheusParam) {
	incrementTotalRequestReceivedPerEndpoint(param.Path)
	if param.IsRequestSuccessful {
		incrementTotalSuccessfulRequestReceivedPerEndpoint(param.Path)
	} else {
		incrementTotalFailedRequestReceivedPerEndpoint(param.Path)
	}
}

func incrementTotalRequestReceivedPerEndpoint(path string) {
	totalRequests.WithLabelValues(path).Inc()
}

func incrementTotalSuccessfulRequestReceivedPerEndpoint(path string) {
	totalSuccessfulRequests.WithLabelValues(path).Inc()
}

func incrementTotalFailedRequestReceivedPerEndpoint(path string) {
	totalFailedRequests.WithLabelValues(path).Inc()
}

func GetTotalRequestPrometheus() *prometheus.CounterVec {
	return totalRequests
}

func GetTotalSuccessfulRequestPrometheus() *prometheus.CounterVec {
	return totalSuccessfulRequests
}

func GetTotalFailedRequestPrometheus() *prometheus.CounterVec {
	return totalFailedRequests
}
