package types

import "watchAlert/internal/models"

// RequestWorkHoursStandardCreate 创建工时标准请求
type RequestWorkHoursStandardCreate struct {
	TenantId      string  `json:"tenantId" binding:"required"`
	Category      string  `json:"category" binding:"required"`
	SubCategory   string  `json:"subCategory" binding:"required"`
	Difficulty    string  `json:"difficulty" binding:"required"`
	StandardHours float64 `json:"standardHours" binding:"required,min=0"`
	Description   string  `json:"description"`
	CreatedBy     string  `json:"createdBy"`
}

// RequestWorkHoursStandardUpdate 更新工时标准请求
type RequestWorkHoursStandardUpdate struct {
	TenantId      string  `json:"tenantId"`
	StandardId    string  `json:"standardId" binding:"required"`
	Category      string  `json:"category"`
	SubCategory   string  `json:"subCategory"`
	Difficulty    string  `json:"difficulty"`
	StandardHours float64 `json:"standardHours"`
	Description   string  `json:"description"`
}

// RequestWorkHoursStandardDelete 删除工时标准请求
type RequestWorkHoursStandardDelete struct {
	TenantId   string `json:"tenantId"`
	StandardId string `json:"standardId" binding:"required"`
}

// RequestWorkHoursStandardQuery 查询工时标准请求
type RequestWorkHoursStandardQuery struct {
	TenantId    string `json:"tenantId" form:"tenantId"`
	StandardId  string `json:"standardId" form:"standardId"`
	Category    string `json:"category" form:"category"`
	SubCategory string `json:"subCategory" form:"subCategory"`
	Difficulty  string `json:"difficulty" form:"difficulty"`
	Page        int    `json:"page" form:"page"`
	Size        int    `json:"size" form:"size"`
}

// RequestWorkHoursCalculate 计算工时请求
type RequestWorkHoursCalculate struct {
	TenantId    string `json:"tenantId" binding:"required"`
	Category    string `json:"category" binding:"required"`
	SubCategory string `json:"subCategory" binding:"required"`
	Difficulty  string `json:"difficulty" binding:"required"`
}

// ResponseWorkHoursStandardList 工时标准列表响应
type ResponseWorkHoursStandardList struct {
	List  []models.WorkHoursStandard `json:"list"`
	Total int64                      `json:"total"`
}

// ResponseWorkHoursCalculate 工时计算响应
type ResponseWorkHoursCalculate struct {
	StandardHours float64 `json:"standardHours"`
	Category      string  `json:"category"`
	SubCategory   string  `json:"subCategory"`
	Difficulty    string  `json:"difficulty"`
}
