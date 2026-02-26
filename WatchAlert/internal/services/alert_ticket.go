package services

import (
	"fmt"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/repo"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

// alertTicketService 告警转工单服务
type alertTicketService struct {
	ctx *ctx.Context
}

type InterAlertTicketService interface {
	// CreateTicketFromAlert 从告警创建工单
	CreateTicketFromAlert(alert *models.AlertCurEvent) error
	// ShouldCreateTicket 判断是否需要创建工单
	ShouldCreateTicket(alert *models.AlertCurEvent) bool
	// GetTicketByEventId 根据事件ID获取工单
	GetTicketByEventId(tenantId, eventId string) (*models.Ticket, error)
	// CreateAlertTicketRule 创建告警转工单规则
	CreateAlertTicketRule(req interface{}) (interface{}, interface{})
	// UpdateAlertTicketRule 更新告警转工单规则
	UpdateAlertTicketRule(req interface{}) (interface{}, interface{})
	// DeleteAlertTicketRule 删除告警转工单规则
	DeleteAlertTicketRule(req interface{}) (interface{}, interface{})
	// GetAlertTicketRule 获取告警转工单规则
	GetAlertTicketRule(req interface{}) (interface{}, interface{})
	// ListAlertTicketRules 获取告警转工单规则列表
	ListAlertTicketRules(req interface{}) (interface{}, interface{})
	// TestAlertTicketRule 测试规则匹配
	TestAlertTicketRule(req interface{}) (interface{}, interface{})
	// GetAlertTicketRuleHistory 获取规则历史记录
	GetAlertTicketRuleHistory(req interface{}) (interface{}, interface{})
	// GetAlertTicketRuleStats 获取规则统计
	GetAlertTicketRuleStats(req interface{}) (interface{}, interface{})
}

func newInterAlertTicketService(ctx *ctx.Context) InterAlertTicketService {
	return &alertTicketService{ctx: ctx}
}

// CreateTicketFromAlert 从告警创建工单
func (s *alertTicketService) CreateTicketFromAlert(alert *models.AlertCurEvent) error {
	// 检查是否已存在工单
	existingTicket, err := s.GetTicketByEventId(alert.TenantId, alert.EventId)
	if err == nil && existingTicket != nil {
		// 工单已存在，检查是否需要更新
		return s.updateExistingTicket(existingTicket, alert)
	}

	// 获取告警转工单规则（按优先级排序）
	rule, err := s.getAlertTicketRuleByPriority(alert)
	if err != nil {
		// 没有规则，使用默认逻辑
		if !s.shouldCreateTicketByDefault(alert) {
			return nil
		}
		rule = &models.AlertTicketRule{
			IsEnabled:     true,
			AutoAssign:    false,
			AutoClose:     false,
			DuplicateRule: "skip",
		}
	}

	if !rule.IsEnabled {
		return nil // 规则未启用，跳过创建
	}

	// 生成工单ID和编号
	ticketId := "tk-" + tools.RandId()
	ticketNo := generateTicketNo()

	// 根据告警严重程度确定工单优先级
	priority := s.mapSeverityToPriority(alert.Severity)
	severity := s.mapSeverityToTicketSeverity(alert.Severity)

	// 获取SLA策略
	var responseSLA, resolutionSLA int64
	var dueTime int64
	slaPolicy, err := s.ctx.DB.Ticket().GetSLAPolicyByPriority(alert.TenantId, priority)
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
	title, description := s.buildTicketContent(alert)

	// 获取默认处理人
	var assignedTo, assignedGroup string
	if rule.AutoAssign {
		assignedTo = rule.DefaultAssignee
		assignedGroup = rule.DefaultAssigneeGroup
	}

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
		TenantId:       alert.TenantId,
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
		AssignedTo:     assignedTo,
		AssignedGroup:  assignedGroup,
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

	// 如果指定了处理人，设置状态为处理中
	if assignedTo != "" {
		ticket.Status = models.TicketStatusProcessing
		ticket.AssignedAt = time.Now().Unix()
	}

	err = s.ctx.DB.Ticket().Create(ticket)
	if err != nil {
		return fmt.Errorf("创建告警工单失败: %v", err)
	}

	// 记录历史
	s.recordRuleHistory(rule, alert.EventId, ticketId, "create", "success", "")

	// 创建工作日志
	s.createWorkLog(ticketId, "system", "create", "告警自动创建工单", "", "")

	// 如果需要自动分配，记录分配日志
	if rule.AutoAssign && assignedTo != "" {
		s.createWorkLog(ticketId, "system", "assign", fmt.Sprintf("自动分配给 %s", assignedTo), "", assignedTo)
	}

	return nil
}

// ShouldCreateTicket 判断是否需要创建工单
func (s *alertTicketService) ShouldCreateTicket(alert *models.AlertCurEvent) bool {
	// 只在告警状态时创建工单
	if alert.Status != models.StateAlerting {
		return false
	}

	// 检查是否有对应的规则
	rule, err := s.getAlertTicketRule(alert)
	if err != nil {
		// 没有规则，使用默认逻辑
		return s.shouldCreateTicketByDefault(alert)
	}

	return rule.IsEnabled
}

// GetTicketByEventId 根据事件ID获取工单
func (s *alertTicketService) GetTicketByEventId(tenantId, eventId string) (*models.Ticket, error) {
	// 通过工单服务的查询方法获取
	query := repo.TicketQuery{
		EventId: eventId,
		Page:    1,
		Size:    1,
	}
	tickets, _, err := s.ctx.DB.Ticket().List(tenantId, query)
	if err != nil {
		return nil, err
	}
	if len(tickets) == 0 {
		return nil, fmt.Errorf("工单不存在")
	}
	return &tickets[0], nil
}

// CreateAlertTicketRule 创建告警转工单规则
func (s *alertTicketService) CreateAlertTicketRule(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAlertTicketRuleCreate)

	rule := models.AlertTicketRule{
		TenantId:             r.TenantId,
		Id:                   "atr-" + tools.RandId(),
		Name:                 r.Name,
		Description:          r.Description,
		Priority:             0,
		FilterConditions:     r.FilterConditions,
		IsEnabled:            r.GetIsEnabled(),
		AutoAssign:           r.GetAutoAssign(),
		DefaultAssignee:      r.DefaultAssignee,
		DefaultAssigneeGroup: r.DefaultAssigneeGroup,
		AutoClose:            r.GetAutoClose(),
		DuplicateRule:        r.DuplicateRule,
		TicketTemplateId:     r.TicketTemplateId,
		PriorityMapping:      r.PriorityMapping,
		SeverityMapping:      r.SeverityMapping,
		CreatedBy:            r.CreatedBy,
		CreatedAt:            time.Now().Unix(),
		UpdatedAt:            time.Now().Unix(),
	}

	err := s.ctx.DB.AlertTicketRule().Create(rule)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// UpdateAlertTicketRule 更新告警转工单规则
func (s *alertTicketService) UpdateAlertTicketRule(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAlertTicketRuleUpdate)

	rule, err := s.ctx.DB.AlertTicketRule().Get(r.TenantId, r.Id)
	if err != nil {
		return nil, fmt.Errorf("规则不存在")
	}

	if r.Name != "" {
		rule.Name = r.Name
	}
	if r.Description != "" {
		rule.Description = r.Description
	}
	if r.FilterConditions != nil {
		rule.FilterConditions = r.FilterConditions
	}
	if r.IsEnabled != nil {
		rule.IsEnabled = *r.IsEnabled
	}
	if r.AutoAssign != nil {
		rule.AutoAssign = *r.AutoAssign
	}
	if r.DefaultAssignee != "" {
		rule.DefaultAssignee = r.DefaultAssignee
	}
	if r.DefaultAssigneeGroup != "" {
		rule.DefaultAssigneeGroup = r.DefaultAssigneeGroup
	}
	if r.AutoClose != nil {
		rule.AutoClose = *r.AutoClose
	}
	if r.DuplicateRule != "" {
		rule.DuplicateRule = r.DuplicateRule
	}
	if r.TicketTemplateId != "" {
		rule.TicketTemplateId = r.TicketTemplateId
	}
	if r.PriorityMapping != nil {
		rule.PriorityMapping = r.PriorityMapping
	}
	if r.SeverityMapping != nil {
		rule.SeverityMapping = r.SeverityMapping
	}
	rule.UpdatedAt = time.Now().Unix()

	err = s.ctx.DB.AlertTicketRule().Update(rule)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// DeleteAlertTicketRule 删除告警转工单规则
func (s *alertTicketService) DeleteAlertTicketRule(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAlertTicketRuleDelete)
	err := s.ctx.DB.AlertTicketRule().Delete(r.TenantId, r.Id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// GetAlertTicketRule 获取告警转工单规则
func (s *alertTicketService) GetAlertTicketRule(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAlertTicketRuleQuery)
	rule, err := s.ctx.DB.AlertTicketRule().Get(r.TenantId, r.Id)
	if err != nil {
		return nil, err
	}
	return rule, nil
}

// ListAlertTicketRules 获取告警转工单规则列表
func (s *alertTicketService) ListAlertTicketRules(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAlertTicketRuleQuery)
	rules, total, err := s.ctx.DB.AlertTicketRule().List(r.TenantId, r.Page, r.Size)
	if err != nil {
		return nil, err
	}

	return types.ResponseAlertTicketRuleList{
		List:  rules,
		Total: total,
	}, nil
}

