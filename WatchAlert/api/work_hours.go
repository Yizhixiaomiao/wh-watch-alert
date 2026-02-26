package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type workHoursController struct{}

var WorkHoursController = new(workHoursController)

/*
工时管理 API
/api/w8t/work-hours
*/
func (whc workHoursController) API(gin *gin.RouterGroup) {
	// 需要审计日志的操作
	a := gin.Group("work-hours")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("standard/create", WorkHoursController.CreateStandard)
		a.POST("standard/update", WorkHoursController.UpdateStandard)
		a.POST("standard/delete", WorkHoursController.DeleteStandard)
		a.POST("calculate", WorkHoursController.CalculateHours)
	}

	// 查询操作
	b := gin.Group("work-hours")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("standard/list", WorkHoursController.ListStandards)
		b.GET("standard/get", WorkHoursController.GetStandard)
	}
}

// CreateStandard 创建工时标准
func (whc workHoursController) CreateStandard(ctx *gin.Context) {
	r := new(types.RequestWorkHoursStandardCreate)
	BindJson(ctx, r)

	tid, exists := ctx.Get("TenantID")
	if !exists {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("租户ID不存在")
		})
		return
	}
	r.TenantId = tid.(string)

	uid, exists := ctx.Get("UserID")
	if !exists {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("用户ID不存在")
		})
		return
	}
	r.CreatedBy = uid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.WorkHoursService.CreateStandard(r)
	})
}

// UpdateStandard 更新工时标准
func (whc workHoursController) UpdateStandard(ctx *gin.Context) {
	r := new(types.RequestWorkHoursStandardUpdate)
	BindJson(ctx, r)

	tid, exists := ctx.Get("TenantID")
	if !exists {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("租户ID不存在")
		})
		return
	}
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.WorkHoursService.UpdateStandard(r)
	})
}

// DeleteStandard 删除工时标准
func (whc workHoursController) DeleteStandard(ctx *gin.Context) {
	r := new(types.RequestWorkHoursStandardDelete)
	BindJson(ctx, r)

	tid, exists := ctx.Get("TenantID")
	if !exists {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("租户ID不存在")
		})
		return
	}
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.WorkHoursService.DeleteStandard(r)
	})
}

// GetStandard 获取工时标准详情
func (whc workHoursController) GetStandard(ctx *gin.Context) {
	r := new(types.RequestWorkHoursStandardQuery)
	BindQuery(ctx, r)

	tid, exists := ctx.Get("TenantID")
	if !exists {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("租户ID不存在")
		})
		return
	}
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.WorkHoursService.GetStandard(r)
	})
}

// ListStandards 获取工时标准列表
func (whc workHoursController) ListStandards(ctx *gin.Context) {
	r := new(types.RequestWorkHoursStandardQuery)
	BindQuery(ctx, r)

	tid, exists := ctx.Get("TenantID")
	if !exists {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("租户ID不存在")
		})
		return
	}
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.WorkHoursService.ListStandards(r)
	})
}

// CalculateHours 计算工时
func (whc workHoursController) CalculateHours(ctx *gin.Context) {
	r := new(types.RequestWorkHoursCalculate)
	BindJson(ctx, r)

	tid, exists := ctx.Get("TenantID")
	if !exists {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("租户ID不存在")
		})
		return
	}
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.WorkHoursService.CalculateHours(r)
	})
}
