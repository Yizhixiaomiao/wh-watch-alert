package types

import "watchAlert/internal/models"

// RequestTicketCreate 创建工单请求
type RequestTicketCreate struct {
	TenantId        string                 `json:"tenantId"`
	Title           string                 `json:"title" binding:"required"`
	Description     string                 `json:"description"`
	Type            models.TicketType      `json:"type" binding:"required"`
	Priority        models.TicketPriority  `json:"priority" binding:"required"`
	Severity        models.TicketSeverity  `json:"severity"`
	Source          models.TicketSource    `json:"source"`
	EventId         string                 `json:"eventId"`
	FaultCenterId   string                 `json:"faultCenterId"`
	RuleId          string                 `json:"ruleId"`
	DatasourceType  string                 `json:"datasourceType"`
	AssignedTo      string                 `json:"assignedTo"`
	AssignedGroup   string                 `json:"assignedGroup"`
	Followers       []string               `json:"followers"`
	Labels          map[string]string      `json:"labels"`
	Tags            []string               `json:"tags"`
	CustomFields    map[string]interface{} `json:"customFields"`
	CreatedBy       string                 `json:"createdBy"`
	RelatedTicketId string                 `json:"relatedTicketId"` // 关联工单ID
	RelationType    string                 `json:"relationType"`    // 关联类型
}

// RequestTicketUpdate 更新工单请求
type RequestTicketUpdate struct {
	TenantId      string                 `json:"tenantId"`
	TicketId      string                 `json:"ticketId" binding:"required"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	Priority      models.TicketPriority  `json:"priority"`
	Severity      models.TicketSeverity  `json:"severity"`
	AssignedTo    string                 `json:"assignedTo"`
	AssignedGroup string                 `json:"assignedGroup"`
	Followers     []string               `json:"followers"`
	Labels        map[string]string      `json:"labels"`
	Tags          []string               `json:"tags"`
	CustomFields  map[string]interface{} `json:"customFields"`
	RootCause     string                 `json:"rootCause"`
	Solution      string                 `json:"solution"`
}

// RequestTicketDelete 删除工单请求
type RequestTicketDelete struct {
	TenantId string `json:"tenantId"`
	TicketId string `json:"ticketId" binding:"required"`
}

// RequestTicketQuery 查询工单请求
type RequestTicketQuery struct {
	TenantId      string                `json:"tenantId" form:"tenantId"`
	TicketId      string                `json:"ticketId" form:"ticketId"`
	TicketNo      string                `json:"ticketNo" form:"ticketNo"`
	Status        models.TicketStatus   `json:"status" form:"status"`
	Priority      models.TicketPriority `json:"priority" form:"priority"`
	Type          models.TicketType     `json:"type" form:"type"`
	AssignedTo    string                `json:"assignedTo" form:"assignedTo"`
	CreatedBy     string                `json:"createdBy" form:"createdBy"`
	EventId       string                `json:"eventId" form:"eventId"`
	FaultCenterId string                `json:"faultCenterId" form:"faultCenterId"`
	Keyword       string                `json:"keyword" form:"keyword"`
	StartTime     int64                 `json:"startTime" form:"startTime"`
	EndTime       int64                 `json:"endTime" form:"endTime"`
	Page          int                   `json:"page" form:"page"`
	Size          int                   `json:"size" form:"size"`
}

// RequestTicketAssign 分配工单请求
type RequestTicketAssign struct {
	TenantId      string `json:"tenantId"`
	TicketId      string `json:"ticketId" binding:"required"`
	AssignedTo    string `json:"assignedTo" binding:"required"`
	AssignedGroup string `json:"assignedGroup"`
	Reason        string `json:"reason"`
	UserId        string `json:"userId"`
	UserName      string `json:"userName"`
}

// RequestTicketClaim 认领工单请求
type RequestTicketClaim struct {
	TenantId string `json:"tenantId"`
	TicketId string `json:"ticketId" binding:"required"`
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
}

// RequestTicketTransfer 转派工单请求
type RequestTicketTransfer struct {
	TenantId   string `json:"tenantId"`
	TicketId   string `json:"ticketId" binding:"required"`
	TransferTo string `json:"transferTo" binding:"required"`
	Reason     string `json:"reason" binding:"required"`
	UserId     string `json:"userId"`
	UserName   string `json:"userName"`
}

// RequestTicketEscalate 升级工单请求
type RequestTicketEscalate struct {
	TenantId   string `json:"tenantId"`
	TicketId   string `json:"ticketId" binding:"required"`
	EscalateTo string `json:"escalateTo" binding:"required"`
	Reason     string `json:"reason" binding:"required"`
	UserId     string `json:"userId"`
	UserName   string `json:"userName"`
}

// RequestTicketResolve 标记解决请求
type RequestTicketResolve struct {
	TenantId  string `json:"tenantId"`
	TicketId  string `json:"ticketId" binding:"required"`
	RootCause string `json:"rootCause"`
	Solution  string `json:"solution" binding:"required"`
	UserId    string `json:"userId"`
	UserName  string `json:"userName"`
}

// RequestTicketClose 关闭工单请求
type RequestTicketClose struct {
	TenantId string `json:"tenantId"`
	TicketId string `json:"ticketId" binding:"required"`
	Reason   string `json:"reason"`
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
}

// RequestTicketReopen 重新打开工单请求
type RequestTicketReopen struct {
	TenantId string `json:"tenantId"`
	TicketId string `json:"ticketId" binding:"required"`
	Reason   string `json:"reason" binding:"required"`
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
}

// RequestTicketComment 添加评论请求
type RequestTicketComment struct {
	TenantId string   `json:"tenantId"`
	TicketId string   `json:"ticketId" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	Mentions []string `json:"mentions"`
	UserId   string   `json:"userId"`
	UserName string   `json:"userName"`
}

