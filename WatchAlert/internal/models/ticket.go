package models

// TicketStatus 工单状态类型
type TicketStatus string

// 工单状态定义
const (
	TicketStatusPending    TicketStatus = "Pending"    // 待处理
	TicketStatusAssigned   TicketStatus = "Assigned"   // 已分配
	TicketStatusProcessing TicketStatus = "Processing" // 处理中
	TicketStatusVerifying  TicketStatus = "Verifying"  // 验证中
	TicketStatusResolved   TicketStatus = "Resolved"   // 已解决
	TicketStatusClosed     TicketStatus = "Closed"     // 已关闭
	TicketStatusCancelled  TicketStatus = "Cancelled"  // 已取消
	TicketStatusEscalated  TicketStatus = "Escalated"  // 已升级
)

// TicketType 工单类型
type TicketType string

const (
	TicketTypeAlert  TicketType = "Alert"  // 告警工单
	TicketTypeFault  TicketType = "Fault"  // 故障工单
	TicketTypeChange TicketType = "Change" // 变更工单
	TicketTypeQuery  TicketType = "Query"  // 咨询工单
)

// TicketPriority 工单优先级
type TicketPriority string

const (
	TicketPriorityP0 TicketPriority = "P0" // 最高优先级
	TicketPriorityP1 TicketPriority = "P1" // 高优先级
	TicketPriorityP2 TicketPriority = "P2" // 中优先级
	TicketPriorityP3 TicketPriority = "P3" // 低优先级
	TicketPriorityP4 TicketPriority = "P4" // 最低优先级
)

// TicketSeverity 工单严重程度
type TicketSeverity string

const (
	TicketSeverityCritical TicketSeverity = "Critical" // 严重
	TicketSeverityHigh     TicketSeverity = "High"     // 高
	TicketSeverityMedium   TicketSeverity = "Medium"   // 中
	TicketSeverityLow      TicketSeverity = "Low"      // 低
)

// TicketSource 工单来源
type TicketSource string

const (
	TicketSourceAuto   TicketSource = "auto"   // 自动创建
	TicketSourceManual TicketSource = "manual" // 手动创建
	TicketSourceAPI    TicketSource = "api"    // API创建
)

