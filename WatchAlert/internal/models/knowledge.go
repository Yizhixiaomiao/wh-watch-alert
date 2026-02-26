package models

type KnowledgeStatus string

const (
	KnowledgeStatusDraft     KnowledgeStatus = "draft"     // 草稿
	KnowledgeStatusPublished KnowledgeStatus = "published" // 已发布
	KnowledgeStatusArchived  KnowledgeStatus = "archived"  // 已归档
)

// Knowledge 知识库表
type Knowledge struct {
	KnowledgeId    string          `json:"knowledgeId" gorm:"column:knowledge_id;primaryKey"`
	TenantId       string          `json:"tenantId" gorm:"column:tenant_id;index:idx_tenant_id"`
	Title          string          `json:"title" gorm:"column:title;index:idx_title"`
	Category       string          `json:"category" gorm:"column:category;index:idx_category"`
	Tags           []string        `json:"tags" gorm:"column:tags;serializer:json"`
	Content        string          `json:"content" gorm:"column:content;type:text"`
	ContentText    string          `json:"contentText" gorm:"column:content_text;type:text"` // 纯文本内容，用于搜索
	SourceTicket   string          `json:"sourceTicket" gorm:"column:source_ticket;index:idx_source_ticket"`
	RelatedTickets []string        `json:"relatedTickets" gorm:"column:related_tickets;serializer:json"` // 关联的工单ID列表
	AuthorId       string          `json:"authorId" gorm:"column:author_id;index:idx_author_id"`
	Status         KnowledgeStatus `json:"status" gorm:"column:status;index:idx_status"`
	ViewCount      int64           `json:"viewCount" gorm:"column:view_count;default:0"`
	LikeCount      int64           `json:"likeCount" gorm:"column:like_count;default:0"`
	UseCount       int64           `json:"useCount" gorm:"column:use_count;default:0"`
	CreatedAt      int64           `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt      int64           `json:"updatedAt" gorm:"column:updated_at"`
}

// TableName 指定表名
func (Knowledge) TableName() string {
	return "knowledge"
}

// KnowledgeLike 知识点赞表
type KnowledgeLike struct {
	KnowledgeId string `json:"knowledgeId" gorm:"column:knowledge_id;primaryKey"`
	UserId      string `json:"userId" gorm:"column:user_id;primaryKey"`
	CreatedAt   int64  `json:"createdAt" gorm:"column:created_at"`
}

// TableName 指定表名
func (KnowledgeLike) TableName() string {
	return "knowledge_like"
}

// KnowledgeCategory 知识分类表
type KnowledgeCategory struct {
	CategoryId   string `json:"categoryId" gorm:"column:category_id;primaryKey"`
	TenantId     string `json:"tenantId" gorm:"column:tenant_id;index:idx_tenant_id"`
	Name         string `json:"name" gorm:"column:name;index:idx_name"`
	Description  string `json:"description" gorm:"column:description;type:text"`
	DisplayOrder int    `json:"displayOrder" gorm:"column:display_order;default:0"`
	IsActive     bool   `json:"isActive" gorm:"column:is_active;default:true"`
	CreatedAt    int64  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    int64  `json:"updatedAt" gorm:"column:updated_at"`
}

// TableName 指定表名
func (KnowledgeCategory) TableName() string {
	return "knowledge_category"
}
