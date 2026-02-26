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

// RequestTicketReviewerCreate 创建评委请求
type RequestTicketReviewerCreate struct {
	TenantId   string `json:"tenantId" binding:"required"`
	ReviewerId string `json:"reviewerId" binding:"required"`
	UserName   string `json:"userName" binding:"required"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Department string `json:"department"`
	Specialty  string `json:"specialty"`
	CreatedBy  string `json:"createdBy"`
}

// RequestTicketReviewerUpdate 更新评委请求
type RequestTicketReviewerUpdate struct {
	TenantId   string `json:"tenantId"`
	ReviewerId string `json:"reviewerId" binding:"required"`
	UserName   string `json:"userName"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Department string `json:"department"`
	Specialty  string `json:"specialty"`
	IsActive   *bool  `json:"isActive"`
}

// RequestTicketReviewerQuery 查询评委请求
type RequestTicketReviewerQuery struct {
	TenantId   string `json:"tenantId" form:"tenantId"`
	ReviewerId string `json:"reviewerId" form:"reviewerId"`
	Department string `json:"department" form:"department"`
	Specialty  string `json:"specialty" form:"specialty"`
	IsActive   *bool  `json:"isActive" form:"isActive"`
	Page       int    `json:"page" form:"page"`
	Size       int    `json:"size" form:"size"`
}

// ResponseTicketReviewList 评审列表响应
type ResponseTicketReviewList struct {
	List  []models.TicketReview `json:"list"`
	Total int64                 `json:"total"`
}

// ResponseTicketReviewerList 评委列表响应
type ResponseTicketReviewerList struct {
	List  []models.TicketReviewer `json:"list"`
	Total int64                   `json:"total"`
}