// Ticket 工单主表
type Ticket struct {
	// 基础字段
	TenantId    string         `json:"tenantId" gorm:"column:tenant_id;index:idx_tenant_status"`
	TicketId    string         `json:"ticketId" gorm:"column:ticket_id;primaryKey"`
	TicketNo    string         `json:"ticketNo" gorm:"column:ticket_no;uniqueIndex;type:varchar(255)"`
	Title       string         `json:"title" gorm:"column:title;type:varchar(255)"`
	Description string         `json:"description" gorm:"column:description;type:text"`
	Type        TicketType     `json:"type" gorm:"column:type;type:varchar(50)"`
	Priority    TicketPriority `json:"priority" gorm:"column:priority;type:varchar(10)"`
	Severity    TicketSeverity `json:"severity" gorm:"column:severity;type:varchar(20)"`
	Status      TicketStatus   `json:"status" gorm:"column:status;index:idx_tenant_status"`

	// 关联信息
	Source         TicketSource `json:"source" gorm:"column:source"`
	EventId        string       `json:"eventId" gorm:"column:event_id;index:idx_event_id"`
	FaultCenterId  string       `json:"faultCenterId" gorm:"column:fault_center_id"`
	RuleId         string       `json:"ruleId" gorm:"column:rule_id"`
	DatasourceType string       `json:"datasourceType" gorm:"column:datasource_type"`

	// 人员信息
	CreatedBy     string   `json:"createdBy" gorm:"column:created_by"`
	AssignedTo    string   `json:"assignedTo" gorm:"column:assigned_to;index:idx_assigned_to"`
	AssignedGroup string   `json:"assignedGroup" gorm:"column:assigned_group"`
	Followers     []string `json:"followers" gorm:"column:followers;serializer:json"`

	// 时间信息
	CreatedAt       int64 `json:"createdAt" gorm:"column:created_at;index:idx_created_at"`
	UpdatedAt       int64 `json:"updatedAt" gorm:"column:updated_at"`
	AssignedAt      int64 `json:"assignedAt" gorm:"column:assigned_at"`
	FirstResponseAt int64 `json:"firstResponseAt" gorm:"column:first_response_at"`
	ResolvedAt      int64 `json:"resolvedAt" gorm:"column:resolved_at"`
	ClosedAt        int64 `json:"closedAt" gorm:"column:closed_at"`
	DueTime         int64 `json:"dueTime" gorm:"column:due_time;index:idx_due_time"`

	// SLA 相关
	SlaPolicy     string `json:"slaPolicy" gorm:"column:sla_policy"`
	ResponseSLA   int64  `json:"responseSLA" gorm:"column:response_sla"`
	ResolutionSLA int64  `json:"resolutionSLA" gorm:"column:resolution_sla"`
	IsOverdue     bool   `json:"isOverdue" gorm:"column:is_overdue"`

	// 扩展信息
	Labels       map[string]string      `json:"labels" gorm:"column:labels;serializer:json"`
	Tags         []string               `json:"tags" gorm:"column:tags;serializer:json"`
	CustomFields map[string]interface{} `json:"customFields" gorm:"column:custom_fields;serializer:json"`

	// 处理信息
	RootCause string `json:"rootCause" gorm:"column:root_cause;type:text"`
	Solution  string `json:"solution" gorm:"column:solution;type:text"`

	// 统计信息
	ResponseTime   int64 `json:"responseTime" gorm:"column:response_time"`
	ResolutionTime int64 `json:"resolutionTime" gorm:"column:resolution_time"`
	ReopenCount    int   `json:"reopenCount" gorm:"column:reopen_count"`

	// 告警状态同步字段
	AlarmActive  bool  `json:"alarmActive" gorm:"column:alarm_active;default:true"`
	LastSyncTime int64 `json:"lastSyncTime" gorm:"column:last_sync_time;default:0"`

	// 工单关联字段
	RelatedTicketId string `json:"relatedTicketId" gorm:"column:related_ticket_id;default:''"`
	RelationType    string `json:"relationType" gorm:"column:relation_type;default:''"`           // "alarm_to_repair", "repair_to_alarm"
	KnowledgeId     string `json:"knowledgeId" gorm:"column:knowledge_id;index:idx_knowledge_id"` // 工单生成的知识ID

	// 处理步骤
	Steps []TicketStep `json:"steps" gorm:"column:steps;serializer:json"`
}

// TableName 指定表名
func (Ticket) TableName() string {
	return "ticket"
}

type TicketWorkLog struct {
	Id        string `json:"id" gorm:"column:id;primaryKey"`
	TicketId  string `json:"ticketId" gorm:"column:ticket_id;index:idx_ticket_id"`
	UserId    string `json:"userId" gorm:"column:user_id"`
	UserName  string `json:"userName" gorm:"column:user_name"`
	Action    string `json:"action" gorm:"column:action"`
	Content   string `json:"content" gorm:"column:content;type:text"`
	OldValue  string `json:"oldValue" gorm:"column:old_value"`
	NewValue  string `json:"newValue" gorm:"column:new_value"`
	TimeSpent int64  `json:"timeSpent" gorm:"column:time_spent"`
	CreatedAt int64  `json:"createdAt" gorm:"column:created_at;index:idx_created_at"`
}

// TableName 指定表名
func (TicketWorkLog) TableName() string {
	return "ticket_work_log"
}

// TicketComment 评论表
type TicketComment struct {
	Id        string   `json:"id" gorm:"column:id;primaryKey"`
	TicketId  string   `json:"ticketId" gorm:"column:ticket_id;index:idx_ticket_id"`
	UserId    string   `json:"userId" gorm:"column:user_id"`
	UserName  string   `json:"userName" gorm:"column:user_name"`
	Content   string   `json:"content" gorm:"column:content;type:text"`
	Mentions  []string `json:"mentions" gorm:"column:mentions;serializer:json"`
	CreatedAt int64    `json:"createdAt" gorm:"column:created_at;index:idx_created_at"`
	UpdatedAt int64    `json:"updatedAt" gorm:"column:updated_at"`
}

// TableName 指定表名
func (TicketComment) TableName() string {
	return "ticket_comment"
}

