package models

type AssignmentRuleType string

const (
	AssignmentRuleTypeAlertType    AssignmentRuleType = "alert_type"    // 告警类型规则
	AssignmentRuleTypeDutySchedule AssignmentRuleType = "duty_schedule" // 值班表规则
)

type AssignmentTargetType string

const (
	AssignmentTargetTypeUser  AssignmentTargetType = "user"  // 分配给用户
	AssignmentTargetTypeGroup AssignmentTargetType = "group" // 分配给组
	AssignmentTargetTypeDuty  AssignmentTargetType = "duty"  // 分配给值班组
)

// AssignmentRule 分配规则表
type AssignmentRule struct {
	RuleId   string             `json:"ruleId" gorm:"column:rule_id;primaryKey"`
	TenantId string             `json:"tenantId" gorm:"column:tenant_id;index:idx_tenant_id"`
	Name     string             `json:"name" gorm:"column:name"`
	RuleType AssignmentRuleType `json:"ruleType" gorm:"column:rule_type;index:idx_rule_type"`

	// 告警类型规则
	AlertType  string `json:"alertType" gorm:"column:alert_type;index:idx_alert_type"`
	DataSource string `json:"dataSource" gorm:"column:data_source;index:idx_data_source"`
	Severity   string `json:"severity" gorm:"column:severity;index:idx_severity"`

	// 分配规则
	AssignmentType AssignmentTargetType `json:"assignmentType" gorm:"column:assignment_type"`
	TargetUserId   string               `json:"targetUserId" gorm:"column:target_user_id"`
	TargetGroupId  string               `json:"targetGroupId" gorm:"column:target_group_id"`
	TargetDutyId   string               `json:"targetDutyId" gorm:"column:target_duty_id"`

	Priority  int    `json:"priority" gorm:"column:priority;index:idx_priority"`
	Enabled   bool   `json:"enabled" gorm:"column:enabled;index:idx_enabled;default:true"`
	CreatedBy string `json:"createdBy" gorm:"column:created_by"`
	CreatedAt int64  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt int64  `json:"updatedAt" gorm:"column:updated_at"`
}

// TableName 指定表名
func (AssignmentRule) TableName() string {
	return "assignment_rule"
}
