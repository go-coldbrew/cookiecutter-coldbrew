package metrics

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func gatherMetric(name string) *dto.MetricFamily {
	families, _ := prometheus.DefaultGatherer.Gather()
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

	mf := gatherMetric(namespace + "_echo_total")
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

	mf := gatherMetric(namespace + "_echo_duration_seconds")
	if mf == nil {
		t.Fatal("metric not found")
	}
	h := mf.GetMetric()[0].GetHistogram()
	if h.GetSampleCount() == 0 {
		t.Fatal("expected at least one observation")
	}
}
