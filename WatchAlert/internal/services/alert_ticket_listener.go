package services

import (
	"encoding/json"
	"fmt"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/repo"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logc"
)

// AlertTicketEventListener 告警工单事件监听器
type AlertTicketEventListener struct {
	ctx    *ctx.Context
	stopCh chan struct{}
}

// NewAlertTicketEventListener 创建告警工单事件监听器
func NewAlertTicketEventListener(ctx *ctx.Context) *AlertTicketEventListener {
	return &AlertTicketEventListener{
		ctx:    ctx,
		stopCh: make(chan struct{}),
	}
}

// Start 启动事件监听器
func (l *AlertTicketEventListener) Start() {
	go l.listenForEvents()
	logc.Info(l.ctx.Ctx, "告警工单事件监听器已启动")
}

// Stop 停止事件监听器
func (l *AlertTicketEventListener) Stop() {
	close(l.stopCh)
	logc.Info(l.ctx.Ctx, "告警工单事件监听器已停止")
}

// listenForEvents 监听Redis事件
func (l *AlertTicketEventListener) listenForEvents() {
	// 订阅告警工单事件频道（使用通配符模式）
	pubsub := l.ctx.Redis.Redis().PSubscribe("alert_ticket_events:*")
	defer pubsub.Close()

	for {
		select {
		case <-l.stopCh:
			return
		default:
			msg, err := pubsub.ReceiveMessage()
			if err != nil {
				logc.Errorf(l.ctx.Ctx, "接收Redis消息失败: %v", err)
				continue
			}

			// 处理事件
			go l.handleEvent(msg)
		}
	}
}

// handleEvent 处理告警工单事件
func (l *AlertTicketEventListener) handleEvent(msg *redis.Message) {
	// 解析事件数据
	var eventData map[string]interface{}
	if err := json.Unmarshal([]byte(msg.Payload), &eventData); err != nil {
		logc.Errorf(l.ctx.Ctx, "解析事件数据失败: %v", err)
		return
	}

	eventType, _ := eventData["event_type"].(string)
	tenantId, _ := eventData["tenant_id"].(string)
	eventId, _ := eventData["event_id"].(string)

	switch eventType {
	case "create_ticket":
		l.createTicketFromAlert(tenantId, eventId)
	case "alert_recovered":
		l.handleAlertRecovered(tenantId, eventId)
	default:
		logc.Infof(l.ctx.Ctx, "未知的事件类型: %s", eventType)
	}
}

// createTicketFromAlert 从告警创建工单
func (l *AlertTicketEventListener) createTicketFromAlert(tenantId, eventId string) {
	// 从Redis缓存中获取告警事件
	// 这里需要遍历所有故障中心来查找事件
	// 为了简化，我们假设事件存在于某个故障中心中

	// 实际实现中，我们需要找到该事件所属的故障中心
	// 这里先简单处理，后续可以优化

	logc.Infof(l.ctx.Ctx, "准备为告警事件创建工单: %s，租户: %s", eventId, tenantId)

	// 检查是否已经存在工单
	existingTicket, err := l.getTicketByEventId(tenantId, eventId)
	if err == nil && existingTicket != nil {
		logc.Infof(l.ctx.Ctx, "工单已存在，跳过创建: %s", existingTicket.TicketId)
		return
	}

	// 获取告警事件（这里需要实际实现查找逻辑）
	alertEvent, err := l.getAlertEvent(tenantId, eventId)
	if err != nil {
		logc.Errorf(l.ctx.Ctx, "获取告警事件失败: %v", err)
		return
	}

	// 创建工单
	if err := l.createAlertTicket(tenantId, alertEvent); err != nil {
		logc.Errorf(l.ctx.Ctx, "创建告警工单失败: %v", err)
		return
	}

	logc.Infof(l.ctx.Ctx, "成功为告警事件创建工单: %s", eventId)
}

