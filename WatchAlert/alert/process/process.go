package process

import (
	"encoding/json"
	"fmt"
	"time"
	"watchAlert/alert/mute"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
)

func BuildEvent(rule models.AlertRule, labels func() map[string]interface{}) models.AlertCurEvent {
	return models.AlertCurEvent{
		TenantId:             rule.TenantId,
		DatasourceType:       rule.DatasourceType,
		RuleGroupId:          rule.RuleGroupId,
		RuleId:               rule.RuleId,
		RuleName:             rule.RuleName,
		Labels:               labels(),
		EvalInterval:         rule.EvalInterval,
		IsRecovered:          false,
		RepeatNoticeInterval: rule.RepeatNoticeInterval,
		Severity:             rule.Severity,
		EffectiveTime:        rule.EffectiveTime,
		FaultCenterId:        rule.FaultCenterId,
	}
}

func PushEventToFaultCenter(ctx *ctx.Context, event *models.AlertCurEvent) {
	if event == nil {
		return
	}

	ctx.Mux.Lock()
	defer ctx.Mux.Unlock()
	if len(event.TenantId) <= 0 || len(event.Fingerprint) <= 0 {
		return
	}

	cache := ctx.Redis
	cacheEvent, _ := cache.Alert().GetEventFromCache(event.TenantId, event.FaultCenterId, event.Fingerprint)

	// 获取基础信息
	event.FirstTriggerTime = cacheEvent.GetFirstTime()
	event.LastEvalTime = cacheEvent.GetLastEvalTime()
	event.LastSendTime = cacheEvent.GetLastSendTime()
	event.ConfirmState = cacheEvent.GetLastConfirmState()
	event.EventId = cacheEvent.GetEventId()
	event.FaultCenter = cache.FaultCenter().GetFaultCenterInfo(models.BuildFaultCenterInfoCacheKey(event.TenantId, event.FaultCenterId))

	// 获取当前缓存中的状态
	currentStatus := cacheEvent.GetEventStatus()

	// 如果是新的告警事件，设置为 StatePreAlert
	if currentStatus == "" {
		event.Status = models.StatePreAlert
	} else {
		event.Status = currentStatus
	}

	// 检查是否处于静默状态
	isSilenced := IsSilencedEvent(event)

	// 根据不同情况处理状态转换
	switch event.Status {
	case models.StatePreAlert:
		// 如果需要静默
		if isSilenced {
			event.TransitionStatus(models.StateSilenced)
		} else if event.IsArriveForDuration() {
			// 如果达到持续时间，转为告警状态
			event.TransitionStatus(models.StateAlerting)
		}
	case models.StateAlerting:
		// 如果需要静默
		if isSilenced {
			event.TransitionStatus(models.StateSilenced)
		}
	case models.StateSilenced:
		// 如果不再静默，转换回预告警状态
		if !isSilenced {
			event.TransitionStatus(models.StatePreAlert)
		}
	}

	// 检查状态是否发生变化，如果是则触发相应的处理
	if currentStatus != event.Status {
		handleAlertStatusChange(ctx, &cacheEvent, event, currentStatus, event.Status)
	}

	// 最终再次校验 fingerprint 非空，避免 push 时使用空 key
	if event.Fingerprint == "" {
		logc.Errorf(ctx.Ctx, "PushEventToFaultCenter: fingerprint became empty before PushAlertEvent, tenant=%s, rule=%s(%s)", event.TenantId, event.RuleName, event.RuleId)
		return
	}

	// 更新缓存
	cacheKey := models.BuildAlertEventCacheKey(event.TenantId, event.FaultCenterId)
	logc.Infof(ctx.Ctx, "告警事件已存储到Redis，键: %s，事件ID: %s，指纹: %s", cacheKey, event.EventId, event.Fingerprint)
	cache.Alert().PushAlertEvent(event)
}