// RequestTicketCommentQuery 查询评论请求
type RequestTicketCommentQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	TicketId string `json:"ticketId" form:"ticketId" binding:"required"`
	Page     int    `json:"page" form:"page"`
	Size     int    `json:"size" form:"size"`
}

// RequestTicketWorkLogQuery 查询工作日志请求
type RequestTicketWorkLogQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	TicketId string `json:"ticketId" form:"ticketId" binding:"required"`
	Page     int    `json:"page" form:"page"`
	Size     int    `json:"size" form:"size"`
}

// RequestTicketStatistics 工单统计请求
type RequestTicketStatistics struct {
	TenantId  string `json:"tenantId" form:"tenantId"`
	StartTime int64  `json:"startTime" form:"startTime"`
	EndTime   int64  `json:"endTime" form:"endTime"`
}

// RequestTicketTemplateCreate 创建工单模板请求
type RequestTicketTemplateCreate struct {
	TenantId        string                 `json:"tenantId"`
	Name            string                 `json:"name" binding:"required"`
	Type            models.TicketType      `json:"type" binding:"required"`
	TitleTemplate   string                 `json:"titleTemplate" binding:"required"`
	DescTemplate    string                 `json:"descTemplate"`
	DefaultPriority models.TicketPriority  `json:"defaultPriority"`
	DefaultAssignee string                 `json:"defaultAssignee"`
	CustomFields    map[string]interface{} `json:"customFields"`
	CreatedBy       string                 `json:"createdBy"`
}

// RequestTicketTemplateUpdate 更新工单模板请求
type RequestTicketTemplateUpdate struct {
	TenantId        string                 `json:"tenantId"`
	Id              string                 `json:"id" binding:"required"`
	Name            string                 `json:"name"`
	TitleTemplate   string                 `json:"titleTemplate"`
	DescTemplate    string                 `json:"descTemplate"`
	DefaultPriority models.TicketPriority  `json:"defaultPriority"`
	DefaultAssignee string                 `json:"defaultAssignee"`
	CustomFields    map[string]interface{} `json:"customFields"`
}

