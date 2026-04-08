package metrics

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func gatherMetric(t *testing.T, name string) *dto.MetricFamily {
	t.Helper()
	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatalf("failed to gather metrics: %v", err)
	}
	for _, f := range families {
		if f.GetName() == name {
			return f
		}
	}
	return nil
}

func TestIncEchoTotal(t *testing.T) {
	m := New()
	m.IncEchoTotal(OutcomeSuccess)
	m.IncEchoTotal(OutcomeError)

	mf := gatherMetric(t, namespace+"_echo_total")
	if mf == nil {
		t.Fatal("metric not found")
	}
	if len(mf.GetMetric()) < 2 {
		t.Fatalf("expected at least 2 label pairs, got %d", len(mf.GetMetric()))
	}
}

func TestObserveEchoDuration(t *testing.T) {
	m := New()
	m.ObserveEchoDuration(OutcomeSuccess, 50*time.Millisecond)

	mf := gatherMetric(t, namespace+"_echo_duration_seconds")
	if mf == nil {
		t.Fatal("metric not found")
	}
	h := mf.GetMetric()[0].GetHistogram()
	if h.GetSampleCount() == 0 {
		t.Fatal("expected at least one observation")
	}
}
