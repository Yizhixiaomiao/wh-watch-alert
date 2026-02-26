package models

type TicketReviewStatus string

const (
	ReviewStatusPending   TicketReviewStatus = "pending"   // 待评审
	ReviewStatusCompleted TicketReviewStatus = "completed" // 已完成
)

// TicketReview 工单评审表
type TicketReview struct {
	ReviewId    string             `json:"reviewId" gorm:"column:review_id;primaryKey"`
	TenantId    string             `json:"tenantId" gorm:"column:tenant_id;index:idx_tenant_id"`
	TicketId    string             `json:"ticketId" gorm:"column:ticket_id;index:idx_ticket_id"`
	ReviewerId  string             `json:"reviewerId" gorm:"column:reviewer_id;index:idx_reviewer_id"`
	Rating      int                `json:"rating" gorm:"column:rating"`        // 评分 1-5
	WorkHours   float64            `json:"workHours" gorm:"column:work_hours"` // 工时（小时）
	Comment     string             `json:"comment" gorm:"column:comment;type:text"`
	Status      TicketReviewStatus `json:"status" gorm:"column:status"`
	CreatedBy   string             `json:"createdBy" gorm:"column:created_by"`
	CreatedAt   int64              `json:"createdAt" gorm:"column:created_at"`
	CompletedAt int64              `json:"completedAt" gorm:"column:completed_at"`
}

// TableName 指定表名
func (TicketReview) TableName() string {
	return "ticket_review"
}

// TicketReviewer 评委表
type TicketReviewer struct {
	TenantId   string `json:"tenantId" gorm:"column:tenant_id;index:idx_tenant_id"`
	ReviewerId string `json:"reviewerId" gorm:"column:reviewer_id;primaryKey"`
	UserName   string `json:"userName" gorm:"column:user_name"`
	Email      string `json:"email" gorm:"column:email"`
	Phone      string `json:"phone" gorm:"column:phone"`
	Department string `json:"department" gorm:"column:department"`
	Specialty  string `json:"specialty" gorm:"column:specialty"` // 专业领域
	IsActive   bool   `json:"isActive" gorm:"column:is_active;default:true"`
	CreatedAt  int64  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt  int64  `json:"updatedAt" gorm:"column:updated_at"`
}

// TableName 指定表名
func (TicketReviewer) TableName() string {
	return "ticket_reviewer"
}
