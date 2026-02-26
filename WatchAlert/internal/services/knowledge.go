package services

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

type knowledgeService struct {
	ctx *ctx.Context
}

type InterKnowledgeService interface {
	// 知识操作
	CreateKnowledge(req interface{}) (interface{}, interface{})
	UpdateKnowledge(req interface{}) (interface{}, interface{})
	DeleteKnowledge(req interface{}) (interface{}, interface{})
	GetKnowledge(req interface{}) (interface{}, interface{})
	ListKnowledges(req interface{}) (interface{}, interface{})
	LikeKnowledge(req interface{}) (interface{}, interface{})
	SaveToTicket(req interface{}) (interface{}, interface{})

	// 分类操作
	CreateCategory(req interface{}) (interface{}, interface{})
	UpdateCategory(req interface{}) (interface{}, interface{})
	DeleteCategory(req interface{}) (interface{}, interface{})
	GetCategory(req interface{}) (interface{}, interface{})
	ListCategories(req interface{}) (interface{}, interface{})
}

func newInterKnowledgeService(ctx *ctx.Context) InterKnowledgeService {
	return &knowledgeService{ctx}
}

// htmlToPlainText 将HTML内容转换为纯文本
func htmlToPlainText(html string) string {
	// 移除HTML标签
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(html, "")

	// 移除多余的空白字符
	text = strings.Join(strings.Fields(text), " ")

	return text
}

// CreateKnowledge 创建知识
func (s knowledgeService) CreateKnowledge(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestKnowledgeCreate)

	knowledge := models.Knowledge{
		KnowledgeId:    "kn-" + tools.RandId(),
		TenantId:       r.TenantId,
		Title:          r.Title,
		Category:       r.Category,
		Tags:           r.Tags,
		Content:        r.Content,
		ContentText:    htmlToPlainText(r.Content),
		SourceTicket:   r.SourceTicket,
		RelatedTickets: []string{},
		AuthorId:       r.AuthorId,
		Status:         models.KnowledgeStatusDraft,
		ViewCount:      0,
		LikeCount:      0,
		UseCount:       0,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
	}

	if r.Status != "" {
		knowledge.Status = r.Status
	}

	// 如果有关联的工单，添加到RelatedTickets
	if r.SourceTicket != "" {
		knowledge.RelatedTickets = []string{r.SourceTicket}
	}

	err := s.ctx.DB.Knowledge().CreateKnowledge(knowledge)
	if err != nil {
		return nil, err
	}

	// 如果有来源工单，更新工单的knowledgeId
	if r.SourceTicket != "" {
		ticket, err := s.ctx.DB.Ticket().Get(r.TenantId, r.SourceTicket)
		if err == nil {
			ticket.KnowledgeId = knowledge.KnowledgeId
			ticket.UpdatedAt = time.Now().Unix()
			s.ctx.DB.Ticket().Update(ticket)
		}
	}

	return knowledge.KnowledgeId, nil
}

// UpdateKnowledge 更新知识
func (s knowledgeService) UpdateKnowledge(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestKnowledgeUpdate)

	knowledge, err := s.ctx.DB.Knowledge().GetKnowledge(r.TenantId, r.KnowledgeId)
	if err != nil {
		return nil, fmt.Errorf("知识不存在")
	}

	if r.Title != "" {
		knowledge.Title = r.Title
	}
	if r.Category != "" {
		knowledge.Category = r.Category
	}
	if r.Tags != nil {
		knowledge.Tags = r.Tags
	}
	if r.Content != "" {
		knowledge.Content = r.Content
		knowledge.ContentText = htmlToPlainText(r.Content)
	}
	if r.Status != "" {
		knowledge.Status = r.Status
	}
	knowledge.UpdatedAt = time.Now().Unix()

	err = s.ctx.DB.Knowledge().UpdateKnowledge(knowledge)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// DeleteKnowledge 删除知识
func (s knowledgeService) DeleteKnowledge(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestKnowledgeDelete)

	err := s.ctx.DB.Knowledge().DeleteKnowledge(r.TenantId, r.KnowledgeId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetKnowledge 获取知识详情
func (s knowledgeService) GetKnowledge(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestKnowledgeQuery)

	knowledge, err := s.ctx.DB.Knowledge().GetKnowledge(r.TenantId, r.KnowledgeId)
	if err != nil {
		return nil, err
	}

	// 增加浏览次数
	s.ctx.DB.Knowledge().IncrementViewCount(r.KnowledgeId)

	return knowledge, nil
}