// RequestTicketTemplateDelete 删除工单模板请求
type RequestTicketTemplateDelete struct {
	TenantId string `json:"tenantId"`
	Id       string `json:"id" binding:"required"`
}

// RequestTicketTemplateQuery 查询工单模板请求
type RequestTicketTemplateQuery struct {
	TemplateId string `json:"templateId"`
	TenantId   string `json:"tenantId"`
	Type       string `json:"type"`
	Page       int    `json:"page"`
	Size       int    `json:"size"`
}

// RequestTicketSLAPolicyCreate 创建SLA策略请求
type RequestTicketSLAPolicyCreate struct {
	TenantId       string                `json:"tenantId"`
	Name           string                `json:"name" binding:"required"`
	Priority       models.TicketPriority `json:"priority" binding:"required"`
	ResponseTime   int64                 `json:"responseTime" binding:"required"`
	ResolutionTime int64                 `json:"resolutionTime" binding:"required"`
	WorkingHours   string                `json:"workingHours"`
	Holidays       []string              `json:"holidays"`
	Enabled        *bool                 `json:"enabled"`
	CreatedBy      string                `json:"createdBy"`
}

func (r *RequestTicketSLAPolicyCreate) GetEnabled() *bool {
	if r.Enabled == nil {
		enabled := true
		return &enabled
	}
	return r.Enabled
}

// RequestTicketSLAPolicyUpdate 更新SLA策略请求
type RequestTicketSLAPolicyUpdate struct {
	TenantId       string                `json:"tenantId"`
	Id             string                `json:"id" binding:"required"`
	Name           string                `json:"name"`
	Priority       models.TicketPriority `json:"priority"`
	ResponseTime   int64                 `json:"responseTime"`
	ResolutionTime int64                 `json:"resolutionTime"`
	WorkingHours   string                `json:"workingHours"`
	Holidays       []string              `json:"holidays"`
	Enabled        *bool                 `json:"enabled"`
}

// RequestTicketSLAPolicyDelete 删除SLA策略请求
type RequestTicketSLAPolicyDelete struct {
	TenantId string `json:"tenantId"`
	Id       string `json:"id" binding:"required"`
}

// RequestTicketSLAPolicyQuery 查询SLA策略请求
type RequestTicketSLAPolicyQuery struct {
	TenantId string                `json:"tenantId" form:"tenantId"`
	Priority models.TicketPriority `json:"priority" form:"priority"`
	Enabled  *bool                 `json:"enabled" form:"enabled"`
	Page     int                   `json:"page" form:"page"`
	Size     int                   `json:"size" form:"size"`
}

// ResponseTicketList 工单列表响应
type ResponseTicketList struct {
	List  []models.Ticket `json:"list"`
	Total int64           `json:"total"`
}

// ResponseTicketCommentList 评论列表响应
type ResponseTicketCommentList struct {
	List  []models.TicketComment `json:"list"`
	Total int64                  `json:"total"`
}

// ResponseTicketWorkLogList 工作日志列表响应
type ResponseTicketWorkLogList struct {
	List  []models.TicketWorkLog `json:"list"`
	Total int64                  `json:"total"`
}

// ResponseTicketTemplateList 工单模板列表响应
type ResponseTicketTemplateList struct {
	List  []models.TicketTemplate `json:"list"`
	Total int64                   `json:"total"`
}

// ResponseTicketSLAPolicyList SLA策略列表响应
type ResponseTicketSLAPolicyList struct {
	List  []models.TicketSLAPolicy `json:"list"`
	Total int64                    `json:"total"`
}

// ResponseTicketStatistics 工单统计响应
type ResponseTicketStatistics struct {
	TotalCount      int64                     `json:"totalCount"`
	PendingCount    int64                     `json:"pendingCount"`
	ProcessingCount int64                     `json:"processingCount"`
	ClosedCount     int64                     `json:"closedCount"`
	OverdueCount    int64                     `json:"overdueCount"`
	PriorityStats   map[string]int64          `json:"priorityStats"`
	TypeStats       map[string]int64          `json:"typeStats"`
	StatusStats     map[string]int64          `json:"statusStats"`
	AvgResponseTime int64                     `json:"avgResponseTime"`
	AvgResolution   int64                     `json:"avgResolution"`
	SLARate         float64                   `json:"slaRate"`
	UserStats       []ResponseTicketUserStats `json:"userStats"`
	TrendData       []ResponseTicketTrendData `json:"trendData"`
}

