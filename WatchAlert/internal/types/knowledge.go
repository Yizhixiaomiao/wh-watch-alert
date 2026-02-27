package types

import "watchAlert/internal/models"

// RequestKnowledgeCreate 创建知识请求
type RequestKnowledgeCreate struct {
	TenantId     string                 `json:"tenantId" binding:"required"`
	Title        string                 `json:"title" binding:"required"`
	Category     string                 `json:"category" binding:"required"`
	Tags         []string               `json:"tags"`
	Content      string                 `json:"content" binding:"required"`
	SourceTicket string                 `json:"sourceTicket"`
	AuthorId     string                 `json:"authorId"`
	Status       models.KnowledgeStatus `json:"status"`
}

// RequestKnowledgeUpdate 更新知识请求
type RequestKnowledgeUpdate struct {
	TenantId    string                 `json:"tenantId"`
	KnowledgeId string                 `json:"knowledgeId" binding:"required"`
	Title       string                 `json:"title"`
	Category    string                 `json:"category"`
	Tags        []string               `json:"tags"`
	Content     string                 `json:"content"`
	Status      models.KnowledgeStatus `json:"status"`
}

// RequestKnowledgeDelete 删除知识请求
type RequestKnowledgeDelete struct {
	TenantId    string `json:"tenantId"`
	KnowledgeId string `json:"knowledgeId" binding:"required"`
}

// RequestKnowledgeQuery 查询知识请求
type RequestKnowledgeQuery struct {
	TenantId     string                 `json:"tenantId" form:"tenantId"`
	KnowledgeId  string                 `json:"knowledgeId" form:"knowledgeId"`
	Title        string                 `json:"title" form:"title"`
	Category     string                 `json:"category" form:"category"`
	Tags         string                 `json:"tags" form:"tags"`
	SourceTicket string                 `json:"sourceTicket" form:"sourceTicket"`
	Status       models.KnowledgeStatus `json:"status" form:"status"`
	AuthorId     string                 `json:"authorId" form:"authorId"`
	Keyword      string                 `json:"keyword" form:"keyword"`
	Page         int                    `json:"page" form:"page"`
	Size         int                    `json:"size" form:"size"`
}

// RequestKnowledgeLike 点赞知识请求
type RequestKnowledgeLike struct {
	TenantId    string `json:"tenantId"`
	KnowledgeId string `json:"knowledgeId" binding:"required"`
	UserId      string `json:"userId" binding:"required"`
}

// RequestKnowledgeSaveToTicket 保存知识到工单请求
type RequestKnowledgeSaveToTicket struct {
	TenantId    string `json:"tenantId"`
	KnowledgeId string `json:"knowledgeId" binding:"required"`
	TicketId    string `json:"ticketId" binding:"required"`
	UserId      string `json:"userId"`
}

// RequestKnowledgeCategoryCreate 创建知识分类请求
type RequestKnowledgeCategoryCreate struct {
	TenantId     string `json:"tenantId"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"displayOrder"`
}

// RequestKnowledgeCategoryUpdate 更新知识分类请求
type RequestKnowledgeCategoryUpdate struct {
	TenantId     string `json:"tenantId"`
	CategoryId   string `json:"categoryId" binding:"required"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"displayOrder"`
	IsActive     *bool  `json:"isActive"`
}

// RequestKnowledgeCategoryDelete 删除知识分类请求
type RequestKnowledgeCategoryDelete struct {
	TenantId   string `json:"tenantId"`
	CategoryId string `json:"categoryId" binding:"required"`
}

// RequestKnowledgeCategoryQuery 查询知识分类请求
type RequestKnowledgeCategoryQuery struct {
	TenantId   string `json:"tenantId" form:"tenantId"`
	CategoryId string `json:"categoryId" form:"categoryId"`
	IsActive   *bool  `json:"isActive" form:"isActive"`
	Page       int    `json:"page" form:"page"`
	Size       int    `json:"size" form:"size"`
}

// ResponseKnowledgeList 知识列表响应
type ResponseKnowledgeList struct {
	List  []models.Knowledge `json:"list"`
	Total int64              `json:"total"`
}

// ResponseKnowledgeCategoryList 知识分类列表响应
type ResponseKnowledgeCategoryList struct {
	List  []models.KnowledgeCategory `json:"list"`
	Total int64                      `json:"total"`
}