// handleAlertRecovered 处理告警恢复事件
func (l *AlertTicketEventListener) handleAlertRecovered(tenantId, eventId string) {
	logc.Infof(l.ctx.Ctx, "处理告警恢复事件: %s", eventId)

	// 获取关联的工单
	ticket, err := l.getTicketByEventId(tenantId, eventId)
	if err != nil {
		// 没有找到工单，跳过
		logc.Infof(l.ctx.Ctx, "未找到关联工单，跳过告警恢复处理: %s", eventId)
		return
	}

	// 检查工单当前状态
	if ticket.Type != models.TicketTypeAlert {
		logc.Infof(l.ctx.Ctx, "工单类型不是告警工单，跳过: %s", ticket.TicketId)
		return
	}

	// 更新工单状态
	now := time.Now().Unix()

	// 如果工单已经是关闭或解决状态，只更新告警状态
	if ticket.Status == models.TicketStatusClosed || ticket.Status == models.TicketStatusResolved {
		ticket.AlarmActive = false
		ticket.LastSyncTime = now
		logc.Infof(l.ctx.Ctx, "告警已恢复，更新工单告警状态: %s", ticket.TicketId)
	} else {
		// 自动将工单状态更新为待验证（已解决）
		ticket.Status = models.TicketStatusResolved
		ticket.AlarmActive = false
		ticket.LastSyncTime = now
		ticket.ResolvedAt = now
		logc.Infof(l.ctx.Ctx, "告警已恢复，自动更新工单状态为已解决: %s", ticket.TicketId)
	}

	// 保存更新
	err = l.ctx.DB.Ticket().Update(*ticket)
	if err != nil {
		logc.Errorf(l.ctx.Ctx, "更新工单失败: %v", err)
		return
	}

	// 创建工作日志
	action := "auto_resolve"
	content := fmt.Sprintf("告警已自动恢复，系统自动将工单状态更新为已解决")
	if ticket.Status == models.TicketStatusClosed {
		action = "sync"
		content = fmt.Sprintf("告警已恢复，同步更新工单告警状态")
	}
	l.createWorkLog(ticket.TicketId, "system", action, content, "", "")

	logc.Infof(l.ctx.Ctx, "成功处理告警恢复事件，工单: %s", ticket.TicketId)
}

// getTicketByEventId 根据事件ID获取工单
func (l *AlertTicketEventListener) getTicketByEventId(tenantId, eventId string) (*models.Ticket, error) {
	// 构造查询条件
	query := repo.TicketQuery{
		EventId: eventId,
		Page:    1,
		Size:    1,
	}

	tickets, _, err := l.ctx.DB.Ticket().List(tenantId, query)
	if err != nil {
		return nil, err
	}

	if len(tickets) == 0 {
		return nil, fmt.Errorf("工单不存在")
	}

	return &tickets[0], nil
}

