package types

import "watchAlert/internal/models"

// RequestWechatRepairCreate 创建微信报修请求
type RequestWechatRepairCreate struct {
	TenantId     string   `json:"tenantId"`
	OpenId       string   `json:"openId"`
	OpenType     string   `json:"openType"`
	ContactName  string   `json:"contactName"`
	ContactPhone string   `json:"contactPhone"`
	ContactEmail string   `json:"contactEmail"`
	Location     string   `json:"location"`
	UrgentLevel  string   `json:"urgentLevel"`
	DeviceInfo   string   `json:"deviceInfo"`
	FaultType    string   `json:"faultType"`
	Description  string   `json:"description"`
	Images       []string `json:"images"`
	UserAgent    string   `json:"userAgent"`
	Platform     string   `json:"platform"`
}

// RequestWechatRepairQuery 查询微信报修请求
type RequestWechatRepairQuery struct {
	TenantId  string `json:"tenantId" form:"tenantId"`
	RequestId string `json:"requestId" form:"requestId"`
	OpenId    string `json:"openId" form:"openId"`
	Status    string `json:"status" form:"status"`
	Platform  string `json:"platform" form:"platform"`
	Page      int    `json:"page" form:"page"`
	Size      int    `json:"size" form:"size"`
}

// RequestWechatReply 微信回复请求
type RequestWechatReply struct {
	TenantId     string   `json:"tenantId"`
	RequestId    string   `json:"requestId" binding:"required"`
	ReplyContent string   `json:"replyContent" binding:"required"`
	ReplyedBy    string   `json:"replyedBy"`
	ReplyImages  []string `json:"replyImages"`
}

// ResponseWechatRepairList 微信报修列表响应
type ResponseWechatRepairList struct {
	List  []models.WechatRepairRequest `json:"list"`
	Total int64                        `json:"total"`
}

// ResponseWechatRepairDetail 微信报修详情响应
type ResponseWechatRepairDetail struct {
	RequestId    string   `json:"requestId"`
	TenantId     string   `json:"tenantId"`
	OpenId       string   `json:"openId"`
	OpenType     string   `json:"openType"`
	ContactName  string   `json:"contactName"`
	ContactPhone string   `json:"contactPhone"`
	ContactEmail string   `json:"contactEmail"`
	Location     string   `json:"location"`
	UrgentLevel  string   `json:"urgentLevel"`
	DeviceInfo   string   `json:"deviceInfo"`
	FaultType    string   `json:"faultType"`
	Description  string   `json:"description"`
	Images       []string `json:"images"`
	UserAgent    string   `json:"userAgent"`
	TicketId     string   `json:"ticketId"`
	TicketNo     string   `json:"ticketNo"`
	Status       string   `json:"status"`
	CreatedAt    int64    `json:"createdAt"`
	UpdatedAt    int64    `json:"updatedAt"`
}

// ResponseWechatReply 微信回复响应
type ResponseWechatReply struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	TicketNo string `json:"ticketNo"`
	TicketId string `json:"ticketId"`
}

// ResponseWechatRepairStatus 微信报修状态响应
type ResponseWechatRepairStatus struct {
	RequestId    string `json:"requestId"`
	Status       string `json:"status"`
	TicketNo     string `json:"ticketNo"`
	TicketId     string `json:"ticketId"`
	Message      string `json:"message"`
	ProcessStep  string `json:"processStep"`
	EstimateTime string `json:"estimateTime"`
}