// ResponseTicketUserStats 用户统计
type ResponseTicketUserStats struct {
	UserId         string  `json:"userId"`
	UserName       string  `json:"userName"`
	TicketCount    int64   `json:"ticketCount"`
	AvgResponseTime int64  `json:"avgResponseTime"`
	AvgResolution  int64   `json:"avgResolution"`
	SLARate        float64 `json:"slaRate"`
}

// ResponseTicketTrendData 趋势数据
type ResponseTicketTrendData struct {
	Date    string `json:"date"`
	Count   int64  `json:"count"`
	Resolved int64  `json:"resolved"`
}

// WechatTemplateMessage 微信模板消息
type WechatTemplateMessage struct {
	ToUser      string                 `json:"touser"`
	TemplateId  string                 `json:"template_id"`
	Url         string                 `json:"url"`
	Data        map[string]interface{} `json:"data"`
	Miniprogram *WechatMiniprogram     `json:"miniprogram,omitempty"`
}

// WechatMiniprogram 微信小程序跳转信息
type WechatMiniprogram struct {
	AppId    string `json:"appid"`
	PagePath string `json:"pagepath"`
}

// RequestMobileTicketCreate 移动端创建工单请求
type RequestMobileTicketCreate struct {
	TenantId     string                 `json:"tenantId"`
	Title        string                 `json:"title" binding:"required"`
	Description  string                 `json:"description"`
	Type         models.TicketType      `json:"type"`
	Priority     models.TicketPriority  `json:"priority"`
	ContactName  string                 `json:"contactName" binding:"required"`
	ContactPhone string                 `json:"contactPhone" binding:"required"`
	ContactEmail string                 `json:"contactEmail"`
	Location     string                 `json:"location"`
	UrgentLevel  string                 `json:"urgentLevel"`
	DeviceInfo   string                 `json:"deviceInfo"`
	FaultType    string                 `json:"faultType"`
	Images       []string               `json:"images"`
	UserAgent    string                 `json:"userAgent"`
	Platform     string                 `json:"platform"`
	Labels       map[string]string      `json:"labels"`
	CustomFields map[string]interface{} `json:"customFields"`
}

// ResponseMobileTicketCreate 移动端创建工单响应
type ResponseMobileTicketCreate struct {
	TicketId   string `json:"ticketId"`
	TicketNo   string `json:"ticketNo"`
	Message    string `json:"message"`
	WaitTime   string `json:"waitTime"`
	ProcessUrl string `json:"processUrl"`
}

// RequestMobileTicketQuery 移动端工单查询请求
type RequestMobileTicketQuery struct {
	TenantId     string              `json:"tenantId" form:"tenantId"`
	TicketId     string              `json:"ticketId" form:"ticketId"`
	TicketNo     string              `json:"ticketNo" form:"ticketNo"`
	ContactPhone string              `json:"contactPhone" form:"contactPhone"`
	Status       models.TicketStatus `json:"status" form:"status"`
	Page         int                 `json:"page" form:"page"`
	Size         int                 `json:"size" form:"size"`
}

// ResponseMobileTicketStatus 移动端工单状态响应
type ResponseMobileTicketStatus struct {
	TicketId     string                `json:"ticketId"`
	TicketNo     string                `json:"ticketNo"`
	Title        string                `json:"title"`
	Status       models.TicketStatus   `json:"status"`
	Priority     models.TicketPriority `json:"priority"`
	CreatedAt    int64                 `json:"createdAt"`
	AssignedTo   string                `json:"assignedTo"`
	ProcessStep  string                `json:"processStep"`
	EstimateTime string                `json:"estimateTime"`
}

