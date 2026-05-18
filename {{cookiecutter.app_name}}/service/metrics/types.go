package metrics

import "time"

// Metrics defines the application metrics interface.
// All methods are safe for concurrent use.
//
// Add new methods here as your service grows. The interface
// enables mocking in tests via mockery.
type Metrics interface {
	// Echo RPC metrics
	IncEchoTotal(outcome string)
	ObserveEchoDuration(outcome string, duration time.Duration)

	// StreamEcho RPC metrics. Time-to-first-token (TTFT) is the standard
	// AI/LLM streaming metric — the gap between request arrival and the first
	// emitted frame. Tracking it separately from total duration surfaces
	// upstream-latency vs. generation-throughput issues independently.
	IncStreamEchoTotal(outcome string)
	ObserveStreamEchoDuration(outcome string, duration time.Duration)
	ObserveStreamEchoTTFT(duration time.Duration)
}
