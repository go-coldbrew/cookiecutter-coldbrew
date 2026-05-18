package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const namespace = "{{ cookiecutter.app_name | replace('-', '_') | replace('.', '_') }}"

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

	streamEchoTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "stream_echo_total",
		Help:      "Total number of StreamEcho RPC calls by outcome (success, error, canceled).",
	}, []string{"outcome"})

	streamEchoDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "stream_echo_duration_seconds",
		Help:      "Total duration of StreamEcho RPC calls in seconds.",
		Buckets:   []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30, 60},
	}, []string{"outcome"})

	streamEchoTTFT = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "stream_echo_ttft_seconds",
		Help:      "Time to first token (TTFT) for StreamEcho RPC calls in seconds. Recorded only when at least one token is emitted.",
		Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
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

func (m *appMetrics) IncStreamEchoTotal(outcome string) {
	streamEchoTotal.WithLabelValues(outcome).Inc()
}

func (m *appMetrics) ObserveStreamEchoDuration(outcome string, duration time.Duration) {
	streamEchoDuration.WithLabelValues(outcome).Observe(duration.Seconds())
}

func (m *appMetrics) ObserveStreamEchoTTFT(duration time.Duration) {
	streamEchoTTFT.Observe(duration.Seconds())
}
