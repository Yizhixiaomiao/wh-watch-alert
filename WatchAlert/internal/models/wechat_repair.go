package models

// WechatRepairRequest 微信报修请求表
type WechatRepairRequest struct {
	RequestId    string   `json:"requestId" gorm:"column:request_id;primaryKey"`
	TenantId     string   `json:"tenantId" gorm:"column:tenant_id;index:idx_tenant_id"`
	OpenId       string   `json:"openId" gorm:"column:"open_id;index:idx_open_id"`
	OpenType     string   `json:"openType" gorm:"column:"open_type;index:idx_open_type"`
	ContactName  string   `json:"contactName" gorm:"column:"contact_name"`
	ContactPhone string   `json:"contactPhone" gorm:"column:"contact_phone"`
	ContactEmail string   `json:"contactEmail" gorm:"column:"contact_email"`
	Location     string   `json:"location" gorm:"column:"location"`
	UrgentLevel  string   `json:"urgentLevel" gorm:"column:"urgent_level"`
	DeviceInfo   string   `json:"deviceInfo" gorm:"column:"device_info"`
	FaultType    string   `json:"faultType" gorm:"column:fault_type"`
	Description  string   `json:"description" gorm:"column:description;type:text"`
	Images       []string `json:"images" gorm:"column:images;serializer:json"`
	UserAgent    string   `json:"userAgent" gorm:"column:"user_agent"`
	TicketId     string   `json:"ticketId" gorm:"column:ticket_id;index:idx_ticket_id"`
	Status       string   `json:"status" gorm:"column:status;default:'pending'`
	TicketNo     string   `json:"ticketNo" gorm:"column:"ticket_no"`
	CreatedAt    int64    `json:"createdAt" gorm:"column:"created_at"`
	UpdatedAt    int64    `json:"updatedAt" gorm:"column:"updated_at"`
}

// TableName 指定表名
func (WechatRepairRequest) TableName() string {
	return "wechat_repair_request"
}
