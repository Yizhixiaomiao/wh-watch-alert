package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	KnowledgeRepo struct {
		entryRepo
	}

	InterKnowledgeRepo interface {
		// 知识操作
		CreateKnowledge(knowledge models.Knowledge) error
		UpdateKnowledge(knowledge models.Knowledge) error
		DeleteKnowledge(tenantId, knowledgeId string) error
		GetKnowledge(tenantId, knowledgeId string) (models.Knowledge, error)
		ListKnowledges(tenantId, title, category, sourceTicket, authorId, keyword string, status models.KnowledgeStatus, page, size int) ([]models.Knowledge, int64, error)
		IncrementViewCount(knowledgeId string) error
		IncrementLikeCount(knowledgeId string) error
		DecrementLikeCount(knowledgeId string) error
		IncrementUseCount(knowledgeId string) error

		// 点赞操作
		CreateLike(like models.KnowledgeLike) error
		DeleteLike(knowledgeId, userId string) error
		CheckLike(knowledgeId, userId string) (bool, error)

		// 分类操作
		CreateCategory(category models.KnowledgeCategory) error
		UpdateCategory(category models.KnowledgeCategory) error
		DeleteCategory(tenantId, categoryId string) error
		GetCategory(tenantId, categoryId string) (models.KnowledgeCategory, error)
		ListCategories(tenantId string, isActive *bool, page, size int) ([]models.KnowledgeCategory, int64, error)
	}
)