// RequestAlertTicketRuleCreate 创建告警转工单规则请求
type RequestAlertTicketRuleCreate struct {
	TenantId             string                       `json:"tenantId"`
	Name                 string                       `json:"name" binding:"required"`
	Description          string                       `json:"description"`
	Priority             int                          `json:"priority"`
	FilterConditions     map[string]map[string]string `json:"filterConditions"`
	IsEnabled            *bool                        `json:"isEnabled"`
	AutoAssign           *bool                        `json:"autoAssign"`
	DefaultAssignee      string                       `json:"defaultAssignee"`
	DefaultAssigneeGroup string                       `json:"defaultAssigneeGroup"`
	AutoClose            *bool                        `json:"autoClose"`
	DuplicateRule        string                       `json:"duplicateRule"`
	TicketTemplateId     string                       `json:"ticketTemplateId"`
	PriorityMapping      map[string]string            `json:"priorityMapping"`
	SeverityMapping      map[string]string            `json:"severityMapping"`
	CreatedBy            string                       `json:"createdBy"`
}

func (r *RequestAlertTicketRuleCreate) GetIsEnabled() bool {
	if r.IsEnabled == nil {
		return true
	}
	return *r.IsEnabled
}

func (r *RequestAlertTicketRuleCreate) GetAutoAssign() bool {
	if r.AutoAssign == nil {
		return false
	}
	return *r.AutoAssign
}

func (r *RequestAlertTicketRuleCreate) GetAutoClose() bool {
	if r.AutoClose == nil {
		return false
	}
	return *r.AutoClose
}

// RequestAlertTicketRuleUpdate 更新告警转工单规则请求
type RequestAlertTicketRuleUpdate struct {
	TenantId             string                       `json:"tenantId"`
	Id                   string                       `json:"id" binding:"required"`
	Name                 string                       `json:"name"`
	Description          string                       `json:"description"`
	FilterConditions     map[string]map[string]string `json:"filterConditions"`
	IsEnabled            *bool                        `json:"isEnabled"`
	AutoAssign           *bool                        `json:"autoAssign"`
	DefaultAssignee      string                       `json:"defaultAssignee"`
	DefaultAssigneeGroup string                       `json:"defaultAssigneeGroup"`
	AutoClose            *bool                        `json:"autoClose"`
	DuplicateRule        string                       `json:"duplicateRule"`
	TicketTemplateId     string                       `json:"ticketTemplateId"`
	PriorityMapping      map[string]string            `json:"priorityMapping"`
	SeverityMapping      map[string]string            `json:"severityMapping"`
}

// RequestAlertTicketRuleDelete 删除告警转工单规则请求
type RequestAlertTicketRuleDelete struct {
	TenantId string `json:"tenantId"`
	Id       string `json:"id" binding:"required"`
}

// RequestAlertTicketRuleQuery 查询告警转工单规则请求
type RequestAlertTicketRuleQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	Id       string `json:"id" form:"id"`
	Page     int    `json:"page" form:"page"`
	Size     int    `json:"size" form:"size"`
}

// ResponseAlertTicketRuleList 告警转工单规则列表响应
type ResponseAlertTicketRuleList struct {
	List  []models.AlertTicketRule `json:"list"`
	Total int64                    `json:"total"`
}

// RequestAlertTicketRuleTest 测试规则匹配请求
type RequestAlertTicketRuleTest struct {
	TenantId  string                 `json:"tenantId"`
	RuleId    string                 `json:"ruleId"`
	AlertData map[string]interface{} `json:"alertData"`
}

// ResponseAlertTicketRuleTest 测试规则匹配响应
type ResponseAlertTicketRuleTest struct {
	Matched       bool                   `json:"matched"`
	RuleName      string                 `json:"ruleName"`
	TicketPreview map[string]interface{} `json:"ticketPreview"`
	Reason        string                 `json:"reason"`
}

