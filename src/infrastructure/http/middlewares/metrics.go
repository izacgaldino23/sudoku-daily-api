package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/common/expfmt"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
)

func MetricsHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		mfs, err := prometheus.DefaultGatherer.Gather()
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		c.Type("text/plain; version=0.0.4")
		for _, mf := range mfs {
			_, err := expfmt.MetricFamilyToText(c, mf)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func RecordMetrics(method, path string, status int, duration time.Duration) {
	httpRequestsTotal.WithLabelValues(method, path, statusLabel(status)).Inc()
	httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

func statusLabel(status int) string {
	switch {
	case status < 200:
		return "1xx"
	case status < 300:
		return "2xx"
	case status < 400:
		return "3xx"
	case status < 500:
		return "4xx"
	default:
		return "5xx"
	}
}
