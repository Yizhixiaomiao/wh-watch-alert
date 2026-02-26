package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type assignmentRuleController struct{}

var AssignmentRuleController = new(assignmentRuleController)

/*
分配规则 API
/api/w8t/assignment-rule
*/
func (arc assignmentRuleController) API(gin *gin.RouterGroup) {
	// 需要审计日志的操作
	a := gin.Group("assignment-rule")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("create", AssignmentRuleController.CreateRule)
		a.POST("update", AssignmentRuleController.UpdateRule)
		a.POST("delete", AssignmentRuleController.DeleteRule)
		a.POST("auto-assign", AssignmentRuleController.AutoAssign)
	}

	// 查询操作
	b := gin.Group("assignment-rule")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("list", AssignmentRuleController.ListRules)
		b.GET("get", AssignmentRuleController.GetRule)
		b.POST("match", AssignmentRuleController.MatchRule)
	}
}

// CreateRule 创建规则
func (arc assignmentRuleController) CreateRule(ctx *gin.Context) {
	r := new(types.RequestAssignmentRuleCreate)
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
		return services.AssignmentRuleService.CreateRule(r)
	})
}

// UpdateRule 更新规则
func (arc assignmentRuleController) UpdateRule(ctx *gin.Context) {
	r := new(types.RequestAssignmentRuleUpdate)
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
		return services.AssignmentRuleService.UpdateRule(r)
	})
}

// DeleteRule 删除规则
func (arc assignmentRuleController) DeleteRule(ctx *gin.Context) {
	r := new(types.RequestAssignmentRuleDelete)
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
		return services.AssignmentRuleService.DeleteRule(r)
	})
}

// GetRule 获取规则详情
func (arc assignmentRuleController) GetRule(ctx *gin.Context) {
	r := new(types.RequestAssignmentRuleQuery)
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
		return services.AssignmentRuleService.GetRule(r)
	})
}

// ListRules 获取规则列表
func (arc assignmentRuleController) ListRules(ctx *gin.Context) {
	r := new(types.RequestAssignmentRuleQuery)
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
		return services.AssignmentRuleService.ListRules(r)
	})
}

// MatchRule 匹配规则
func (arc assignmentRuleController) MatchRule(ctx *gin.Context) {
	r := new(types.RequestAssignmentRuleMatch)
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
		return services.AssignmentRuleService.MatchRule(r)
	})
}

// AutoAssign 自动分配
func (arc assignmentRuleController) AutoAssign(ctx *gin.Context) {
	r := new(types.RequestAutoAssign)
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
		return services.AssignmentRuleService.AutoAssign(r)
	})
}
