package types

import "watchAlert/internal/models"

// RequestTicketReviewAssign 分配评委请求
type RequestTicketReviewAssign struct {
	TenantId    string   `json:"tenantId"`
	TicketId    string   `json:"ticketId" binding:"required"`
	ReviewerIds []string `json:"reviewerIds" binding:"required"`
	AssignedBy  string   `json:"assignedBy"`
}

// RequestTicketReviewSubmit 提交评审请求
type RequestTicketReviewSubmit struct {
	TenantId   string  `json:"tenantId"`
	ReviewId   string  `json:"reviewId" binding:"required"`
	Rating     int     `json:"rating" binding:"required,min=1,max=5"`
	WorkHours  float64 `json:"workHours" binding:"required,min=0"`
	Comment    string  `json:"comment"`
	ReviewedBy string  `json:"reviewedBy"`
}

// RequestTicketReviewQuery 查询评审请求
type RequestTicketReviewQuery struct {
	TenantId   string                    `json:"tenantId" form:"tenantId"`
	TicketId   string                    `json:"ticketId" form:"ticketId"`
	ReviewerId string                    `json:"reviewerId" form:"reviewerId"`
	Status     models.TicketReviewStatus `json:"status" form:"status"`
	Page       int                       `json:"page" form:"page"`
	Size       int                       `json:"size" form:"size"`
}

// ResponseTicketReviewList 评审列表响应
type ResponseTicketReviewList struct {
	List  []models.TicketReview `json:"list"`
	Total int64                 `json:"total"`
}
