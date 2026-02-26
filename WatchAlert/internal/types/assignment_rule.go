package types

import "watchAlert/internal/models"

// RequestAssignmentRuleCreate 创建分配规则请求
type RequestAssignmentRuleCreate struct {
	TenantId       string                      `json:"tenantId" binding:"required"`
	Name           string                      `json:"name" binding:"required"`
	RuleType       models.AssignmentRuleType   `json:"ruleType" binding:"required"`
	AlertType      string                      `json:"alertType"`
	DataSource     string                      `json:"dataSource"`
	Severity       string                      `json:"severity"`
	AssignmentType models.AssignmentTargetType `json:"assignmentType" binding:"required"`
	TargetUserId   string                      `json:"targetUserId"`
	TargetGroupId  string                      `json:"targetGroupId"`
	TargetDutyId   string                      `json:"targetDutyId"`
	Priority       int                         `json:"priority"`
	Enabled        bool                        `json:"enabled"`
	CreatedBy      string                      `json:"createdBy"`
}

// RequestAssignmentRuleUpdate 更新分配规则请求
type RequestAssignmentRuleUpdate struct {
	TenantId       string                      `json:"tenantId"`
	RuleId         string                      `json:"ruleId" binding:"required"`
	Name           string                      `json:"name"`
	RuleType       models.AssignmentRuleType   `json:"ruleType"`
	AlertType      string                      `json:"alertType"`
	DataSource     string                      `json:"dataSource"`
	Severity       string                      `json:"severity"`
	AssignmentType models.AssignmentTargetType `json:"assignmentType"`
	TargetUserId   string                      `json:"targetUserId"`
	TargetGroupId  string                      `json:"targetGroupId"`
	TargetDutyId   string                      `json:"targetDutyId"`
	Priority       int                         `json:"priority"`
	Enabled        *bool                       `json:"enabled"`
}

// RequestAssignmentRuleDelete 删除分配规则请求
type RequestAssignmentRuleDelete struct {
	TenantId string `json:"tenantId"`
	RuleId   string `json:"ruleId" binding:"required"`
}

// RequestAssignmentRuleQuery 查询分配规则请求
type RequestAssignmentRuleQuery struct {
	TenantId string                    `json:"tenantId" form:"tenantId"`
	RuleId   string                    `json:"ruleId" form:"ruleId"`
	RuleType models.AssignmentRuleType `json:"ruleType" form:"ruleType"`
	Enabled  *bool                     `json:"enabled" form:"enabled"`
	Page     int                       `json:"page" form:"page"`
	Size     int                       `json:"size" form:"size"`
}

// RequestAssignmentRuleMatch 匹配分配规则请求
type RequestAssignmentRuleMatch struct {
	TenantId   string `json:"tenantId" binding:"required"`
	AlertType  string `json:"alertType"`
	DataSource string `json:"dataSource"`
	Severity   string `json:"severity"`
	TicketId   string `json:"ticketId"`
}

// RequestAutoAssign 自动分配请求
type RequestAutoAssign struct {
	TenantId string `json:"tenantId" binding:"required"`
	TicketId string `json:"ticketId" binding:"required"`
}

// ResponseAssignmentRuleList 分配规则列表响应
type ResponseAssignmentRuleList struct {
	List  []models.AssignmentRule `json:"list"`
	Total int64                   `json:"total"`
}

// ResponseAssignmentRuleMatch 分配规则匹配响应
type ResponseAssignmentRuleMatch struct {
	Matched     bool                   `json:"matched"`
	RuleName    string                 `json:"ruleName"`
	AssignTo    string                 `json:"assignTo"`
	AssignType  string                 `json:"assignType"`
	Assignee    string                 `json:"assignee"`
	Reason      string                 `json:"reason"`
	RulePreview map[string]interface{} `json:"rulePreview"`
}
