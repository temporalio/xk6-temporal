package metrics

import (
	"context"
	"time"

	"go.k6.io/k6/metrics"

	sdkclient "go.temporal.io/sdk/client"
)

// Metrics keys, copied from https://github.com/temporalio/sdk-go/blob/master/internal/common/metrics/constants.go
const (
	TemporalMetricsPrefix = "temporal_"

	WorkflowCompletedCounter     = TemporalMetricsPrefix + "workflow_completed"
	WorkflowCanceledCounter      = TemporalMetricsPrefix + "workflow_canceled"
	WorkflowFailedCounter        = TemporalMetricsPrefix + "workflow_failed"
	WorkflowContinueAsNewCounter = TemporalMetricsPrefix + "workflow_continue_as_new"
	WorkflowEndToEndLatency      = TemporalMetricsPrefix + "workflow_endtoend_latency"

	WorkflowTaskReplayLatency           = TemporalMetricsPrefix + "workflow_task_replay_latency"
	WorkflowTaskQueuePollEmptyCounter   = TemporalMetricsPrefix + "workflow_task_queue_poll_empty"
	WorkflowTaskQueuePollSucceedCounter = TemporalMetricsPrefix + "workflow_task_queue_poll_succeed"
	WorkflowTaskScheduleToStartLatency  = TemporalMetricsPrefix + "workflow_task_schedule_to_start_latency"
	WorkflowTaskExecutionLatency        = TemporalMetricsPrefix + "workflow_task_execution_latency"
	WorkflowTaskExecutionFailureCounter = TemporalMetricsPrefix + "workflow_task_execution_failed"
	WorkflowTaskNoCompletionCounter     = TemporalMetricsPrefix + "workflow_task_no_completion"

	ActivityPollNoTaskCounter             = TemporalMetricsPrefix + "activity_poll_no_task"
	ActivityScheduleToStartLatency        = TemporalMetricsPrefix + "activity_schedule_to_start_latency"
	ActivityExecutionFailedCounter        = TemporalMetricsPrefix + "activity_execution_failed"
	UnregisteredActivityInvocationCounter = TemporalMetricsPrefix + "unregistered_activity_invocation"
	ActivityExecutionLatency              = TemporalMetricsPrefix + "activity_execution_latency"
	ActivitySucceedEndToEndLatency        = TemporalMetricsPrefix + "activity_succeed_endtoend_latency"
	ActivityTaskErrorCounter              = TemporalMetricsPrefix + "activity_task_error"

	LocalActivityTotalCounter           = TemporalMetricsPrefix + "local_activity_total"
	LocalActivityCanceledCounter        = TemporalMetricsPrefix + "local_activity_canceled"
	LocalActivityFailedCounter          = TemporalMetricsPrefix + "local_activity_failed"
	LocalActivityErrorCounter           = TemporalMetricsPrefix + "local_activity_error"
	LocalActivityExecutionLatency       = TemporalMetricsPrefix + "local_activity_execution_latency"
	LocalActivitySucceedEndToEndLatency = TemporalMetricsPrefix + "local_activity_succeed_endtoend_latency"

	CorruptedSignalsCounter = TemporalMetricsPrefix + "corrupted_signals"

	WorkerStartCounter       = TemporalMetricsPrefix + "worker_start"
	WorkerTaskSlotsAvailable = TemporalMetricsPrefix + "worker_task_slots_available"
	PollerStartCounter       = TemporalMetricsPrefix + "poller_start"

	TemporalRequest            = TemporalMetricsPrefix + "request"
	TemporalRequestFailure     = TemporalRequest + "_failure"
	TemporalRequestLatency     = TemporalRequest + "_latency"
	TemporalLongRequest        = TemporalMetricsPrefix + "long_request"
	TemporalLongRequestFailure = TemporalLongRequest + "_failure"
	TemporalLongRequestLatency = TemporalLongRequest + "_latency"

	StickyCacheHit                 = TemporalMetricsPrefix + "sticky_cache_hit"
	StickyCacheMiss                = TemporalMetricsPrefix + "sticky_cache_miss"
	StickyCacheTotalForcedEviction = TemporalMetricsPrefix + "sticky_cache_total_forced_eviction"
	StickyCacheSize                = TemporalMetricsPrefix + "sticky_cache_size"

	WorkflowActiveThreadCount = TemporalMetricsPrefix + "workflow_active_thread_count"
)