// getAlertEvent 获取告警事件
func (l *AlertTicketEventListener) getAlertEvent(tenantId, eventId string) (*models.AlertCurEvent, error) {
	// 获取所有故障中心
	faultCenters, err := l.ctx.DB.FaultCenter().List(tenantId, "")
	if err != nil {
		return nil, fmt.Errorf("获取故障中心列表失败: %v", err)
	}

	logc.Infof(l.ctx.Ctx, "获取到 %d 个故障中心，租户: %s", len(faultCenters), tenantId)
	for _, fc := range faultCenters {
		logc.Infof(l.ctx.Ctx, "  - 故障中心: %s (ID: %s)", fc.Name, fc.ID)
	}

	// 遍历所有故障中心，查找匹配的事件
	for _, fc := range faultCenters {
		cacheKey := models.BuildAlertEventCacheKey(tenantId, fc.ID)
		logc.Infof(l.ctx.Ctx, "查找故障中心 %s (ID: %s)，Redis键: %s", fc.Name, fc.ID, cacheKey)
		events, err := l.ctx.Redis.Alert().GetAllEvents(cacheKey)
		if err != nil {
			logc.Infof(l.ctx.Ctx, "  - 查找失败: %v", err)
			continue
		}
		logc.Infof(l.ctx.Ctx, "  - 找到 %d 个事件", len(events))

		// 查找匹配的事件ID
		for _, event := range events {
			if event.EventId == eventId {
				logc.Infof(l.ctx.Ctx, "找到告警事件: %s，故障中心: %s (ID: %s)，规则: %s", eventId, fc.Name, fc.ID, event.RuleName)
				return event, nil
			}
		}
	}

	logc.Infof(l.ctx.Ctx, "未在故障中心中找到告警事件: %s，尝试从历史记录查找", eventId)

	// 如果在 Redis 中找不到，尝试从历史事件中查找
	historyQuery := types.RequestAlertHisEventQuery{
		TenantId: tenantId,
		Page:     models.Page{Index: 1, Size: 100},
	}

	historyResult, err := l.ctx.DB.Event().GetHistoryEvent(historyQuery)
	if err != nil {
		return nil, fmt.Errorf("获取历史事件失败: %v", err)
	}

	for _, historyEvent := range historyResult.List {
		if historyEvent.EventId == eventId {
			// 转换历史事件为当前事件结构
			event := &models.AlertCurEvent{
				TenantId:         historyEvent.TenantId,
				EventId:          historyEvent.EventId,
				RuleId:           historyEvent.RuleId,
				RuleName:         historyEvent.RuleName,
				DatasourceType:   historyEvent.DatasourceType,
				DatasourceId:     historyEvent.DatasourceId,
				Fingerprint:      historyEvent.Fingerprint,
				Severity:         historyEvent.Severity,
				Status:           models.StateRecovered,
				Labels:           historyEvent.Labels,
				EvalInterval:     historyEvent.EvalInterval,
				FaultCenterId:    historyEvent.FaultCenterId,
				FirstTriggerTime: historyEvent.FirstTriggerTime,
				LastEvalTime:     historyEvent.LastEvalTime,
				IsRecovered:      true,
			}
			logc.Infof(l.ctx.Ctx, "从历史记录中找到告警事件: %s", eventId)
			return event, nil
		}
	}

	return nil, fmt.Errorf("未找到告警事件: %s", eventId)
}

// createAlertTicket 创建告警工单
func (l *AlertTicketEventListener) createAlertTicket(tenantId string, alert *models.AlertCurEvent) error {
	// 生成工单ID和编号
	ticketId := "tk-" + tools.RandId()
	ticketNo := generateAlertTicketNo()

	// 根据告警严重程度确定工单优先级
	priority := l.mapSeverityToPriority(alert.Severity)
	severity := l.mapSeverityToTicketSeverity(alert.Severity)

	// 获取SLA策略
	var responseSLA, resolutionSLA int64
	var dueTime int64
	slaPolicy, err := l.ctx.DB.Ticket().GetSLAPolicyByPriority(tenantId, priority)
	if err == nil {
		responseSLA = slaPolicy.ResponseTime
		resolutionSLA = slaPolicy.ResolutionTime
		
		// 使用工作日计算SLA截止时间
		workingHoursConfig, err := tools.ParseWorkingHoursConfig(slaPolicy.WorkingHours)
		if err == nil && workingHoursConfig != nil {
			// 设置节假日
			workingHoursConfig.Holidays = slaPolicy.Holidays
			
			// 计算考虑工作日的截止时间
			dueTimeObj := tools.CalculateSLADueTime(time.Now(), resolutionSLA, workingHoursConfig)
			dueTime = dueTimeObj.Unix()
		} else {
			// 如果配置解析失败，使用简单计算
			dueTime = time.Now().Unix() + resolutionSLA
		}
	}

	// 构建工单标题和描述
	title, description := l.buildTicketContent(alert)

	// 构建标签
	labels := map[string]string{
		"alert_event_id":   alert.EventId,
		"alert_rule_id":    alert.RuleId,
		"alert_severity":   alert.Severity,
		"alert_datasource": alert.DatasourceType,
		"created_from":     "alert",
	}

	// 添加告警标签到工单标签
	for k, v := range alert.Labels {
		if strVal, ok := v.(string); ok {
			labels[k] = strVal
		}
	}

	// 构建自定义字段
	customFields := map[string]interface{}{
		"alert_first_trigger_time": alert.FirstTriggerTime,
		"alert_fingerprint":        alert.Fingerprint,
		"alert_for_duration":       alert.ForDuration,
		"alert_eval_interval":      alert.EvalInterval,
		"auto_created":             true,
	}

	ticket := models.Ticket{
		TenantId:       tenantId,
		TicketId:       ticketId,
		TicketNo:       ticketNo,
		Title:          title,
		Description:    description,
		Type:           models.TicketTypeAlert,
		Priority:       priority,
		Severity:       severity,
		Status:         models.TicketStatusPending,
		Source:         models.TicketSourceAuto,
		EventId:        alert.EventId,
		FaultCenterId:  alert.FaultCenterId,
		RuleId:         alert.RuleId,
		DatasourceType: alert.DatasourceType,
		CreatedBy:      "system",
		AssignedTo:     "",
		AssignedGroup:  "",
		Followers:      []string{},
		Labels:         labels,
		Tags:           []string{"告警工单", alert.DatasourceType, alert.Severity},
		CustomFields:   customFields,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
		ResponseSLA:    responseSLA,
		ResolutionSLA:  resolutionSLA,
		DueTime:        dueTime,
		IsOverdue:      false,
		AlarmActive:    true,
		LastSyncTime:   time.Now().Unix(),
	}

	err = l.ctx.DB.Ticket().Create(ticket)
	if err != nil {
		return fmt.Errorf("创建告警工单失败: %v", err)
	}

	// 创建工作日志
	l.createWorkLog(ticketId, "system", "create", "告警自动创建工单", "", "")

	return nil
}