func newKnowledgeInterface(db *gorm.DB, g InterGormDBCli) InterKnowledgeRepo {
	return &KnowledgeRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

// CreateKnowledge 创建知识
func (kr KnowledgeRepo) CreateKnowledge(knowledge models.Knowledge) error {
	return kr.g.Create(&models.Knowledge{}, &knowledge)
}

// UpdateKnowledge 更新知识
func (kr KnowledgeRepo) UpdateKnowledge(knowledge models.Knowledge) error {
	return kr.g.Updates(Updates{
		Table:   &models.Knowledge{},
		Where:   map[string]interface{}{"tenant_id": knowledge.TenantId, "knowledge_id": knowledge.KnowledgeId},
		Updates: knowledge,
	})
}

// DeleteKnowledge 删除知识
func (kr KnowledgeRepo) DeleteKnowledge(tenantId, knowledgeId string) error {
	return kr.g.Delete(Delete{
		Table: &models.Knowledge{},
		Where: map[string]interface{}{"tenant_id": tenantId, "knowledge_id": knowledgeId},
	})
}

// GetKnowledge 获取知识详情
func (kr KnowledgeRepo) GetKnowledge(tenantId, knowledgeId string) (models.Knowledge, error) {
	var knowledge models.Knowledge
	db := kr.db.Model(&models.Knowledge{})
	db.Where("tenant_id = ? AND knowledge_id = ?", tenantId, knowledgeId)
	err := db.First(&knowledge).Error
	if err != nil {
		return knowledge, err
	}
	return knowledge, nil
}

// ListKnowledges 获取知识列表
func (kr KnowledgeRepo) ListKnowledges(tenantId, title, category, sourceTicket, authorId, keyword string, status models.KnowledgeStatus, page, size int) ([]models.Knowledge, int64, error) {
	var (
		knowledges []models.Knowledge
		count      int64
	)

	db := kr.db.Model(&models.Knowledge{})
	db.Where("tenant_id = ?", tenantId)

	if title != "" {
		db.Where("title LIKE ?", "%"+title+"%")
	}
	if category != "" {
		db.Where("category = ?", category)
	}
	if sourceTicket != "" {
		db.Where("source_ticket = ?", sourceTicket)
	}
	if authorId != "" {
		db.Where("author_id = ?", authorId)
	}
	if status != "" {
		db.Where("status = ?", status)
	}
	if keyword != "" {
		db.Where("title LIKE ? OR content_text LIKE ? OR JSON_CONTAINS(tags, ?)", "%"+keyword+"%", "%"+keyword+"%", `["`+keyword+`"]`)
	}

	db.Count(&count)

	if page > 0 && size > 0 {
		db.Limit(size).Offset((page - 1) * size)
	}

	db.Order("updated_at DESC")

	err := db.Find(&knowledges).Error
	if err != nil {
		return nil, 0, err
	}

	return knowledges, count, nil
}

// IncrementViewCount 增加浏览次数
func (kr KnowledgeRepo) IncrementViewCount(knowledgeId string) error {
	return kr.db.Model(&models.Knowledge{}).
		Where("knowledge_id = ?", knowledgeId).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// IncrementLikeCount 增加点赞次数
func (kr KnowledgeRepo) IncrementLikeCount(knowledgeId string) error {
	return kr.db.Model(&models.Knowledge{}).
		Where("knowledge_id = ?", knowledgeId).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

// DecrementLikeCount 减少点赞次数
func (kr KnowledgeRepo) DecrementLikeCount(knowledgeId string) error {
	return kr.db.Model(&models.Knowledge{}).
		Where("knowledge_id = ?", knowledgeId).
		UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error
}

// IncrementUseCount 增加使用次数
func (kr KnowledgeRepo) IncrementUseCount(knowledgeId string) error {
	return kr.db.Model(&models.Knowledge{}).
		Where("knowledge_id = ?", knowledgeId).
		UpdateColumn("use_count", gorm.Expr("use_count + 1")).Error
}

// CreateLike 创建点赞
func (kr KnowledgeRepo) CreateLike(like models.KnowledgeLike) error {
	return kr.g.Create(&models.KnowledgeLike{}, &like)
}

// DeleteLike 删除点赞
func (kr KnowledgeRepo) DeleteLike(knowledgeId, userId string) error {
	return kr.g.Delete(Delete{
		Table: &models.KnowledgeLike{},
		Where: map[string]interface{}{"knowledge_id": knowledgeId, "user_id": userId},
	})
}

// CheckLike 检查是否已点赞
func (kr KnowledgeRepo) CheckLike(knowledgeId, userId string) (bool, error) {
	var count int64
	err := kr.db.Model(&models.KnowledgeLike{}).
		Where("knowledge_id = ? AND user_id = ?", knowledgeId, userId).
		Count(&count).Error
	return count > 0, err
}

// CreateCategory 创建分类
func (kr KnowledgeRepo) CreateCategory(category models.KnowledgeCategory) error {
	return kr.g.Create(&models.KnowledgeCategory{}, &category)
}

// UpdateCategory 更新分类
func (kr KnowledgeRepo) UpdateCategory(category models.KnowledgeCategory) error {
	return kr.g.Updates(Updates{
		Table:   &models.KnowledgeCategory{},
		Where:   map[string]interface{}{"tenant_id": category.TenantId, "category_id": category.CategoryId},
		Updates: category,
	})
}

// DeleteCategory 删除分类
func (kr KnowledgeRepo) DeleteCategory(tenantId, categoryId string) error {
	return kr.g.Delete(Delete{
		Table: &models.KnowledgeCategory{},
		Where: map[string]interface{}{"tenant_id": tenantId, "category_id": categoryId},
	})
}

// GetCategory 获取分类详情
func (kr KnowledgeRepo) GetCategory(tenantId, categoryId string) (models.KnowledgeCategory, error) {
	var category models.KnowledgeCategory
	db := kr.db.Model(&models.KnowledgeCategory{})
	db.Where("tenant_id = ? AND category_id = ?", tenantId, categoryId)
	err := db.First(&category).Error
	if err != nil {
		return category, err
	}
	return category, nil
}

// ListCategories 获取分类列表
func (kr KnowledgeRepo) ListCategories(tenantId string, isActive *bool, page, size int) ([]models.KnowledgeCategory, int64, error) {
	var (
		categories []models.KnowledgeCategory
		count      int64
	)

	db := kr.db.Model(&models.KnowledgeCategory{})
	db.Where("tenant_id = ?", tenantId)

	if isActive != nil {
		db.Where("is_active = ?", *isActive)
	}

	db.Count(&count)

	if page > 0 && size > 0 {
		db.Limit(size).Offset((page - 1) * size)
	}

	db.Order("display_order ASC, created_at DESC")

	err := db.Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}
