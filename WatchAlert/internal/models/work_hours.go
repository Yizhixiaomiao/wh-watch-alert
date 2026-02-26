package models

// WorkHoursStandard 工时标准表
type WorkHoursStandard struct {
	StandardId    string  `json:"standardId" gorm:"column:standard_id;primaryKey"`
	TenantId      string  `json:"tenantId" gorm:"column:tenant_id;index:idx_tenant_id"`
	Category      string  `json:"category" gorm:"column:category;index:idx_category"`
	SubCategory   string  `json:"subCategory" gorm:"column:sub_category;index:idx_sub_category"`
	Difficulty    string  `json:"difficulty" gorm:"column:difficulty;index:idx_difficulty"`
	StandardHours float64 `json:"standardHours" gorm:"column:standard_hours"`
	Description   string  `json:"description" gorm:"column:description;type:text"`
	CreatedBy     string  `json:"createdBy" gorm:"column:created_by"`
	CreatedAt     int64   `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt     int64   `json:"updatedAt" gorm:"column:updated_at"`
}

// TableName 指定表名
func (WorkHoursStandard) TableName() string {
	return "work_hours_standard"
}
