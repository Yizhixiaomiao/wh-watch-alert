package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	AssignmentRuleRepo struct {
		entryRepo
	}

	InterAssignmentRuleRepo interface {
		// 规则操作
		CreateRule(rule models.AssignmentRule) error
		UpdateRule(rule models.AssignmentRule) error
		DeleteRule(tenantId, ruleId string) error
		GetRule(tenantId, ruleId string) (models.AssignmentRule, error)
		ListRules(tenantId string, ruleType models.AssignmentRuleType, enabled *bool, page, size int) ([]models.AssignmentRule, int64, error)

		// 规则匹配
		MatchRule(tenantId, alertType, dataSource, severity string) ([]models.AssignmentRule, error)
		GetEnabledRulesByType(tenantId string, ruleType models.AssignmentRuleType) ([]models.AssignmentRule, error)
	}
)

func newAssignmentRuleInterface(db *gorm.DB, g InterGormDBCli) InterAssignmentRuleRepo {
	return &AssignmentRuleRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

// CreateRule 创建规则
func (ar AssignmentRuleRepo) CreateRule(rule models.AssignmentRule) error {
	return ar.g.Create(&models.AssignmentRule{}, &rule)
}

// UpdateRule 更新规则
func (ar AssignmentRuleRepo) UpdateRule(rule models.AssignmentRule) error {
	return ar.g.Updates(Updates{
		Table:   &models.AssignmentRule{},
		Where:   map[string]interface{}{"tenant_id": rule.TenantId, "rule_id": rule.RuleId},
		Updates: rule,
	})
}

// DeleteRule 删除规则
func (ar AssignmentRuleRepo) DeleteRule(tenantId, ruleId string) error {
	return ar.g.Delete(Delete{
		Table: &models.AssignmentRule{},
		Where: map[string]interface{}{"tenant_id": tenantId, "rule_id": ruleId},
	})
}

// GetRule 获取规则详情
func (ar AssignmentRuleRepo) GetRule(tenantId, ruleId string) (models.AssignmentRule, error) {
	var rule models.AssignmentRule
	db := ar.db.Model(&models.AssignmentRule{})
	db.Where("tenant_id = ? AND rule_id = ?", tenantId, ruleId)
	err := db.First(&rule).Error
	if err != nil {
		return rule, err
	}
	return rule, nil
}

// ListRules 获取规则列表
func (ar AssignmentRuleRepo) ListRules(tenantId string, ruleType models.AssignmentRuleType, enabled *bool, page, size int) ([]models.AssignmentRule, int64, error) {
	var (
		rules []models.AssignmentRule
		count int64
	)

	db := ar.db.Model(&models.AssignmentRule{})
	db.Where("tenant_id = ?", tenantId)

	if ruleType != "" {
		db.Where("rule_type = ?", ruleType)
	}
	if enabled != nil {
		db.Where("enabled = ?", *enabled)
	}

	db.Count(&count)

	if page > 0 && size > 0 {
		db.Limit(size).Offset((page - 1) * size)
	}

	db.Order("priority ASC, created_at DESC")

	err := db.Find(&rules).Error
	if err != nil {
		return nil, 0, err
	}

	return rules, count, nil
}

// MatchRule 匹配规则
func (ar AssignmentRuleRepo) MatchRule(tenantId, alertType, dataSource, severity string) ([]models.AssignmentRule, error) {
	var rules []models.AssignmentRule

	db := ar.db.Model(&models.AssignmentRule{})
	db.Where("tenant_id = ? AND enabled = ?", tenantId, true)

	// 构建匹配条件
	where := "1=1"
	args := []interface{}{tenantId, true}

	if alertType != "" {
		where += " AND alert_type = ?"
		args = append(args, alertType)
	}
	if dataSource != "" {
		where += " AND data_source = ?"
		args = append(args, dataSource)
	}
	if severity != "" {
		where += " AND severity = ?"
		args = append(args, severity)
	}

	db.Where(where, args...)
	db.Order("priority ASC")

	err := db.Find(&rules).Error
	if err != nil {
		return nil, err
	}

	return rules, nil
}

// GetEnabledRulesByType 按类型获取启用规则
func (ar AssignmentRuleRepo) GetEnabledRulesByType(tenantId string, ruleType models.AssignmentRuleType) ([]models.AssignmentRule, error) {
	var rules []models.AssignmentRule

	db := ar.db.Model(&models.AssignmentRule{})
	db.Where("tenant_id = ? AND rule_type = ? AND enabled = ?", tenantId, ruleType, true)
	db.Order("priority ASC")

	err := db.Find(&rules).Error
	if err != nil {
		return nil, err
	}

	return rules, nil
}
