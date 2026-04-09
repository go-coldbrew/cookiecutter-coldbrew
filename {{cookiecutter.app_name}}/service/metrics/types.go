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
}
