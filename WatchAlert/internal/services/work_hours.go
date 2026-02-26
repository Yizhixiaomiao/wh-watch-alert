package services

import (
	"fmt"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

type workHoursService struct {
	ctx *ctx.Context
}

type InterWorkHoursService interface {
	// 工时标准操作
	CreateStandard(req interface{}) (interface{}, interface{})
	UpdateStandard(req interface{}) (interface{}, interface{})
	DeleteStandard(req interface{}) (interface{}, interface{})
	GetStandard(req interface{}) (interface{}, interface{})
	ListStandards(req interface{}) (interface{}, interface{})
	CalculateHours(req interface{}) (interface{}, interface{})
}

func newInterWorkHoursService(ctx *ctx.Context) InterWorkHoursService {
	return &workHoursService{ctx}
}

// CreateStandard 创建工时标准
func (s workHoursService) CreateStandard(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestWorkHoursStandardCreate)

	standard := models.WorkHoursStandard{
		StandardId:    "whs-" + tools.RandId(),
		TenantId:      r.TenantId,
		Category:      r.Category,
		SubCategory:   r.SubCategory,
		Difficulty:    r.Difficulty,
		StandardHours: r.StandardHours,
		Description:   r.Description,
		CreatedBy:     r.CreatedBy,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}

	err := s.ctx.DB.WorkHours().CreateStandard(standard)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// UpdateStandard 更新工时标准
func (s workHoursService) UpdateStandard(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestWorkHoursStandardUpdate)

	standard, err := s.ctx.DB.WorkHours().GetStandard(r.TenantId, r.StandardId)
	if err != nil {
		return nil, fmt.Errorf("工时标准不存在")
	}

	if r.Category != "" {
		standard.Category = r.Category
	}
	if r.SubCategory != "" {
		standard.SubCategory = r.SubCategory
	}
	if r.Difficulty != "" {
		standard.Difficulty = r.Difficulty
	}
	if r.StandardHours > 0 {
		standard.StandardHours = r.StandardHours
	}
	if r.Description != "" {
		standard.Description = r.Description
	}
	standard.UpdatedAt = time.Now().Unix()

	err = s.ctx.DB.WorkHours().UpdateStandard(standard)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// DeleteStandard 删除工时标准
func (s workHoursService) DeleteStandard(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestWorkHoursStandardDelete)

	err := s.ctx.DB.WorkHours().DeleteStandard(r.TenantId, r.StandardId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetStandard 获取工时标准详情
func (s workHoursService) GetStandard(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestWorkHoursStandardQuery)

	standard, err := s.ctx.DB.WorkHours().GetStandard(r.TenantId, r.StandardId)
	if err != nil {
		return nil, err
	}

	return standard, nil
}

// ListStandards 获取工时标准列表
func (s workHoursService) ListStandards(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestWorkHoursStandardQuery)

	standards, total, err := s.ctx.DB.WorkHours().ListStandards(r.TenantId, r.Category, r.SubCategory, r.Difficulty, r.Page, r.Size)
	if err != nil {
		return nil, err
	}

	return types.ResponseWorkHoursStandardList{
		List:  standards,
		Total: total,
	}, nil
}

// CalculateHours 计算工时
func (s workHoursService) CalculateHours(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestWorkHoursCalculate)

	standard, err := s.ctx.DB.WorkHours().GetStandardByCategory(r.TenantId, r.Category, r.SubCategory, r.Difficulty)
	if err != nil {
		return nil, fmt.Errorf("未找到匹配的工时标准")
	}

	return types.ResponseWorkHoursCalculate{
		StandardHours: standard.StandardHours,
		Category:      standard.Category,
		SubCategory:   standard.SubCategory,
		Difficulty:    standard.Difficulty,
	}, nil
}
