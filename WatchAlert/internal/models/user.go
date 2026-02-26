package models

type Member struct {
	UserId       string   `json:"userid"`
	UserName     string   `json:"username"`
	Email        string   `json:"email"`
	Phone        string   `json:"phone"`
	Password     string   `json:"password"`
	Role         string   `json:"role"`
	CreateBy     string   `json:"create_by"`
	CreateAt     int64    `json:"create_at"`
	JoinDuty     string   `json:"joinDuty"`
	DutyUserId   string   `json:"dutyUserId"`
	Tenants      []string `json:"tenants" gorm:"tenants;serializer:json"`
	Status       string   `json:"status" gorm:"column:status;default:'enabled'"`
	StatusReason string   `json:"statusReason" gorm:"column:status_reason;type:text"`
}

type ResponseLoginInfo struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	UserId   string `json:"userId"`
}

// UserStatusHistory 用户状态变更历史
type UserStatusHistory struct {
	Id           string `json:"id" gorm:"column:id;primaryKey"`
	UserId       string `json:"userId" gorm:"column:user_id;index:idx_user_id"`
	OldStatus    string `json:"oldStatus" gorm:"column:old_status"`
	NewStatus    string `json:"newStatus" gorm:"column:new_status;index:idx_status"`
	Reason       string `json:"reason" gorm:"column:reason"`
	OperatorId   string `json:"operatorId" gorm:"column:operator_id;index:idx_operator_id"`
	OperatorName string `json:"operatorName" gorm:"column:operator_name"`
	CreatedAt    int64  `json:"createdAt" gorm:"column:created_at;index:idx_created_at"`
}

// UserActivityLog 用户活动审计日志
type UserActivityLog struct {
	Id           string `json:"id" gorm:"column:id;primaryKey"`
	UserId       string `json:"userId" gorm:"column:user_id;index:idx_user_id"`
	UserName     string `json:"userName" gorm:"column:user_name"`
	Action       string `json:"action" gorm:"column:action;index:idx_action"`
	ResourceType string `json:"resourceType" gorm:"column:resource_type;index:idx_resource_type"`
	ResourceId   string `json:"resourceId" gorm:"column:resource_id"`
	ResourceName string `json:"resourceName" gorm:"column:resource_name"`
	IpAddress    string `json:"ipAddress" gorm:"column:ip_address"`
	UserAgent    string `json:"userAgent" gorm:"column:user_agent"`
	Details      string `json:"details" gorm:"column:details;type:text"`
	Status       string `json:"status" gorm:"column:status"`
	ErrorMessage string `json:"errorMessage" gorm:"column:error_message;type:text"`
	CreatedAt    int64  `json:"createdAt" gorm:"column:created_at;index:idx_created_at"`
}

// TableName 指定表名
func (ResponseLoginInfo) TableName() string {
	return "response_login_info"
}

func (UserStatusHistory) TableName() string {
	return "user_status_history"
}

func (UserActivityLog) TableName() string {
	return "user_activity_log"
}
