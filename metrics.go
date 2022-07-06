package temporal

import (
	"time"

	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"
	"go.temporal.io/sdk/client"
)

type MetricsHandler struct {
	vu      modules.VU
	tags    map[string]string
	metrics CustomMetrics
}

type CustomMetrics map[string]*metrics.Metric

func RegisterMetrics(registry *metrics.Registry) CustomMetrics {
	return CustomMetrics{
		"temporal_request":              registry.MustNewMetric("temporal_request", metrics.Counter),
		"temporal_request_latency":      registry.MustNewMetric("temporal_request_latency", metrics.Trend, metrics.Time),
		"temporal_long_request":         registry.MustNewMetric("temporal_long_request", metrics.Counter),
		"temporal_long_request_latency": registry.MustNewMetric("temporal_long_request_latency", metrics.Trend, metrics.Time),
	}
}

func NewMetricsHandler(vu modules.VU, tags map[string]string, customMetrics CustomMetrics) client.MetricsHandler {
	return &MetricsHandler{
		vu:      vu,
		tags:    tags,
		metrics: customMetrics,
	}
}

type metricWrapper struct {
	vu     modules.VU
	tags   *metrics.SampleTags
	metric *metrics.Metric
}

func (w metricWrapper) Inc(v int64) {
	state := w.vu.State()

	metrics.PushIfNotDone(
		w.vu.Context(),
		state.Samples,
		metrics.Sample{
			Time:   time.Now(),
			Metric: w.metric,
			Tags:   w.tags,
			Value:  float64(v),
		},
	)
}

func (w metricWrapper) Update(v float64) {
	state := w.vu.State()

	metrics.PushIfNotDone(
		w.vu.Context(),
		state.Samples,
		metrics.Sample{
			Time:   time.Now(),
			Metric: w.metric,
			Tags:   w.tags,
			Value:  v,
		},
	)
}

func (w metricWrapper) Record(v time.Duration) {
	state := w.vu.State()

	metrics.PushIfNotDone(
		w.vu.Context(),
		state.Samples,
		metrics.Sample{
			Time:   time.Now(),
			Metric: w.metric,
			Tags:   w.tags,
			Value:  metrics.D(v),
		},
	)
}

func (h *MetricsHandler) WithTags(tags map[string]string) client.MetricsHandler {
	mergedTags := h.tags
	for k, v := range tags {
		mergedTags[k] = v
	}

	return NewMetricsHandler(h.vu, mergedTags, h.metrics)
}

func (h *MetricsHandler) Counter(name string) client.MetricsCounter {
	if m, ok := h.metrics[name]; ok {
		return metricWrapper{h.vu, metrics.NewSampleTags(h.tags), m}
	}

	return client.MetricsNopHandler.Counter(name)
}

func (h *MetricsHandler) Gauge(name string) client.MetricsGauge {
	if m, ok := h.metrics[name]; ok {
		return metricWrapper{h.vu, metrics.NewSampleTags(h.tags), m}
	}

	return client.MetricsNopHandler.Gauge(name)
}

func (h *MetricsHandler) Timer(name string) client.MetricsTimer {
	if m, ok := h.metrics[name]; ok {
		return metricWrapper{h.vu, metrics.NewSampleTags(h.tags), m}
	}

	return client.MetricsNopHandler.Timer(name)
}
