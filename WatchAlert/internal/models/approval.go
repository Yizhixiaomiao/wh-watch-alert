package models

// AutoEscalateRule 自动升级规则
type AutoEscalateRule struct {
	ID         string `json:"id" gorm:"column:id;primaryKey"`
	TenantId   string `json:"tenantId" gorm:"column:tenant_id;index"`
	Name       string `json:"name" gorm:"column:name"`
	Conditions struct {
		Status       string `json:"status" gorm:"column:status"`
		Priority     string `json:"priority" gorm:"column:priority"`
		OverdueHours int    `json:"overdueHours" gorm:"column:overdue_hours"`
	} `json:"conditions" gorm:"embedded;embeddedPrefix:condition_"`
	Actions struct {
		EscalateTo  string `json:"escalateTo" gorm:"column:escalate_to"`
		NotifyLevel string `json:"notifyLevel" gorm:"column:notify_level"`
	} `json:"actions" gorm:"embedded;embeddedPrefix:action_"`
	Enabled   bool  `json:"enabled" gorm:"column:enabled;default:true"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at"`
}

func (AutoEscalateRule) TableName() string {
	return "auto_escalate_rule"
}

// ApprovalWorkflow 审批工作流
type ApprovalWorkflow struct {
	ID        string         `json:"id" gorm:"column:id;primaryKey"`
	TenantId  string         `json:"tenantId" gorm:"column:tenant_id;index"`
	Name      string         `json:"name" gorm:"column:name"`
	Type      string         `json:"type" gorm:"column:type"`
	Steps     []ApprovalStep `json:"steps" gorm:"column:steps;serializer:json"`
	Enabled   bool           `json:"enabled" gorm:"column:enabled;default:true"`
	CreatedAt int64          `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt int64          `json:"updatedAt" gorm:"column:updated_at"`
}

func (ApprovalWorkflow) TableName() string {
	return "approval_workflow"
}

// ApprovalStep 审批步骤
type ApprovalStep struct {
	StepId     string   `json:"stepId" gorm:"column:step_id;primaryKey"`
	Name       string   `json:"name" gorm:"column:name"`
	Approvers  []string `json:"approvers" gorm:"column:approvers;serializer:json"`
	RequireAll bool     `json:"requireAll" gorm:"column:require_all;default:false"`
}

// ApprovalRequest 审批请求
type ApprovalRequest struct {
	ID          string `json:"id" gorm:"column:id;primaryKey"`
	WorkflowId  string `json:"workflowId" gorm:"column:workflow_id;index"`
	ResourceId  string `json:"resourceId" gorm:"column:resource_id;index"`
	Status      string `json:"status" gorm:"column:status;default:pending"` // pending, approved, rejected
	CurrentStep int    `json:"currentStep" gorm:"column:current_step;default:0"`
	CreatedBy   string `json:"createdBy" gorm:"column:created_by"`
	CreatedAt   int64  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   int64  `json:"updatedAt" gorm:"column:updated_at"`
}

func (ApprovalRequest) TableName() string {
	return "approval_request"
}

// ApprovalStepResult 审批步骤结果
type ApprovalStepResult struct {
	ID         string `json:"id" gorm:"column:id;primaryKey"`
	RequestId  string `json:"requestId" gorm:"column:request_id;index"`
	StepId     string `json:"stepId" gorm:"column:step_id"`
	ApproverId string `json:"approverId" gorm:"column:approver_id"`
	Approved   bool   `json:"approved" gorm:"column:approved"`
	Comment    string `json:"comment" gorm:"column:comment;type:text"`
	ApprovedAt int64  `json:"approvedAt" gorm:"column:approved_at"`
}

func (ApprovalStepResult) TableName() string {
	return "approval_step_result"
}