// TicketAttachment 附件表
type TicketAttachment struct {
	Id        string `json:"id" gorm:"column:id;primaryKey"`
	TicketId  string `json:"ticketId" gorm:"column:ticket_id;index:idx_ticket_id"`
	FileName  string `json:"fileName" gorm:"column:file_name"`
	FileSize  int64  `json:"fileSize" gorm:"column:file_size"`
	FileType  string `json:"fileType" gorm:"column:file_type"`
	FilePath  string `json:"filePath" gorm:"column:file_path"`
	UploadBy  string `json:"uploadBy" gorm:"column:upload_by"`
	CreatedAt int64  `json:"createdAt" gorm:"column:created_at"`
}

// TableName 指定表名
func (TicketAttachment) TableName() string {
	return "ticket_attachment"
}

// TicketTemplate 工单模板表
type TicketTemplate struct {
	TenantId        string                 `json:"tenantId" gorm:"column:tenant_id;index:idx_tenant_id"`
	Id              string                 `json:"id" gorm:"column:id;primaryKey"`
	Name            string                 `json:"name" gorm:"column:name"`
	Type            TicketType             `json:"type" gorm:"column:type"`
	TitleTemplate   string                 `json:"titleTemplate" gorm:"column:title_template"`
	DescTemplate    string                 `json:"descTemplate" gorm:"column:desc_template;type:text"`
	DefaultPriority TicketPriority         `json:"defaultPriority" gorm:"column:default_priority"`
	DefaultAssignee string                 `json:"defaultAssignee" gorm:"column:default_assignee"`
	CustomFields    map[string]interface{} `json:"customFields" gorm:"column:custom_fields;serializer:json"`
	CreatedBy       string                 `json:"createdBy" gorm:"column:created_by"`
	CreatedAt       int64                  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt       int64                  `json:"updatedAt" gorm:"column:updated_at"`
}

// TableName 指定表名
func (TicketTemplate) TableName() string {
	return "ticket_template"
}

// TicketSLAPolicy SLA策略表
type TicketSLAPolicy struct {
	TenantId       string         `json:"tenantId" gorm:"column:tenant_id;index:idx_tenant_id"`
	Id             string         `json:"id" gorm:"column:id;primaryKey"`
	Name           string         `json:"name" gorm:"column:name"`
	Priority       TicketPriority `json:"priority" gorm:"column:priority"`
	ResponseTime   int64          `json:"responseTime" gorm:"column:response_time"`
	ResolutionTime int64          `json:"resolutionTime" gorm:"column:resolution_time"`
	WorkingHours   string         `json:"workingHours" gorm:"column:working_hours"`
	Holidays       []string       `json:"holidays" gorm:"column:holidays;serializer:json"`
	Enabled        *bool          `json:"enabled" gorm:"column:enabled"`
	CreatedBy      string         `json:"createdBy" gorm:"column:created_by"`
	CreatedAt      int64          `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt      int64          `json:"updatedAt" gorm:"column:updated_at"`
}

// TableName 指定表名
func (TicketSLAPolicy) TableName() string {
	return "ticket_sla_policy"
}

// TicketStep 工单处理步骤
type TicketStep struct {
	StepId       string   `json:"stepId" gorm:"column:step_id"`
	TicketId     string   `json:"ticketId" gorm:"column:ticket_id;index:idx_ticket_id"`
	Order        int      `json:"order" gorm:"column:order"`
	Title        string   `json:"title" gorm:"column:title"`
	Description  string   `json:"description" gorm:"column:description;type:text"`
	Method       string   `json:"method" gorm:"column:method;type:text"`
	Result       string   `json:"result" gorm:"column:result;type:text"`
	Attachments  []string `json:"attachments" gorm:"column:attachments;serializer:json"`
	CreatedBy    string   `json:"createdBy" gorm:"column:created_by"`
	CreatedAt    int64    `json:"createdAt" gorm:"column:created_at"`
	TenantId     string   `json:"tenantId" gorm:"column:tenant_id;index:idx_tenant_id"`
	KnowledgeIds []string `json:"knowledgeIds" gorm:"column:knowledge_ids;serializer:json"` // 关联的知识ID列表
}

// TableName 指定表名
func (TicketStep) TableName() string {
	return "ticket_step"
}
