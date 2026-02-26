package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type alertTicketRuleController struct{}

var AlertTicketRuleController = new(alertTicketRuleController)

/*
告警转工单规则管理 API
/api/w8t/alertTicketRule
*/
func (atr alertTicketRuleController) API(gin *gin.RouterGroup) {
	a := gin.Group("alertTicketRule")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("create", AlertTicketRuleController.Create)
		a.POST("update", AlertTicketRuleController.Update)
		a.POST("delete", AlertTicketRuleController.Delete)
		a.POST("test", AlertTicketRuleController.Test)
	}

	b := gin.Group("alertTicketRule")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("list", AlertTicketRuleController.List)
		b.GET("get", AlertTicketRuleController.Get)
		b.GET("history", AlertTicketRuleController.History)
		b.GET("stats", AlertTicketRuleController.Stats)
	}
}

// Create 创建告警转工单规则
func (atr alertTicketRuleController) Create(ctx *gin.Context) {
	r := new(types.RequestAlertTicketRuleCreate)
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
		return services.AlertTicketService.CreateAlertTicketRule(r)
	})
}

// Update 更新告警转工单规则
func (atr alertTicketRuleController) Update(ctx *gin.Context) {
	r := new(types.RequestAlertTicketRuleUpdate)
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
		return services.AlertTicketService.UpdateAlertTicketRule(r)
	})
}

// Delete 删除告警转工单规则
func (atr alertTicketRuleController) Delete(ctx *gin.Context) {
	r := new(types.RequestAlertTicketRuleDelete)
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
		return services.AlertTicketService.DeleteAlertTicketRule(r)
	})
}

// Get 获取告警转工单规则
func (atr alertTicketRuleController) Get(ctx *gin.Context) {
	r := new(types.RequestAlertTicketRuleQuery)
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
		return services.AlertTicketService.GetAlertTicketRule(r)
	})
}

// List 获取告警转工单规则列表
func (atr alertTicketRuleController) List(ctx *gin.Context) {
	r := new(types.RequestAlertTicketRuleQuery)
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
		return services.AlertTicketService.ListAlertTicketRules(r)
	})
}

// Test 测试规则匹配
func (atr alertTicketRuleController) Test(ctx *gin.Context) {
	r := new(types.RequestAlertTicketRuleTest)
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
		return services.AlertTicketService.TestAlertTicketRule(r)
	})
}

// History 获取规则历史记录
func (atr alertTicketRuleController) History(ctx *gin.Context) {
	r := new(types.RequestAlertTicketRuleHistoryQuery)
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
		return services.AlertTicketService.GetAlertTicketRuleHistory(r)
	})
}

// Stats 获取规则统计报表
func (atr alertTicketRuleController) Stats(ctx *gin.Context) {
	r := new(types.RequestAlertTicketRuleStats)
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
		return services.AlertTicketService.GetAlertTicketRuleStats(r)
	})
}
