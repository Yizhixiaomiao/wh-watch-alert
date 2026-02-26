package types

type RequestUserLogin struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type RequestUserCreate struct {
	UserId     string   `json:"userid"`
	UserName   string   `json:"username"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Password   string   `json:"password"`
	Role       string   `json:"role"`
	CreateBy   string   `json:"create_by"`
	CreateAt   int64    `json:"create_at"`
	JoinDuty   string   `json:"joinDuty" `
	DutyUserId string   `json:"dutyUserId"`
	Tenants    []string `json:"tenants" gorm:"tenants;serializer:json"`
}

type RequestUserUpdate struct {
	UserId     string   `json:"userid"`
	UserName   string   `json:"username"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Password   string   `json:"password"`
	Role       string   `json:"role"`
	CreateBy   string   `json:"create_by"`
	CreateAt   int64    `json:"create_at"`
	JoinDuty   string   `json:"joinDuty" `
	DutyUserId string   `json:"dutyUserId"`
	Tenants    []string `json:"tenants" gorm:"tenants;serializer:json"`
}

type RequestUserQuery struct {
	UserId   string `json:"userid" form:"userid"`
	UserName string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Phone    string `json:"phone" form:"phone"`
	Query    string `json:"query" form:"query"`
	JoinDuty string `json:"joinDuty" form:"joinDuty"`
	TenantId string `json:"tenantId" form:"tenantId"`
}

type RequestUserChangePassword struct {
	UserId   string `json:"userid"`
	Password string `json:"password"`
}

type RequestUserStatusUpdate struct {
	UserId       string `json:"userid"`
	Status       string `json:"status"`
	StatusReason string `json:"statusReason"`
}

type RequestUserStatusQuery struct {
	UserId   string `json:"userid" form:"userid"`
	UserName string `json:"username" form:"username"`
	Status   string `json:"status" form:"status"`
	Page     int    `json:"page" form:"page"`
	Size     int    `json:"size" form:"size"`
}

type RequestUserPermissionsQuery struct {
	UserId   string `json:"userid" form:"userid" binding:"required"`
	TenantId string `json:"tenantId" form:"tenantId"`
}

type RequestUserBatchOperation struct {
	UserIds      []string `json:"userIds"`
	Operation    string   `json:"operation"`    // status, role, delete
	Status       string   `json:"status"`       // for status operation
	StatusReason string   `json:"statusReason"` // for status operation
	Role         string   `json:"role"`         // for role operation
}

type RequestUserActivityLogQuery struct {
	UserId       string `json:"userid" form:"userid"`
	UserName     string `json:"username" form:"username"`
	Action       string `json:"action" form:"action"`
	ResourceType string `json:"resourceType" form:"resourceType"`
	Page         int    `json:"page" form:"page"`
	Size         int    `json:"size" form:"size"`
}

type RequestUserActivityLogCreate struct {
	UserId       string `json:"userId"`
	UserName     string `json:"userName"`
	Action       string `json:"action"`
	ResourceType string `json:"resourceType"`
	ResourceId   string `json:"resourceId"`
	ResourceName string `json:"resourceName"`
	IpAddress    string `json:"ipAddress"`
	UserAgent    string `json:"userAgent"`
	Details      string `json:"details"`
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}