// ListKnowledges 获取知识列表
func (s knowledgeService) ListKnowledges(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestKnowledgeQuery)

	knowledges, total, err := s.ctx.DB.Knowledge().ListKnowledges(
		r.TenantId,
		r.Title,
		r.Category,
		r.SourceTicket,
		r.AuthorId,
		r.Keyword,
		r.Status,
		r.Page,
		r.Size,
	)
	if err != nil {
		return nil, err
	}

	return types.ResponseKnowledgeList{
		List:  knowledges,
		Total: total,
	}, nil
}

// LikeKnowledge 点赞知识
func (s knowledgeService) LikeKnowledge(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestKnowledgeLike)

	// 检查是否已点赞
	hasLiked, err := s.ctx.DB.Knowledge().CheckLike(r.KnowledgeId, r.UserId)
	if err != nil {
		return nil, err
	}

	if hasLiked {
		// 取消点赞
		err = s.ctx.DB.Knowledge().DeleteLike(r.KnowledgeId, r.UserId)
		if err != nil {
			return nil, err
		}
		s.ctx.DB.Knowledge().DecrementLikeCount(r.KnowledgeId)
		return map[string]bool{"liked": false}, nil
	} else {
		// 添加点赞
		like := models.KnowledgeLike{
			KnowledgeId: r.KnowledgeId,
			UserId:      r.UserId,
			CreatedAt:   time.Now().Unix(),
		}
		err = s.ctx.DB.Knowledge().CreateLike(like)
		if err != nil {
			return nil, err
		}
		s.ctx.DB.Knowledge().IncrementLikeCount(r.KnowledgeId)
		return map[string]bool{"liked": true}, nil
	}
}

// SaveToTicket 保存知识到工单
func (s knowledgeService) SaveToTicket(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestKnowledgeSaveToTicket)

	// 获取知识
	knowledge, err := s.ctx.DB.Knowledge().GetKnowledge(r.TenantId, r.KnowledgeId)
	if err != nil {
		return nil, fmt.Errorf("知识不存在")
	}

	// 获取工单
	ticket, err := s.ctx.DB.Ticket().Get(r.TenantId, r.TicketId)
	if err != nil {
		return nil, fmt.Errorf("工单不存在")
	}

	// 将知识内容添加到工单描述
	newDescription := ticket.Description
	if newDescription != "" {
		newDescription += "\n\n---\n\n参考知识：\n"
	}
	newDescription += fmt.Sprintf("【%s】\n%s\n\n%s", knowledge.Title, knowledge.Category, knowledge.Content)

	// 更新工单
	ticket.Description = newDescription
	ticket.UpdatedAt = time.Now().Unix()

	err = s.ctx.DB.Ticket().Update(ticket)
	if err != nil {
		return nil, err
	}

	// 增加知识使用次数
	s.ctx.DB.Knowledge().IncrementUseCount(r.KnowledgeId)

	return nil, nil
}

// CreateCategory 创建分类
func (s knowledgeService) CreateCategory(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestKnowledgeCategoryCreate)

	category := models.KnowledgeCategory{
		CategoryId:   "kc-" + tools.RandId(),
		TenantId:     r.TenantId,
		Name:         r.Name,
		Description:  r.Description,
		DisplayOrder: r.DisplayOrder,
		IsActive:     true,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	err := s.ctx.DB.Knowledge().CreateCategory(category)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// UpdateCategory 更新分类
func (s knowledgeService) UpdateCategory(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestKnowledgeCategoryUpdate)

	category, err := s.ctx.DB.Knowledge().GetCategory(r.TenantId, r.CategoryId)
	if err != nil {
		return nil, fmt.Errorf("分类不存在")
	}

	if r.Name != "" {
		category.Name = r.Name
	}
	if r.Description != "" {
		category.Description = r.Description
	}
	if r.DisplayOrder > 0 {
		category.DisplayOrder = r.DisplayOrder
	}
	if r.IsActive != nil {
		category.IsActive = *r.IsActive
	}
	category.UpdatedAt = time.Now().Unix()

	err = s.ctx.DB.Knowledge().UpdateCategory(category)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// DeleteCategory 删除分类
func (s knowledgeService) DeleteCategory(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestKnowledgeCategoryDelete)

	err := s.ctx.DB.Knowledge().DeleteCategory(r.TenantId, r.CategoryId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetCategory 获取分类详情
func (s knowledgeService) GetCategory(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestKnowledgeCategoryQuery)

	category, err := s.ctx.DB.Knowledge().GetCategory(r.TenantId, r.CategoryId)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// ListCategories 获取分类列表
func (s knowledgeService) ListCategories(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestKnowledgeCategoryQuery)

	categories, total, err := s.ctx.DB.Knowledge().ListCategories(r.TenantId, r.IsActive, r.Page, r.Size)
	if err != nil {
		return nil, err
	}

	return types.ResponseKnowledgeCategoryList{
		List:  categories,
		Total: total,
	}, nil
}
