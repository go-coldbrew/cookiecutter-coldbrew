package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const namespace = "{{cookiecutter.app_name}}"

var (
	echoTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "echo_total",
		Help:      "Total number of Echo RPC calls by outcome.",
	}, []string{"outcome"})

	echoDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "echo_duration_seconds",
		Help:      "Duration of Echo RPC calls in seconds.",
		Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
	}, []string{"outcome"})

	activeRequests = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "active_requests",
		Help:      "Number of currently active requests.",
	})
)

type appMetrics struct{}

// New returns a new Metrics implementation.
func New() Metrics {
	return &appMetrics{}
}

func (m *appMetrics) IncEchoTotal(outcome string) {
	echoTotal.WithLabelValues(outcome).Inc()
}

func (m *appMetrics) ObserveEchoDuration(outcome string, duration time.Duration) {
	echoDuration.WithLabelValues(outcome).Observe(duration.Seconds())
}

func (m *appMetrics) SetActiveRequests(count int) {
	activeRequests.Set(float64(count))
}
