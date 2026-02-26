package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	AlertTicketRuleRepo struct {
		entryRepo
	}

	InterAlertTicketRuleRepo interface {
		Create(rule models.AlertTicketRule) error
		Update(rule models.AlertTicketRule) error
		Delete(tenantId, id string) error
		Get(tenantId, id string) (models.AlertTicketRule, error)
		List(tenantId string, page, size int) ([]models.AlertTicketRule, int64, error)
		CreateHistory(history models.AlertTicketRuleHistory) error
		ListHistory(tenantId string, page, size int) ([]models.AlertTicketRuleHistory, int64, error)
	}
)

func newAlertTicketRuleInterface(db *gorm.DB, g InterGormDBCli) InterAlertTicketRuleRepo {
	return &AlertTicketRuleRepo{
		entryRepo: entryRepo{
			g:  g,
			db: db,
		},
	}
}

// Create 创建告警转工单规则
func (r *AlertTicketRuleRepo) Create(rule models.AlertTicketRule) error {
	return r.g.Create(&models.AlertTicketRule{}, &rule)
}

// Update 更新告警转工单规则
func (r *AlertTicketRuleRepo) Update(rule models.AlertTicketRule) error {
	return r.g.Updates(Updates{
		Table: &models.AlertTicketRule{},
		Where: map[string]interface{}{
			"tenant_id": rule.TenantId,
			"id":        rule.Id,
		},
		Updates: rule,
	})
}

// Delete 删除告警转工单规则
func (r *AlertTicketRuleRepo) Delete(tenantId, id string) error {
	return r.g.Delete(Delete{
		Table: &models.AlertTicketRule{},
		Where: map[string]interface{}{
			"tenant_id": tenantId,
			"id":        id,
		},
	})
}

// Get 获取告警转工单规则
func (r *AlertTicketRuleRepo) Get(tenantId, id string) (models.AlertTicketRule, error) {
	var rule models.AlertTicketRule
	err := r.db.Where("tenant_id = ? AND id = ?", tenantId, id).First(&rule).Error
	return rule, err
}

// List 获取告警转工单规则列表
func (r *AlertTicketRuleRepo) List(tenantId string, page, size int) ([]models.AlertTicketRule, int64, error) {
	var rules []models.AlertTicketRule
	var total int64

	db := r.db.Model(&models.AlertTicketRule{}).Where("tenant_id = ?", tenantId)

	// 获取总数
	db.Count(&total)

	// 分页查询
	offset := (page - 1) * size
	err := db.Offset(offset).Limit(size).Find(&rules).Error

	return rules, total, err
}

// CreateHistory 创建历史记录
func (r *AlertTicketRuleRepo) CreateHistory(history models.AlertTicketRuleHistory) error {
	return r.g.Create(&models.AlertTicketRuleHistory{}, &history)
}

// ListHistory 获取历史记录列表
func (r *AlertTicketRuleRepo) ListHistory(tenantId string, page, size int) ([]models.AlertTicketRuleHistory, int64, error) {
	var histories []models.AlertTicketRuleHistory
	var total int64

	db := r.db.Model(&models.AlertTicketRuleHistory{}).Where("tenant_id = ?", tenantId)

	// 获取总数
	db.Count(&total)

	// 分页查询
	offset := (page - 1) * size
	err := db.Offset(offset).Limit(size).Order("created_at DESC").Find(&histories).Error

	return histories, total, err
}
