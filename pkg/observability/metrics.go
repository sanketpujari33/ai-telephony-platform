package observability

import (
	"github.com/ai-telephony-platform/ai-telephony-platform/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CallLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "ai_telephony_call_latency_seconds",
			Help: "Latency of calls processed by the platform",
			Buckets: []float64{0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0},
		},
		[]string{"provider", "operation", "tenant"},
	)

	ConcurrentCalls = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ai_telephony_concurrent_calls",
			Help: "Number of currently active calls",
		},
		[]string{"service", "tenant"},
	)

	ErrorRates = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ai_telephony_errors_total",
			Help: "Total number of errors encountered",
		},
		[]string{"service", "error_type", "tenant"},
	)
)

// InitMetrics initializes metrics configurations
func InitMetrics(cfg config.ObservabilityMetricsConfig) {
	// Prometheus metrics are automatically registered via promauto when they are defined above.
	// We could use the cfg here to selectively register/unregister if desired.
}