func RegisterClientMetrics(registry *metrics.Registry) CustomMetrics {
	return CustomMetrics{
		TemporalRequest:            registry.MustNewMetric(TemporalRequest, metrics.Counter),
		TemporalRequestFailure:     registry.MustNewMetric(TemporalRequestFailure, metrics.Counter),
		TemporalRequestLatency:     registry.MustNewMetric(TemporalRequestLatency, metrics.Trend, metrics.Time),
		TemporalLongRequest:        registry.MustNewMetric(TemporalLongRequest, metrics.Counter),
		TemporalLongRequestFailure: registry.MustNewMetric(TemporalLongRequestFailure, metrics.Counter),
		TemporalLongRequestLatency: registry.MustNewMetric(TemporalLongRequestLatency, metrics.Trend, metrics.Time),

		WorkflowTaskReplayLatency:           registry.MustNewMetric(WorkflowTaskReplayLatency, metrics.Trend, metrics.Time),
		WorkflowTaskQueuePollEmptyCounter:   registry.MustNewMetric(WorkflowTaskQueuePollEmptyCounter, metrics.Counter),
		WorkflowTaskQueuePollSucceedCounter: registry.MustNewMetric(WorkflowTaskQueuePollSucceedCounter, metrics.Counter),
		WorkflowTaskScheduleToStartLatency:  registry.MustNewMetric(WorkflowTaskScheduleToStartLatency, metrics.Trend, metrics.Time),
		WorkflowTaskExecutionLatency:        registry.MustNewMetric(WorkflowTaskExecutionLatency, metrics.Trend, metrics.Time),
		WorkflowTaskExecutionFailureCounter: registry.MustNewMetric(WorkflowTaskExecutionFailureCounter, metrics.Counter),
		WorkflowTaskNoCompletionCounter:     registry.MustNewMetric(WorkflowTaskNoCompletionCounter, metrics.Counter),

		ActivityPollNoTaskCounter:             registry.MustNewMetric(ActivityPollNoTaskCounter, metrics.Counter),
		ActivityScheduleToStartLatency:        registry.MustNewMetric(ActivityScheduleToStartLatency, metrics.Trend, metrics.Time),
		ActivityExecutionFailedCounter:        registry.MustNewMetric(ActivityExecutionFailedCounter, metrics.Counter),
		UnregisteredActivityInvocationCounter: registry.MustNewMetric(UnregisteredActivityInvocationCounter, metrics.Counter),
		ActivityExecutionLatency:              registry.MustNewMetric(ActivityExecutionLatency, metrics.Trend, metrics.Time),
		ActivitySucceedEndToEndLatency:        registry.MustNewMetric(ActivitySucceedEndToEndLatency, metrics.Trend, metrics.Time),
		ActivityTaskErrorCounter:              registry.MustNewMetric(ActivityTaskErrorCounter, metrics.Counter),

		WorkflowCompletedCounter:     registry.MustNewMetric(WorkflowCompletedCounter, metrics.Counter),
		WorkflowCanceledCounter:      registry.MustNewMetric(WorkflowCanceledCounter, metrics.Counter),
		WorkflowFailedCounter:        registry.MustNewMetric(WorkflowFailedCounter, metrics.Counter),
		WorkflowContinueAsNewCounter: registry.MustNewMetric(WorkflowContinueAsNewCounter, metrics.Counter),
		WorkflowEndToEndLatency:      registry.MustNewMetric(WorkflowEndToEndLatency, metrics.Trend, metrics.Time),

		WorkerTaskSlotsAvailable: registry.MustNewMetric(WorkerTaskSlotsAvailable, metrics.Gauge),
	}
}

type ClientMetricsHandler struct {
	ctx     context.Context
	samples chan<- metrics.SampleContainer
	tags    map[string]string
	metrics CustomMetrics
}

type CustomMetrics map[string]*metrics.Metric

func NewClientMetricsHandler(ctx context.Context, samples chan<- metrics.SampleContainer, tags map[string]string, customMetrics CustomMetrics) sdkclient.MetricsHandler {
	return &ClientMetricsHandler{
		ctx:     ctx,
		samples: samples,
		tags:    tags,
		metrics: customMetrics,
	}
}

type clientMetricWrapper struct {
	ctx     context.Context
	samples chan<- metrics.SampleContainer
	tags    *metrics.SampleTags
	metric  *metrics.Metric
}

func (w clientMetricWrapper) Inc(v int64) {
	tags := w.tags.CloneTags()

	metrics.PushIfNotDone(
		w.ctx,
		w.samples,
		metrics.Sample{
			Time:   time.Now(),
			Metric: w.metric,
			Tags:   metrics.IntoSampleTags(&tags),
			Value:  float64(v),
		},
	)
}

func (w clientMetricWrapper) Update(v float64) {
	tags := w.tags.CloneTags()

	metrics.PushIfNotDone(
		w.ctx,
		w.samples,
		metrics.Sample{
			Time:   time.Now(),
			Metric: w.metric,
			Tags:   metrics.IntoSampleTags(&tags),
			Value:  v,
		},
	)
}

func (w clientMetricWrapper) Record(v time.Duration) {
	tags := w.tags.CloneTags()

	metrics.PushIfNotDone(
		w.ctx,
		w.samples,
		metrics.Sample{
			Time:   time.Now(),
			Metric: w.metric,
			Tags:   metrics.IntoSampleTags(&tags),
			Value:  metrics.D(v),
		},
	)
}

func (h *ClientMetricsHandler) WithTags(tags map[string]string) sdkclient.MetricsHandler {
	mergedTags := make(map[string]string, len(h.tags)+len(tags))
	for k, v := range h.tags {
		mergedTags[k] = v
	}
	for k, v := range tags {
		mergedTags[k] = v
	}

	return NewClientMetricsHandler(h.ctx, h.samples, mergedTags, h.metrics)
}

func (h *ClientMetricsHandler) Counter(name string) sdkclient.MetricsCounter {
	if m, ok := h.metrics[name]; ok {
		return clientMetricWrapper{h.ctx, h.samples, metrics.NewSampleTags(h.tags), m}
	}

	return sdkclient.MetricsNopHandler.Counter(name)
}

func (h *ClientMetricsHandler) Gauge(name string) sdkclient.MetricsGauge {
	if m, ok := h.metrics[name]; ok {
		return clientMetricWrapper{h.ctx, h.samples, metrics.NewSampleTags(h.tags), m}
	}

	return sdkclient.MetricsNopHandler.Gauge(name)
}

func (h *ClientMetricsHandler) Timer(name string) sdkclient.MetricsTimer {
	if m, ok := h.metrics[name]; ok {
		return clientMetricWrapper{h.ctx, h.samples, metrics.NewSampleTags(h.tags), m}
	}

	return sdkclient.MetricsNopHandler.Timer(name)
}
