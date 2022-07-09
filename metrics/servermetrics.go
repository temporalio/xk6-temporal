package metrics

import (
	"context"
	"time"

	"go.k6.io/k6/metrics"

	"go.temporal.io/server/common/log"
	servermetrics "go.temporal.io/server/common/metrics"
)

func RegisterServerMetrics(registry *metrics.Registry) CustomMetrics {
	return CustomMetrics{}
}

type ServerMetricsHandler struct {
	ctx     context.Context
	samples chan<- metrics.SampleContainer
	tags    map[string]string
	metrics CustomMetrics
}

func NewServerMetricsHandler(ctx context.Context, samples chan<- metrics.SampleContainer, tags map[string]string, customMetrics CustomMetrics) servermetrics.MetricsHandler {
	return &ServerMetricsHandler{
		ctx:     ctx,
		samples: samples,
		tags:    tags,
		metrics: customMetrics,
	}
}

type serverCounterMetricWrapper struct {
	ctx     context.Context
	samples chan<- metrics.SampleContainer
	tags    *metrics.SampleTags
	metric  *metrics.Metric
}

type serverGaugeMetricWrapper struct {
	ctx     context.Context
	samples chan<- metrics.SampleContainer
	tags    *metrics.SampleTags
	metric  *metrics.Metric
}

type serverTimerMetricWrapper struct {
	ctx     context.Context
	samples chan<- metrics.SampleContainer
	tags    *metrics.SampleTags
	metric  *metrics.Metric
}

type serverHistogramMetricWrapper struct {
	ctx     context.Context
	samples chan<- metrics.SampleContainer
	tags    *metrics.SampleTags
	metric  *metrics.Metric
}

func (w serverCounterMetricWrapper) Record(v int64, tags ...servermetrics.Tag) {
	mergedTags := w.tags.CloneTags()
	for _, t := range tags {
		mergedTags[t.Key()] = t.Value()
	}

	metrics.PushIfNotDone(
		w.ctx,
		w.samples,
		metrics.Sample{
			Time:   time.Now(),
			Metric: w.metric,
			Tags:   metrics.IntoSampleTags(&mergedTags),
			Value:  float64(v),
		},
	)
}

func (w serverGaugeMetricWrapper) Record(v float64, tags ...servermetrics.Tag) {
	mergedTags := w.tags.CloneTags()
	for _, t := range tags {
		mergedTags[t.Key()] = t.Value()
	}

	metrics.PushIfNotDone(
		w.ctx,
		w.samples,
		metrics.Sample{
			Time:   time.Now(),
			Metric: w.metric,
			Tags:   metrics.IntoSampleTags(&mergedTags),
			Value:  v,
		},
	)
}

func (w serverTimerMetricWrapper) Record(v time.Duration, tags ...servermetrics.Tag) {
	mergedTags := w.tags.CloneTags()
	for _, t := range tags {
		mergedTags[t.Key()] = t.Value()
	}

	metrics.PushIfNotDone(
		w.ctx,
		w.samples,
		metrics.Sample{
			Time:   time.Now(),
			Metric: w.metric,
			Tags:   metrics.IntoSampleTags(&mergedTags),
			Value:  metrics.D(v),
		},
	)
}

func (w serverHistogramMetricWrapper) Record(v int64, tags ...servermetrics.Tag) {
	mergedTags := w.tags.CloneTags()
	for _, t := range tags {
		mergedTags[t.Key()] = t.Value()
	}

	metrics.PushIfNotDone(
		w.ctx,
		w.samples,
		metrics.Sample{
			Time:   time.Now(),
			Metric: w.metric,
			Tags:   metrics.IntoSampleTags(&mergedTags),
			Value:  float64(v),
		},
	)
}

func (h *ServerMetricsHandler) WithTags(tags ...servermetrics.Tag) servermetrics.MetricsHandler {
	mergedTags := make(map[string]string, len(h.tags)+len(tags))
	for k, v := range h.tags {
		mergedTags[k] = v
	}
	for _, t := range tags {
		mergedTags[t.Key()] = t.Value()
	}

	return NewServerMetricsHandler(h.ctx, h.samples, mergedTags, h.metrics)
}

func (h *ServerMetricsHandler) Counter(name string) servermetrics.CounterMetric {
	if m, ok := h.metrics[name]; ok {
		return serverCounterMetricWrapper{h.ctx, h.samples, metrics.NewSampleTags(h.tags), m}
	}

	return servermetrics.NoopMetricsHandler.Counter(name)
}

func (h *ServerMetricsHandler) Gauge(name string) servermetrics.GaugeMetric {
	if m, ok := h.metrics[name]; ok {
		return serverGaugeMetricWrapper{h.ctx, h.samples, metrics.NewSampleTags(h.tags), m}
	}

	return servermetrics.NoopMetricsHandler.Gauge(name)
}

func (h *ServerMetricsHandler) Timer(name string) servermetrics.TimerMetric {
	if m, ok := h.metrics[name]; ok {
		return serverTimerMetricWrapper{h.ctx, h.samples, metrics.NewSampleTags(h.tags), m}
	}

	return servermetrics.NoopMetricsHandler.Timer(name)
}

func (h *ServerMetricsHandler) Histogram(name string, unit servermetrics.MetricUnit) servermetrics.HistogramMetric {
	if m, ok := h.metrics[name]; ok {
		return serverHistogramMetricWrapper{h.ctx, h.samples, metrics.NewSampleTags(h.tags), m}
	}

	return servermetrics.NoopMetricsHandler.Histogram(name, unit)
}

func (h *ServerMetricsHandler) Stop(log.Logger) {}