// GetDutyUsers 获取值班用户列表
func GetDutyUsers(ctx *ctx.Context, noticeData models.AlertNotice) []string {
	if noticeData.DutyId == nil || *noticeData.DutyId == "" {
		return []string{}
	}

	// 获取当前值班用户
	users, ok := ctx.DB.DutyCalendar().GetDutyUserInfo(*noticeData.DutyId, time.Now().Format("2006-01-02"))
	if !ok || len(users) == 0 {
		logc.Errorf(ctx.Ctx, "获取值班用户失败或无值班用户")
		return []string{}
	}

	var usernames []string
	for _, user := range users {
		usernames = append(usernames, user.UserName)
	}

	return usernames
}

// RecordAlertHisEvent 记录告警历史事件
func RecordAlertHisEvent(ctx *ctx.Context, alert models.AlertCurEvent) error {
	hisEvent := models.AlertHisEvent{
		TenantId:         alert.TenantId,
		EventId:          alert.EventId,
		DatasourceType:   alert.DatasourceType,
		Fingerprint:      alert.Fingerprint,
		RuleId:           alert.RuleId,
		RuleName:         alert.RuleName,
		Severity:         alert.Severity,
		Labels:           alert.Labels,
		Annotations:      tools.JsonMarshalToString(alert.Annotations),
		EvalInterval:     alert.EvalInterval,
		FirstTriggerTime: alert.FirstTriggerTime,
		LastEvalTime:     alert.LastEvalTime,
		LastSendTime:     alert.LastSendTime,
		RecoverTime:      time.Now().Unix(),
		FaultCenterId:    alert.FaultCenterId,
		ConfirmState:     alert.ConfirmState,
	}

	// 直接创建记录，数据库会根据event_id自动处理重复
	err := ctx.DB.DB().Create(&hisEvent).Error
	if err != nil {
		logc.Errorf(ctx.Ctx, "记录告警历史失败: %v", err)
		return err
	}

	return nil
}

// handleAlertStatusChange 处理告警状态变化
func handleAlertStatusChange(ctx *ctx.Context, oldEvent *models.AlertCurEvent, newEvent *models.AlertCurEvent, oldStatus, newStatus models.AlertStatus) {
	// 当告警状态变为alerting时，尝试创建工单
	if newStatus == models.StateAlerting {
		// 发布告警转工单事件
		publishAlertTicketEvent(ctx, newEvent, "create_ticket")
		logc.Infof(ctx.Ctx, "发布告警转工单事件，事件ID: %s，规则: %s，租户: %s，故障中心: %s，严重程度: %s",
			newEvent.EventId, newEvent.RuleName, newEvent.TenantId, newEvent.FaultCenterId, newEvent.Severity)
	}

	// 当告警恢复时，记录相关信息
	if newStatus == models.StateRecovered {
		// 发布告警恢复事件
		publishAlertTicketEvent(ctx, newEvent, "alert_recovered")
		logc.Infof(ctx.Ctx, "告警已恢复，事件ID: %s，规则: %s", newEvent.EventId, newEvent.RuleName)
	}
}

// publishAlertTicketEvent 发布告警工单事件到Redis
func publishAlertTicketEvent(ctx *ctx.Context, event *models.AlertCurEvent, eventType string) {
	// 构造事件数据
	eventData := map[string]interface{}{
		"event_type": eventType,
		"event_id":   event.EventId,
		"tenant_id":  event.TenantId,
		"rule_id":    event.RuleId,
		"rule_name":  event.RuleName,
		"severity":   event.Severity,
		"status":     string(event.Status),
		"timestamp":  time.Now().Unix(),
	}

	// 将事件数据序列化为JSON
	eventJSON, _ := json.Marshal(eventData)

	// 发布到Redis频道
	channel := fmt.Sprintf("alert_ticket_events:%s", event.TenantId)
	err := ctx.Redis.Redis().Publish(channel, string(eventJSON)).Err()
	if err != nil {
		logc.Errorf(ctx.Ctx, "发布告警工单事件失败: %v", err)
	}
}

// IsSilencedEvent 静默检查
func IsSilencedEvent(event *models.AlertCurEvent) bool {
	return mute.IsSilence(mute.MuteParams{
		EffectiveTime: event.EffectiveTime,
		Labels:        event.Labels,
		TenantId:      event.TenantId,
		FaultCenterId: event.FaultCenterId,
	})
}