// 内部辅助方法

// updateExistingTicket 更新现有工单
func (s *alertTicketService) updateExistingTicket(ticket *models.Ticket, alert *models.AlertCurEvent) error {
	// 根据重复规则处理
	rule, _ := s.getAlertTicketRule(alert)
	duplicateRule := "skip" // 默认跳过
	if rule != nil {
		duplicateRule = rule.DuplicateRule
	}

	switch duplicateRule {
	case "skip":
		// 跳过，不创建新工单
		return nil
	case "update":
		// 更新现有工单
		ticket.UpdatedAt = time.Now().Unix()
		ticket.CustomFields["alert_last_update"] = time.Now().Unix()
		err := s.ctx.DB.Ticket().Update(*ticket)
		if err != nil {
			return err
		}
		// 记录更新日志
		s.createWorkLog(ticket.TicketId, "system", "update", "告警更新工单信息", "", "")
		return nil
	case "create":
		// 创建新工单（继续执行）
		return nil
	default:
		return nil
	}
}

// mapSeverityToPriority 将告警严重程度映射为工单优先级
func (s *alertTicketService) mapSeverityToPriority(severity string) models.TicketPriority {
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
func (s *alertTicketService) mapSeverityToTicketSeverity(severity string) models.TicketSeverity {
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
func (s *alertTicketService) buildTicketContent(alert *models.AlertCurEvent) (title, description string) {
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

// getAlertTicketRule 获取告警转工单规则
func (s *alertTicketService) getAlertTicketRule(alert *models.AlertCurEvent) (*models.AlertTicketRule, error) {
	rules, _, err := s.ctx.DB.AlertTicketRule().List(alert.TenantId, 1, 100) // 获取所有规则
	if err != nil {
		return nil, err
	}

	// 查找第一个匹配的规则
	for _, rule := range rules {
		if s.matchRule(alert, &rule) {
			return &rule, nil
		}
	}

	return nil, fmt.Errorf("没有匹配的规则")
}

// getAlertTicketRuleByPriority 按优先级获取告警转工单规则
func (s *alertTicketService) getAlertTicketRuleByPriority(alert *models.AlertCurEvent) (*models.AlertTicketRule, error) {
	rules, _, err := s.ctx.DB.AlertTicketRule().List(alert.TenantId, 1, 100) // 获取所有规则
	if err != nil {
		return nil, err
	}

	// 查找所有匹配的规则
	var matchedRules []models.AlertTicketRule
	for _, rule := range rules {
		if s.matchRule(alert, &rule) {
			matchedRules = append(matchedRules, rule)
		}
	}

	if len(matchedRules) == 0 {
		return nil, fmt.Errorf("没有匹配的规则")
	}

	// 按优先级排序（优先级数字越小，优先级越高）
	for i := 0; i < len(matchedRules)-1; i++ {
		for j := i + 1; j < len(matchedRules); j++ {
			if matchedRules[i].Priority > matchedRules[j].Priority {
				matchedRules[i], matchedRules[j] = matchedRules[j], matchedRules[i]
			}
		}
	}

	// 返回优先级最高的规则
	return &matchedRules[0], nil
}

// recordRuleHistory 记录规则历史
func (s *alertTicketService) recordRuleHistory(rule *models.AlertTicketRule, eventId, ticketId, action, result, reason string) {
	history := models.AlertTicketRuleHistory{
		TenantId:    rule.TenantId,
		Id:          "h-" + tools.RandId(),
		RuleId:      rule.Id,
		EventId:     eventId,
		TicketId:    ticketId,
		Action:      action,
		Result:      result,
		Reason:      reason,
		ProcessTime: time.Now().Unix(),
		CreatedAt:   time.Now().Unix(),
	}
	s.ctx.DB.AlertTicketRule().CreateHistory(history)
}

// matchRule 检查告警是否匹配规则
func (s *alertTicketService) matchRule(alert *models.AlertCurEvent, rule *models.AlertTicketRule) bool {
	// 如果规则没有启用，跳过
	if !rule.IsEnabled {
		return false
	}

	// 如果没有过滤条件，匹配所有
	if rule.FilterConditions == nil || len(rule.FilterConditions) == 0 {
		return true
	}

	// 检查过滤条件
	for key, condition := range rule.FilterConditions {
		value := ""
		switch key {
		case "severity":
			value = alert.Severity
		case "datasource_type":
			value = alert.DatasourceType
		case "rule_id":
			value = alert.RuleId
		case "fault_center_id":
			value = alert.FaultCenterId
		default:
			// 检查标签
			if alertVal, ok := alert.Labels[key]; ok {
				if strVal, ok := alertVal.(string); ok {
					value = strVal
				}
			}
		}

		// 简单的字符串匹配
		if condition["operator"] == "equals" && value != condition["value"] {
			return false
		}
		if condition["operator"] == "contains" && !contains(value, condition["value"]) {
			return false
		}
	}

	return true
}

// shouldCreateTicketByDefault 默认的创建工单逻辑
func (s *alertTicketService) shouldCreateTicketByDefault(alert *models.AlertCurEvent) bool {
	// 默认只对critical和warning级别的告警创建工单
	return alert.Severity == "critical" || alert.Severity == "warning"
}

// createWorkLog 创建工作日志
func (s *alertTicketService) createWorkLog(ticketId, userId, action, content, oldValue, newValue string) {
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
	s.ctx.DB.Ticket().CreateWorkLog(log)
}

// contains 检查字符串包含关系
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}

// TestAlertTicketRule 测试规则匹配
func (s *alertTicketService) TestAlertTicketRule(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAlertTicketRuleTest)

	// 获取规则
	rule, err := s.ctx.DB.AlertTicketRule().Get(r.TenantId, r.RuleId)
	if err != nil {
		return nil, fmt.Errorf("规则不存在")
	}

	// 构造模拟告警事件
	alert := &models.AlertCurEvent{
		TenantId: r.TenantId,
	}

	// 从 AlertData 中提取字段
	if severity, ok := r.AlertData["severity"].(string); ok {
		alert.Severity = severity
	}
	if datasourceType, ok := r.AlertData["datasource_type"].(string); ok {
		alert.DatasourceType = datasourceType
	}
	if ruleId, ok := r.AlertData["rule_id"].(string); ok {
		alert.RuleId = ruleId
	}
	if faultCenterId, ok := r.AlertData["fault_center_id"].(string); ok {
		alert.FaultCenterId = faultCenterId
	}
	if labels, ok := r.AlertData["labels"].(map[string]interface{}); ok {
		alert.Labels = labels
	}

	// 检查是否匹配
	matched := s.matchRule(alert, &rule)

	result := types.ResponseAlertTicketRuleTest{
		Matched:  matched,
		RuleName: rule.Name,
	}

	if matched {
		result.Reason = "规则匹配成功"
		// 生成工单预览
		priority := s.mapSeverityToPriority(alert.Severity)
		severity := s.mapSeverityToTicketSeverity(alert.Severity)
		title, description := s.buildTicketContent(alert)

		result.TicketPreview = map[string]interface{}{
			"title":       title,
			"priority":    priority,
			"severity":    severity,
			"description": description,
		}
	} else {
		result.Reason = "规则不匹配"
	}

	return result, nil
}

// GetAlertTicketRuleHistory 获取规则历史记录
func (s *alertTicketService) GetAlertTicketRuleHistory(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAlertTicketRuleHistoryQuery)

	histories, total, err := s.ctx.DB.AlertTicketRule().ListHistory(r.TenantId, r.Page, r.Size)
	if err != nil {
		return nil, err
	}

	// 如果指定了规则ID，过滤结果
	if r.RuleId != "" {
		filtered := make([]models.AlertTicketRuleHistory, 0)
		for _, h := range histories {
			if h.RuleId == r.RuleId {
				filtered = append(filtered, h)
			}
		}
		total = int64(len(filtered))
		histories = filtered
	}

	// 如果指定了结果，过滤结果
	if r.Result != "" {
		filtered := make([]models.AlertTicketRuleHistory, 0)
		for _, h := range histories {
			if h.Result == r.Result {
				filtered = append(filtered, h)
			}
		}
		total = int64(len(filtered))
		histories = filtered
	}

	return types.ResponseAlertTicketRuleHistory{
		List:  histories,
		Total: total,
	}, nil
}

