package models

// AlertTicketRule 告警转工单规则表
type AlertTicketRule struct {
	TenantId             string                       `json:"tenantId" gorm:"column:tenant_id;index:idx_tenant_id"`
	Id                   string                       `json:"id" gorm:"column:id;primaryKey"`
	Name                 string                       `json:"name" gorm:"column:name"`
	Description          string                       `json:"description" gorm:"column:description;type:text"`
	Priority             int                          `json:"priority" gorm:"column:priority;default:0;index:idx_priority"`
	FilterConditions     map[string]map[string]string `json:"filterConditions" gorm:"column:filter_conditions;serializer:json"`
	IsEnabled            bool                         `json:"isEnabled" gorm:"column:is_enabled;index:idx_is_enabled"`
	AutoAssign           bool                         `json:"autoAssign" gorm:"column:auto_assign"`
	DefaultAssignee      string                       `json:"defaultAssignee" gorm:"column:default_assignee"`
	DefaultAssigneeGroup string                       `json:"defaultAssigneeGroup" gorm:"column:default_assignee_group"`
	AutoClose            bool                         `json:"autoClose" gorm:"column:auto_close"`
	DuplicateRule        string                       `json:"duplicateRule" gorm:"column:duplicate_rule"` // skip, update, create
	TicketTemplateId     string                       `json:"ticketTemplateId" gorm:"column:ticket_template_id"`
	PriorityMapping      map[string]string            `json:"priorityMapping" gorm:"column:priority_mapping;serializer:json"`
	SeverityMapping      map[string]string            `json:"severityMapping" gorm:"column:severity_mapping;serializer:json"`
	CreatedBy            string                       `json:"createdBy" gorm:"column:created_by"`
	CreatedAt            int64                        `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt            int64                        `json:"updatedAt" gorm:"column:updated_at"`
}

// TableName 指定表名
func (AlertTicketRule) TableName() string {
	return "alert_ticket_rule"
}

// AlertTicketRuleHistory 告警转工单历史记录表
type AlertTicketRuleHistory struct {
	TenantId    string `json:"tenantId" gorm:"column:tenant_id;index:idx_tenant_id"`
	Id          string `json:"id" gorm:"column:id;primaryKey"`
	RuleId      string `json:"ruleId" gorm:"column:rule_id;index:idx_rule_id"`
	EventId     string `json:"eventId" gorm:"column:event_id;index:idx_event_id"`
	TicketId    string `json:"ticketId" gorm:"column:ticket_id;index:idx_ticket_id"`
	Action      string `json:"action" gorm:"column:action"` // create, update, skip
	Result      string `json:"result" gorm:"column:result"` // success, failed, skipped
	Reason      string `json:"reason" gorm:"column:reason;type:text"`
	ProcessTime int64  `json:"processTime" gorm:"column:process_time"`
	CreatedAt   int64  `json:"createdAt" gorm:"column:created_at;index:idx_created_at"`
}

// TableName 指定表名
func (AlertTicketRuleHistory) TableName() string {
	return "alert_ticket_rule_history"
}

// AlertTicketDuplicateRule 重复处理规则常量
const (
	DuplicateRuleSkip   = "skip"   // 跳过，不创建新工单
	DuplicateRuleUpdate = "update" // 更新现有工单
	DuplicateRuleCreate = "create" // 创建新工单
)

// GetFilterConditionTypes 获取过滤条件类型
func GetFilterConditionTypes() []string {
	return []string{
		"severity",        // 严重程度
		"datasource_type", // 数据源类型
		"rule_id",         // 规则ID
		"fault_center_id", // 故障中心ID
		"labels",          // 标签
	}
}

// GetFilterOperators 获取过滤操作符
func GetFilterOperators() []string {
	return []string{
		"equals",     // 等于
		"contains",   // 包含
		"not_equals", // 不等于
		"regex",      // 正则匹配
	}
}
