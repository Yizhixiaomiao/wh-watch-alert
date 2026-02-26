package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	WorkHoursRepo struct {
		entryRepo
	}

	InterWorkHoursRepo interface {
		// 工时标准操作
		CreateStandard(standard models.WorkHoursStandard) error
		UpdateStandard(standard models.WorkHoursStandard) error
		DeleteStandard(tenantId, standardId string) error
		GetStandard(tenantId, standardId string) (models.WorkHoursStandard, error)
		ListStandards(tenantId, category, subCategory, difficulty string, page, size int) ([]models.WorkHoursStandard, int64, error)
		GetStandardByCategory(tenantId, category, subCategory, difficulty string) (models.WorkHoursStandard, error)
	}
)

func newWorkHoursInterface(db *gorm.DB, g InterGormDBCli) InterWorkHoursRepo {
	return &WorkHoursRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

// CreateStandard 创建工时标准
func (wh WorkHoursRepo) CreateStandard(standard models.WorkHoursStandard) error {
	return wh.g.Create(&models.WorkHoursStandard{}, &standard)
}

// UpdateStandard 更新工时标准
func (wh WorkHoursRepo) UpdateStandard(standard models.WorkHoursStandard) error {
	return wh.g.Updates(Updates{
		Table:   &models.WorkHoursStandard{},
		Where:   map[string]interface{}{"tenant_id": standard.TenantId, "standard_id": standard.StandardId},
		Updates: standard,
	})
}

// DeleteStandard 删除工时标准
func (wh WorkHoursRepo) DeleteStandard(tenantId, standardId string) error {
	return wh.g.Delete(Delete{
		Table: &models.WorkHoursStandard{},
		Where: map[string]interface{}{"tenant_id": tenantId, "standard_id": standardId},
	})
}

// GetStandard 获取工时标准详情
func (wh WorkHoursRepo) GetStandard(tenantId, standardId string) (models.WorkHoursStandard, error) {
	var standard models.WorkHoursStandard
	db := wh.db.Model(&models.WorkHoursStandard{})
	db.Where("tenant_id = ? AND standard_id = ?", tenantId, standardId)
	err := db.First(&standard).Error
	if err != nil {
		return standard, err
	}
	return standard, nil
}

// ListStandards 获取工时标准列表
func (wh WorkHoursRepo) ListStandards(tenantId, category, subCategory, difficulty string, page, size int) ([]models.WorkHoursStandard, int64, error) {
	var (
		standards []models.WorkHoursStandard
		count     int64
	)

	db := wh.db.Model(&models.WorkHoursStandard{})
	db.Where("tenant_id = ?", tenantId)

	if category != "" {
		db.Where("category = ?", category)
	}
	if subCategory != "" {
		db.Where("sub_category = ?", subCategory)
	}
	if difficulty != "" {
		db.Where("difficulty = ?", difficulty)
	}

	db.Count(&count)

	if page > 0 && size > 0 {
		db.Limit(size).Offset((page - 1) * size)
	}

	db.Order("category, sub_category, difficulty")

	err := db.Find(&standards).Error
	if err != nil {
		return nil, 0, err
	}

	return standards, count, nil
}

// GetStandardByCategory 根据分类获取工时标准
func (wh WorkHoursRepo) GetStandardByCategory(tenantId, category, subCategory, difficulty string) (models.WorkHoursStandard, error) {
	var standard models.WorkHoursStandard
	db := wh.db.Model(&models.WorkHoursStandard{})
	db.Where("tenant_id = ? AND category = ? AND sub_category = ? AND difficulty = ?", tenantId, category, subCategory, difficulty)
	err := db.First(&standard).Error
	if err != nil {
		return standard, err
	}
	return standard, nil
}