// GetAlertTicketRuleStats 获取规则统计
func (s *alertTicketService) GetAlertTicketRuleStats(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAlertTicketRuleStats)

	// 获取所有规则
	rules, _, err := s.ctx.DB.AlertTicketRule().List(r.TenantId, 1, 1000)
	if err != nil {
		return nil, err
	}

	// 获取历史记录
	histories, _, err := s.ctx.DB.AlertTicketRule().ListHistory(r.TenantId, 1, 10000)
	if err != nil {
		return nil, err
	}

	// 统计数据
	stats := types.ResponseAlertTicketRuleStats{
		TotalRules:          int64(len(rules)),
		EnabledRules:        0,
		TotalMatches:        0,
		TotalTicketsCreated: 0,
		RuleUsage:           make([]types.RuleUsageStats, 0),
		FailureReasons:      make(map[string]int64),
	}

	// 启用的规则数量
	for _, rule := range rules {
		if rule.IsEnabled {
			stats.EnabledRules++
		}
	}

	// 构建规则使用统计
	ruleUsageMap := make(map[string]*types.RuleUsageStats)
	for _, rule := range rules {
		ruleUsageMap[rule.Id] = &types.RuleUsageStats{
			RuleId:   rule.Id,
			RuleName: rule.Name,
		}
	}

	// 统计历史记录
	for _, history := range histories {
		stats.TotalMatches++

		if history.Result == "success" {
			stats.TotalTicketsCreated++
			if usage, ok := ruleUsageMap[history.RuleId]; ok {
				usage.MatchCount++
				usage.TicketCount++
				usage.SuccessCount++
			}
		} else if history.Result == "failed" {
			if usage, ok := ruleUsageMap[history.RuleId]; ok {
				usage.MatchCount++
				usage.FailCount++
			}
			// 统计失败原因
			if history.Reason != "" {
				stats.FailureReasons[history.Reason]++
			}
		}
	}

	// 计算成功率
	if stats.TotalMatches > 0 {
		stats.SuccessRate = float64(stats.TotalTicketsCreated) / float64(stats.TotalMatches) * 100
	}

	// 转换规则使用统计为列表
	for _, usage := range ruleUsageMap {
		stats.RuleUsage = append(stats.RuleUsage, *usage)
	}

	return stats, nil
}
