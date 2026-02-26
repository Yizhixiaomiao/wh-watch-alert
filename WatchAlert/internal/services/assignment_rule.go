package services

import (
	"fmt"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

type assignmentRuleService struct {
	ctx *ctx.Context
}

type InterAssignmentRuleService interface {
	// 规则操作
	CreateRule(req interface{}) (interface{}, interface{})
	UpdateRule(req interface{}) (interface{}, interface{})
	DeleteRule(req interface{}) (interface{}, interface{})
	GetRule(req interface{}) (interface{}, interface{})
	ListRules(req interface{}) (interface{}, interface{})
	MatchRule(req interface{}) (interface{}, interface{})
	AutoAssign(req interface{}) (interface{}, interface{})
}

func newInterAssignmentRuleService(ctx *ctx.Context) InterAssignmentRuleService {
	return &assignmentRuleService{ctx}
}

// CreateRule 创建规则
func (s assignmentRuleService) CreateRule(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAssignmentRuleCreate)

	rule := models.AssignmentRule{
		RuleId:         "ar-" + tools.RandId(),
		TenantId:       r.TenantId,
		Name:           r.Name,
		RuleType:       r.RuleType,
		AlertType:      r.AlertType,
		DataSource:     r.DataSource,
		Severity:       r.Severity,
		AssignmentType: r.AssignmentType,
		TargetUserId:   r.TargetUserId,
		TargetGroupId:  r.TargetGroupId,
		TargetDutyId:   r.TargetDutyId,
		Priority:       r.Priority,
		Enabled:        r.Enabled,
		CreatedBy:      r.CreatedBy,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
	}

	err := s.ctx.DB.AssignmentRule().CreateRule(rule)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// UpdateRule 更新规则
func (s assignmentRuleService) UpdateRule(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAssignmentRuleUpdate)

	rule, err := s.ctx.DB.AssignmentRule().GetRule(r.TenantId, r.RuleId)
	if err != nil {
		return nil, fmt.Errorf("规则不存在")
	}

	if r.Name != "" {
		rule.Name = r.Name
	}
	if r.RuleType != "" {
		rule.RuleType = r.RuleType
	}
	if r.AlertType != "" {
		rule.AlertType = r.AlertType
	}
	if r.DataSource != "" {
		rule.DataSource = r.DataSource
	}
	if r.Severity != "" {
		rule.Severity = r.Severity
	}
	if r.AssignmentType != "" {
		rule.AssignmentType = r.AssignmentType
	}
	if r.TargetUserId != "" {
		rule.TargetUserId = r.TargetUserId
	}
	if r.TargetGroupId != "" {
		rule.TargetGroupId = r.TargetGroupId
	}
	if r.TargetDutyId != "" {
		rule.TargetDutyId = r.TargetDutyId
	}
	if r.Priority > 0 {
		rule.Priority = r.Priority
	}
	if r.Enabled != nil {
		rule.Enabled = *r.Enabled
	}
	rule.UpdatedAt = time.Now().Unix()

	err = s.ctx.DB.AssignmentRule().UpdateRule(rule)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// DeleteRule 删除规则
func (s assignmentRuleService) DeleteRule(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAssignmentRuleDelete)

	err := s.ctx.DB.AssignmentRule().DeleteRule(r.TenantId, r.RuleId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetRule 获取规则详情
func (s assignmentRuleService) GetRule(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAssignmentRuleQuery)

	rule, err := s.ctx.DB.AssignmentRule().GetRule(r.TenantId, r.RuleId)
	if err != nil {
		return nil, err
	}

	return rule, nil
}

// ListRules 获取规则列表
func (s assignmentRuleService) ListRules(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAssignmentRuleQuery)

	rules, total, err := s.ctx.DB.AssignmentRule().ListRules(r.TenantId, r.RuleType, r.Enabled, r.Page, r.Size)
	if err != nil {
		return nil, err
	}

	return types.ResponseAssignmentRuleList{
		List:  rules,
		Total: total,
	}, nil
}

// MatchRule 匹配规则
func (s assignmentRuleService) MatchRule(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAssignmentRuleMatch)

	rules, err := s.ctx.DB.AssignmentRule().MatchRule(r.TenantId, r.AlertType, r.DataSource, r.Severity)
	if err != nil {
		return nil, err
	}

	if len(rules) == 0 {
		return types.ResponseAssignmentRuleMatch{
			Matched: false,
			Reason:  "未找到匹配的分配规则",
		}, nil
	}

	// 返回优先级最高的规则
	bestRule := rules[0]
	var assignee string
	var assignType string

	switch bestRule.AssignmentType {
	case models.AssignmentTargetTypeUser:
		assignee = bestRule.TargetUserId
		assignType = "user"
	case models.AssignmentTargetTypeGroup:
		assignee = bestRule.TargetGroupId
		assignType = "group"
	case models.AssignmentTargetTypeDuty:
		assignee, err = s.getDutyUser(bestRule.TargetDutyId)
		assignType = "duty"
		if err != nil {
			return types.ResponseAssignmentRuleMatch{
				Matched: false,
				Reason:  fmt.Sprintf("获取值班用户失败: %v", err),
			}, nil
		}
	}

	return types.ResponseAssignmentRuleMatch{
		Matched:    true,
		RuleName:   bestRule.Name,
		AssignTo:   assignee,
		AssignType: assignType,
		Assignee:   assignee,
		Reason:     "匹配成功",
		RulePreview: map[string]interface{}{
			"ruleId":         bestRule.RuleId,
			"ruleType":       bestRule.RuleType,
			"alertType":      bestRule.AlertType,
			"assignmentType": bestRule.AssignmentType,
		},
	}, nil
}