// mapSeverityToPriority 将告警严重程度映射为工单优先级
func (l *AlertTicketEventListener) mapSeverityToPriority(severity string) models.TicketPriority {
	switch severity {
	case "critical":
		return models.TicketPriorityP0
	case "warning":
		return models.TicketPriorityP1
	case "info":
		return models.TicketPriorityP2
	default:
		return models.TicketPriorityP2
	}
}

// mapSeverityToTicketSeverity 将告警严重程度映射为工单严重程度
func (l *AlertTicketEventListener) mapSeverityToTicketSeverity(severity string) models.TicketSeverity {
	switch severity {
	case "critical":
		return models.TicketSeverityCritical
	case "warning":
		return models.TicketSeverityHigh
	case "info":
		return models.TicketSeverityMedium
	default:
		return models.TicketSeverityMedium
	}
}

// buildTicketContent 构建工单标题和描述
func (l *AlertTicketEventListener) buildTicketContent(alert *models.AlertCurEvent) (title, description string) {
	// 构建标题
	title = fmt.Sprintf("[告警] %s - %s", alert.RuleName, alert.Severity)

	// 构建描述
	description = fmt.Sprintf(`
## 告警详情

**规则名称**: %s
**严重程度**: %s
**数据源**: %s
**首次触发时间**: %s
**持续时间**: %s

## 告警标签
`, alert.RuleName, alert.Severity, alert.DatasourceType,
		time.Unix(alert.FirstTriggerTime, 0).Format("2006-01-02 15:04:05"),
		fmt.Sprintf("%d秒", alert.ForDuration))

	// 添加标签信息
	if len(alert.Labels) > 0 {
		for k, v := range alert.Labels {
			description += fmt.Sprintf("- **%s**: %v\n", k, v)
		}
	}

	description += fmt.Sprintf(`
## 处理建议

1. 检查相关服务状态
2. 查看详细日志信息
3. 确认影响范围
4. 制定解决方案

---
*此工单由告警系统自动创建*
`)

	return title, description
}

// createWorkLog 创建工作日志
func (l *AlertTicketEventListener) createWorkLog(ticketId, userId, action, content, oldValue, newValue string) {
	log := models.TicketWorkLog{
		Id:        "log-" + tools.RandId(),
		TicketId:  ticketId,
		UserId:    userId,
		Action:    action,
		Content:   content,
		OldValue:  oldValue,
		NewValue:  newValue,
		CreatedAt: time.Now().Unix(),
	}
	l.ctx.DB.Ticket().CreateWorkLog(log)
}

// generateAlertTicketNo 生成告警工单编号
func generateAlertTicketNo() string {
	return fmt.Sprintf("TK%s%s", time.Now().Format("20060102"), tools.RandId()[:6])
}
