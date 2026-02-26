package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// 业务指标
	TicketsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "watchalert_tickets_total",
			Help: "Total number of tickets by status",
		},
		[]string{"status", "tenant_id"},
	)

	TicketsCreatedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "watchalert_tickets_created_total",
			Help: "Total number of tickets created",
		},
		[]string{"tenant_id"},
	)

	TicketsResolvedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "watchalert_tickets_resolved_total",
			Help: "Total number of tickets resolved",
		},
		[]string{"tenant_id"},
	)

	TicketsOverdueGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "watchalert_tickets_overdue",
			Help: "Number of overdue tickets",
		},
		[]string{"tenant_id"},
	)

	TicketResolutionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "watchalert_ticket_resolution_duration_seconds",
			Help:    "Ticket resolution duration in seconds",
			Buckets: []float64{60, 300, 900, 3600, 7200, 86400}, // 1m, 5m, 15m, 1h, 2h, 24h
		},
		[]string{"tenant_id"},
	)

	AlertsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "watchalert_alerts_total",
			Help: "Total number of alerts by status",
		},
		[]string{"status", "tenant_id"},
	)

	AlertsActiveGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "watchalert_alerts_active",
			Help: "Number of active alerts",
		},
		[]string{"tenant_id"},
	)

	SLAComplianceRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "watchalert_sla_compliance_rate",
			Help: "SLA compliance rate (0-1)",
		},
		[]string{"tenant_id", "priority"},
	)

	// 性能指标
	HTTPRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "watchalert_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "watchalert_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "watchalert_db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	RedisOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "watchalert_redis_operation_duration_seconds",
			Help:    "Redis operation duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

// 工单相关指标
func IncTicketsCreated(tenantId string) {
	TicketsCreatedTotal.WithLabelValues(tenantId).Inc()
	TicketsTotal.WithLabelValues("created", tenantId).Inc()
}

func IncTicketsResolved(tenantId string) {
	TicketsResolvedTotal.WithLabelValues(tenantId).Inc()
	TicketsTotal.WithLabelValues("resolved", tenantId).Inc()
}

func SetTicketsOverdue(tenantId string, count int) {
	TicketsOverdueGauge.WithLabelValues(tenantId).Set(float64(count))
}

func RecordTicketResolution(tenantId string, duration float64) {
	TicketResolutionDuration.WithLabelValues(tenantId).Observe(duration)
}

// 告警相关指标
func IncAlertsFiring(tenantId string) {
	AlertsTotal.WithLabelValues("firing", tenantId).Inc()
}

func IncAlertsResolved(tenantId string) {
	AlertsTotal.WithLabelValues("resolved", tenantId).Inc()
}

func SetAlertsActive(tenantId string, count int) {
	AlertsActiveGauge.WithLabelValues(tenantId).Set(float64(count))
}

// SLA相关指标
func SetSLAComplianceRate(tenantId, priority string, rate float64) {
	SLAComplianceRate.WithLabelValues(tenantId, priority).Set(rate)
}

// HTTP请求相关指标
func RecordHTTPRequest(method, path, status string, duration float64) {
	HTTPRequestTotal.WithLabelValues(method, path, status).Inc()
	HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
}

// 数据库相关指标
func RecordDBQuery(operation string, duration float64) {
	DBQueryDuration.WithLabelValues(operation).Observe(duration)
}

// Redis相关指标
func RecordRedisOperation(operation string, duration float64) {
	RedisOperationDuration.WithLabelValues(operation).Observe(duration)
}