// RequestAlertTicketRuleHistoryQuery 查询规则历史记录请求
type RequestAlertTicketRuleHistoryQuery struct {
	TenantId  string `json:"tenantId" form:"tenantId"`
	RuleId    string `json:"ruleId" form:"ruleId"`
	EventId   string `json:"eventId" form:"eventId"`
	Result    string `json:"result" form:"result"`
	StartTime int64  `json:"startTime" form:"startTime"`
	EndTime   int64  `json:"endTime" form:"endTime"`
	Page      int    `json:"page" form:"page"`
	Size      int    `json:"size" form:"size"`
}

// ResponseAlertTicketRuleHistory 规则历史记录响应
type ResponseAlertTicketRuleHistory struct {
	List  []models.AlertTicketRuleHistory `json:"list"`
	Total int64                           `json:"total"`
}

// RequestAlertTicketRuleStats 规则统计请求
type RequestAlertTicketRuleStats struct {
	TenantId  string `json:"tenantId" form:"tenantId"`
	RuleId    string `json:"ruleId" form:"ruleId"`
	StartTime int64  `json:"startTime" form:"startTime"`
	EndTime   int64  `json:"endTime" form:"endTime"`
}

// ResponseAlertTicketRuleStats 规则统计响应
type ResponseAlertTicketRuleStats struct {
	TotalRules          int64            `json:"totalRules"`
	EnabledRules        int64            `json:"enabledRules"`
	TotalMatches        int64            `json:"totalMatches"`
	TotalTicketsCreated int64            `json:"totalTicketsCreated"`
	SuccessRate         float64          `json:"successRate"`
	RuleUsage           []RuleUsageStats `json:"ruleUsage"`
	FailureReasons      map[string]int64 `json:"failureReasons"`
}

// RuleUsageStats 规则使用统计
type RuleUsageStats struct {
	RuleId       string `json:"ruleId"`
	RuleName     string `json:"ruleName"`
	MatchCount   int64  `json:"matchCount"`
	TicketCount  int64  `json:"ticketCount"`
	SuccessCount int64  `json:"successCount"`
	FailCount    int64  `json:"failCount"`
}

// RequestTicketStepCreate 创建处理步骤请求
type RequestTicketStepCreate struct {
	TenantId     string   `json:"tenantId"`
	TicketId     string   `json:"ticketId" binding:"required"`
	Order        int      `json:"order" binding:"required"`
	Title        string   `json:"title" binding:"required"`
	Description  string   `json:"description"`
	Method       string   `json:"method"`
	Result       string   `json:"result"`
	Attachments  []string `json:"attachments"`
	CreatedBy    string   `json:"createdBy"`
	KnowledgeIds []string `json:"knowledgeIds"` // 关联的知识ID列表
}

// RequestTicketStepUpdate 更新处理步骤请求
type RequestTicketStepUpdate struct {
	TenantId     string   `json:"tenantId"`
	TicketId     string   `json:"ticketId" binding:"required"`
	StepId       string   `json:"stepId" binding:"required"`
	Order        int      `json:"order"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Method       string   `json:"method"`
	Result       string   `json:"result"`
	Attachments  []string `json:"attachments"`
	KnowledgeIds []string `json:"knowledgeIds"` // 关联的知识ID列表
}

// RequestTicketStepDelete 删除处理步骤请求
type RequestTicketStepDelete struct {
	TenantId string `json:"tenantId"`
	TicketId string `json:"ticketId" binding:"required"`
	StepId   string `json:"stepId" binding:"required"`
}

// RequestTicketStepReorder 重新排序步骤请求
type RequestTicketStepReorder struct {
	TenantId string          `json:"tenantId"`
	TicketId string          `json:"ticketId" binding:"required"`
	Steps    []StepOrderItem `json:"steps" binding:"required"`
}

// StepOrderItem 步骤排序项
type StepOrderItem struct {
	StepId string `json:"stepId"`
	Order  int    `json:"order"`
}

// RequestTicketStepsQuery 查询处理步骤请求
type RequestTicketStepsQuery struct {
	TenantId string `json:"tenantId" form:"tenantId"`
	TicketId string `json:"ticketId" form:"ticketId" binding:"required"`
}
