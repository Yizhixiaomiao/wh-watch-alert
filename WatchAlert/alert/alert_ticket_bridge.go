package alert

import (
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
)

// AlertTicketBridge 告警转工单桥接服务
type AlertTicketBridge struct {
	ctx *ctx.Context
}

// NewAlertTicketBridge 创建告警转工单桥接服务
func NewAlertTicketBridge(ctx *ctx.Context) *AlertTicketBridge {
	return &AlertTicketBridge{
		ctx: ctx,
	}
}

// HandleAlertStatusChange 处理告警状态变更
func (b *AlertTicketBridge) HandleAlertStatusChange(alert *models.AlertCurEvent, oldStatus, newStatus models.AlertStatus) error {
	// 当告警状态变为alerting时，尝试创建工单
	if newStatus == models.StateAlerting {
		// 这里需要调用工单服务创建工单
		// 由于循环依赖问题，暂时留空，后续在处理层实现
	}

	// 当告警恢复时，检查是否需要自动关闭工单
	if newStatus == models.StateRecovered {
		return b.handleAlertRecovered(alert)
	}

	return nil
}

// handleAlertRecovered 处理告警恢复
func (b *AlertTicketBridge) handleAlertRecovered(alert *models.AlertCurEvent) error {
	// 告警恢复处理逻辑
	// 暂时只记录日志
	return nil
}