// AutoAssign 自动分配
func (s assignmentRuleService) AutoAssign(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestAutoAssign)

	// 获取工单信息
	ticket, err := s.ctx.DB.Ticket().Get(r.TenantId, r.TicketId)
	if err != nil {
		return nil, fmt.Errorf("工单不存在")
	}

	// 根据工单信息匹配规则
	alertType := string(ticket.Type)
	dataSource := ticket.DatasourceType
	severity := string(ticket.Severity)

	rules, err := s.ctx.DB.AssignmentRule().MatchRule(r.TenantId, alertType, dataSource, severity)
	if err != nil {
		return nil, fmt.Errorf("匹配规则失败: %v", err)
	}

	if len(rules) == 0 {
		// 没有匹配的规则，使用默认值班表
		return s.assignByDutySchedule(r.TenantId, r.TicketId)
	}

	// 使用优先级最高的规则
	bestRule := rules[0]
	var assignedTo string
	var assignedGroup string

	switch bestRule.AssignmentType {
	case models.AssignmentTargetTypeUser:
		assignedTo = bestRule.TargetUserId
	case models.AssignmentTargetTypeGroup:
		assignedGroup = bestRule.TargetGroupId
	case models.AssignmentTargetTypeDuty:
		assignedTo, err = s.getDutyUser(bestRule.TargetDutyId)
		if err != nil {
			return nil, fmt.Errorf("获取值班用户失败: %v", err)
		}
		assignedGroup = bestRule.TargetDutyId
	}

	// 更新工单
	ticket.AssignedTo = assignedTo
	ticket.AssignedGroup = assignedGroup
	ticket.Status = models.TicketStatusAssigned
	ticket.AssignedAt = time.Now().Unix()
	ticket.UpdatedAt = time.Now().Unix()

	err = s.ctx.DB.Ticket().Update(ticket)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ticketId":      r.TicketId,
		"assignedTo":    assignedTo,
		"assignedGroup": assignedGroup,
		"ruleId":        bestRule.RuleId,
		"ruleName":      bestRule.Name,
	}, nil
}

// getDutyUser 获取值班用户
func (s assignmentRuleService) getDutyUser(dutyId string) (string, error) {
	// 获取当前值班人员
	duty, err := s.ctx.DB.Duty().GetDuty(dutyId)
	if err != nil {
		return "", err
	}

	if len(duty.CurDutyUser) == 0 {
		return "", fmt.Errorf("当前无值班人员")
	}

	// 选择工单数最少的人员
	var selectedUser string
	minTicketCount := int64(1<<63 - 1)

	for _, user := range duty.CurDutyUser {
		// 统计该用户的当前工单数
		ticketCount, _ := s.ctx.DB.Ticket().CountByAssignedTo(user.UserId)
		if ticketCount < minTicketCount {
			minTicketCount = ticketCount
			selectedUser = user.UserId
		}
	}

	if selectedUser == "" {
		selectedUser = duty.CurDutyUser[0].UserId
	}

	return selectedUser, nil
}

// assignByDutySchedule 按值班表分配
func (s assignmentRuleService) assignByDutySchedule(tenantId, ticketId string) (interface{}, interface{}) {
	// 获取默认值班组
	rules, err := s.ctx.DB.AssignmentRule().GetEnabledRulesByType(tenantId, models.AssignmentRuleTypeDutySchedule)
	if err != nil || len(rules) == 0 {
		return nil, fmt.Errorf("未找到可用的值班表规则")
	}

	// 使用优先级最高的规则
	bestRule := rules[0]
	assignedTo, err := s.getDutyUser(bestRule.TargetDutyId)
	if err != nil {
		return nil, err
	}

	// 更新工单
	ticket, err := s.ctx.DB.Ticket().Get(tenantId, ticketId)
	if err != nil {
		return nil, err
	}

	ticket.AssignedTo = assignedTo
	ticket.AssignedGroup = bestRule.TargetDutyId
	ticket.Status = models.TicketStatusAssigned
	ticket.AssignedAt = time.Now().Unix()
	ticket.UpdatedAt = time.Now().Unix()

	err = s.ctx.DB.Ticket().Update(ticket)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ticketId":      ticketId,
		"assignedTo":    assignedTo,
		"assignedGroup": bestRule.TargetDutyId,
		"ruleId":        bestRule.RuleId,
		"ruleName":      bestRule.Name,
	}, nil
}
