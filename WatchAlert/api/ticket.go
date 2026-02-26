package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type ticketController struct{}

var TicketController = new(ticketController)

/*
工单管理 API
/api/w8t/ticket
*/
func (ticketController ticketController) API(gin *gin.RouterGroup) {
	// 需要审计日志的操作
	a := gin.Group("ticket")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("create", TicketController.Create)
		a.POST("update", TicketController.Update)
		a.POST("delete", TicketController.Delete)
		a.POST("assign", TicketController.Assign)
		a.POST("claim", TicketController.Claim)
		a.POST("transfer", TicketController.Transfer)
		a.POST("escalate", TicketController.Escalate)
		a.POST("resolve", TicketController.Resolve)
		a.POST("close", TicketController.Close)
		a.POST("reopen", TicketController.Reopen)
		a.POST("comment", TicketController.AddComment)
		a.POST("step/add", TicketController.AddStep)
		a.POST("step/update", TicketController.UpdateStep)
		a.POST("step/delete", TicketController.DeleteStep)
		a.POST("step/reorder", TicketController.ReorderSteps)
	}

	// 查询操作
	b := gin.Group("ticket")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("list", TicketController.List)
		b.GET("get", TicketController.Get)
		b.GET("comments", TicketController.GetComments)
		b.GET("worklog", TicketController.GetWorkLogs)
		b.GET("statistics", TicketController.GetStatistics)
		b.GET("steps", TicketController.GetSteps)
	}

	// 模板管理
	tmpl := gin.Group("ticket/template")
	tmpl.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		tmpl.POST("create", TicketController.CreateTemplate)
		tmpl.POST("update", TicketController.UpdateTemplate)
		tmpl.POST("delete", TicketController.DeleteTemplate)
		tmpl.GET("list", TicketController.ListTemplates)
	}

	// SLA策略管理
	sla := gin.Group("ticket/sla")
	sla.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		sla.POST("create", TicketController.CreateSLAPolicy)
		sla.POST("update", TicketController.UpdateSLAPolicy)
		sla.POST("delete", TicketController.DeleteSLAPolicy)
		sla.GET("list", TicketController.ListSLAPolicies)
	}

	// 移动端接口 (不需要认证，支持公开访问)
	mobile := gin.Group("mobile/ticket")
	{
		mobile.POST("create", TicketController.MobileCreate)
		mobile.GET("query", TicketController.MobileQuery)
	}
}

// Create 创建工单
func (tc ticketController) Create(ctx *gin.Context) {
	r := new(types.RequestTicketCreate)
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

	if r.Source == "" {
		r.Source = models.TicketSourceManual
	}

	Service(ctx, func() (interface{}, interface{}) {
		return services.TicketService.Create(r)
	})
}

// Update 更新工单
func (tc ticketController) Update(ctx *gin.Context) {
	r := new(types.RequestTicketUpdate)
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
		return services.TicketService.Update(r)
	})
}

// Delete 删除工单
func (tc ticketController) Delete(ctx *gin.Context) {
	r := new(types.RequestTicketDelete)
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
		return services.TicketService.Delete(r)
	})
}

// Get 获取工单详情
func (tc ticketController) Get(ctx *gin.Context) {
	r := new(types.RequestTicketQuery)
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
		return services.TicketService.Get(r)
	})
}

// List 获取工单列表
func (tc ticketController) List(ctx *gin.Context) {
	r := new(types.RequestTicketQuery)
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
		return services.TicketService.List(r)
	})
}

// Assign 分配工单
func (tc ticketController) Assign(ctx *gin.Context) {
	r := new(types.RequestTicketAssign)
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
	r.UserId = uid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TicketService.Assign(r)
	})
}

// Claim 认领工单
func (tc ticketController) Claim(ctx *gin.Context) {
	r := new(types.RequestTicketClaim)
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
	r.UserId = uid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TicketService.Claim(r)
	})
}

// Transfer 转派工单
func (tc ticketController) Transfer(ctx *gin.Context) {
	r := new(types.RequestTicketTransfer)
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
	r.UserId = uid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TicketService.Transfer(r)
	})
}

// Escalate 升级工单
func (tc ticketController) Escalate(ctx *gin.Context) {
	r := new(types.RequestTicketEscalate)
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
	r.UserId = uid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TicketService.Escalate(r)
	})
}

// Resolve 标记解决
func (tc ticketController) Resolve(ctx *gin.Context) {
	r := new(types.RequestTicketResolve)
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
	r.UserId = uid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TicketService.Resolve(r)
	})
}

// Close 关闭工单
func (tc ticketController) Close(ctx *gin.Context) {
	r := new(types.RequestTicketClose)
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
	r.UserId = uid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TicketService.Close(r)
	})
}

// Reopen 重新打开工单
func (tc ticketController) Reopen(ctx *gin.Context) {
	r := new(types.RequestTicketReopen)
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
	r.UserId = uid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TicketService.Reopen(r)
	})
}

// AddComment 添加评论
func (tc ticketController) AddComment(ctx *gin.Context) {
	r := new(types.RequestTicketComment)
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
	r.UserId = uid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TicketService.AddComment(r)
	})
}

// GetComments 获取评论列表
func (tc ticketController) GetComments(ctx *gin.Context) {
	r := new(types.RequestTicketCommentQuery)
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
		return services.TicketService.GetComments(r)
	})
}

// GetWorkLogs 获取工作日志
func (tc ticketController) GetWorkLogs(ctx *gin.Context) {
	r := new(types.RequestTicketWorkLogQuery)
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
		return services.TicketService.GetWorkLogs(r)
	})
}

// GetStatistics 获取工单统计
func (tc ticketController) GetStatistics(ctx *gin.Context) {
	r := new(types.RequestTicketStatistics)
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
		return services.TicketService.GetStatistics(r)
	})
}

// CreateTemplate 创建工单模板
func (tc ticketController) CreateTemplate(ctx *gin.Context) {
	r := new(types.RequestTicketTemplateCreate)
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
		return services.TicketService.CreateTemplate(r)
	})
}

// UpdateTemplate 更新工单模板
func (tc ticketController) UpdateTemplate(ctx *gin.Context) {
	r := new(types.RequestTicketTemplateUpdate)
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
		return services.TicketService.UpdateTemplate(r)
	})
}

// DeleteTemplate 删除工单模板
func (tc ticketController) DeleteTemplate(ctx *gin.Context) {
	r := new(types.RequestTicketTemplateDelete)
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
		return services.TicketService.DeleteTemplate(r)
	})
}

// ListTemplates 获取工单模板列表
func (tc ticketController) ListTemplates(ctx *gin.Context) {
	r := new(types.RequestTicketTemplateQuery)
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
		return services.TicketService.ListTemplates(r)
	})
}

// CreateSLAPolicy 创建SLA策略
func (tc ticketController) CreateSLAPolicy(ctx *gin.Context) {
	r := new(types.RequestTicketSLAPolicyCreate)
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
		return services.TicketService.CreateSLAPolicy(r)
	})
}

// UpdateSLAPolicy 更新SLA策略
func (tc ticketController) UpdateSLAPolicy(ctx *gin.Context) {
	r := new(types.RequestTicketSLAPolicyUpdate)
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
		return services.TicketService.UpdateSLAPolicy(r)
	})
}

// DeleteSLAPolicy 删除SLA策略
func (tc ticketController) DeleteSLAPolicy(ctx *gin.Context) {
	r := new(types.RequestTicketSLAPolicyDelete)
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
		return services.TicketService.DeleteSLAPolicy(r)
	})
}

// ListSLAPolicies 获取SLA策略列表
func (tc ticketController) ListSLAPolicies(ctx *gin.Context) {
	r := new(types.RequestTicketSLAPolicyQuery)
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
		return services.TicketService.ListSLAPolicies(r)
	})
}

// MobileCreate 移动端创建工单
func (tc ticketController) MobileCreate(ctx *gin.Context) {
	r := new(types.RequestMobileTicketCreate)
	BindJson(ctx, r)

	// 设置默认值
	if r.TenantId == "" {
		r.TenantId = "default"
	}
	if r.Type == "" {
		r.Type = models.TicketTypeFault
	}
	if r.Priority == "" {
		r.Priority = models.TicketPriorityP2
	}

	// 获取客户端信息
	r.UserAgent = ctx.GetHeader("User-Agent")
	r.Platform = detectPlatform(ctx.GetHeader("User-Agent"))

	Service(ctx, func() (interface{}, interface{}) {
		return services.TicketService.MobileCreate(r)
	})
}

// MobileQuery 移动端查询工单状态
func (tc ticketController) MobileQuery(ctx *gin.Context) {
	r := new(types.RequestMobileTicketQuery)
	BindQuery(ctx, r)

	// 设置默认租户ID
	if r.TenantId == "" {
		r.TenantId = "default"
	}

	Service(ctx, func() (interface{}, interface{}) {
		return services.TicketService.MobileQuery(r)
	})
}

// detectPlatform 检测客户端平台
func detectPlatform(userAgent string) string {
	if userAgent == "" {
		return "unknown"
	}

	// 检测微信
	if strings.Contains(strings.ToLower(userAgent), "micromessenger") {
		return "wechat"
	}

	// 检测移动设备
	if strings.Contains(strings.ToLower(userAgent), "mobile") ||
		strings.Contains(strings.ToLower(userAgent), "android") ||
		strings.Contains(strings.ToLower(userAgent), "iphone") ||
		strings.Contains(strings.ToLower(userAgent), "ipad") {
		return "mobile"
	}

	return "web"
}

// AddStep 添加处理步骤
func (tc ticketController) AddStep(ctx *gin.Context) {
	r := new(types.RequestTicketStepCreate)
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
		return services.TicketService.AddStep(r)
	})
}

// UpdateStep 更新处理步骤
func (tc ticketController) UpdateStep(ctx *gin.Context) {
	r := new(types.RequestTicketStepUpdate)
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
		return services.TicketService.UpdateStep(r)
	})
}

// DeleteStep 删除处理步骤
func (tc ticketController) DeleteStep(ctx *gin.Context) {
	r := new(types.RequestTicketStepDelete)
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
		return services.TicketService.DeleteStep(r)
	})
}

// ReorderSteps 重新排序步骤
func (tc ticketController) ReorderSteps(ctx *gin.Context) {
	r := new(types.RequestTicketStepReorder)
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
		return services.TicketService.ReorderSteps(r)
	})
}

// GetSteps 获取处理步骤列表
func (tc ticketController) GetSteps(ctx *gin.Context) {
	r := new(types.RequestTicketStepsQuery)
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
		return services.TicketService.GetSteps(r)
	})
}
