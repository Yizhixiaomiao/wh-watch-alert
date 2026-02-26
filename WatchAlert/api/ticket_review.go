package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
	"watchAlert/pkg/response"
)

type ticketReviewController struct{}

var TicketReviewController = new(ticketReviewController)

/*
工单评审 API
/api/w8t/ticket/review
*/
func (trc ticketReviewController) API(gin *gin.RouterGroup) {
	// 需要审计日志的操作
	a := gin.Group("ticket/review")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("assign", TicketReviewController.AssignReviewers)
		a.POST("submit", TicketReviewController.SubmitReview)
	}

	// 查询操作
	b := gin.Group("ticket/review")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("list", TicketReviewController.ListReviews)
		b.GET("get", TicketReviewController.GetReview)
	}

	// 评委管理
	c := gin.Group("ticket/reviewer")
	c.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		c.POST("create", TicketReviewController.CreateReviewer)
		c.POST("update", TicketReviewController.UpdateReviewer)
		c.POST("delete", TicketReviewController.DeleteReviewer)
		c.GET("list", TicketReviewController.ListReviewers)
		c.GET("get", TicketReviewController.GetReviewer)
	}
}

// AssignReviewers 分配评委
func (trc ticketReviewController) AssignReviewers(ctx *gin.Context) {
	r := new(types.RequestTicketReviewAssign)
	if err := ctx.ShouldBindJSON(&r); err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

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
	r.AssignedBy = uid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TicketReviewService.AssignReviewers(r)
	})
}

// SubmitReview 提交评审
func (trc ticketReviewController) SubmitReview(ctx *gin.Context) {
	r := new(types.RequestTicketReviewSubmit)
	if err := ctx.ShouldBindJSON(&r); err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	uid, exists := ctx.Get("UserID")
	if !exists {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("用户ID不存在")
		})
		return
	}
	r.ReviewedBy = uid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TicketReviewService.SubmitReview(r)
	})
}

// GetReview 获取评审详情
func (trc ticketReviewController) GetReview(ctx *gin.Context) {
	r := new(types.RequestTicketReviewQuery)
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
		return services.TicketReviewService.GetReview(r)
	})
}

// ListReviews 获取评审列表
func (trc ticketReviewController) ListReviews(ctx *gin.Context) {
	r := new(types.RequestTicketReviewQuery)
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
		return services.TicketReviewService.ListReviews(r)
	})
}

// CreateReviewer 创建评委
func (trc ticketReviewController) CreateReviewer(ctx *gin.Context) {
	r := new(types.RequestTicketReviewerCreate)
	if err := ctx.ShouldBindJSON(&r); err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

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
		return services.TicketReviewService.CreateReviewer(r)
	})
}

// UpdateReviewer 更新评委
func (trc ticketReviewController) UpdateReviewer(ctx *gin.Context) {
	r := new(types.RequestTicketReviewerUpdate)
	if err := ctx.ShouldBindJSON(&r); err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	tid, exists := ctx.Get("TenantID")
	if !exists {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("租户ID不存在")
		})
		return
	}
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TicketReviewService.UpdateReviewer(r)
	})
}

// DeleteReviewer 删除评委
func (trc ticketReviewController) DeleteReviewer(ctx *gin.Context) {
	r := new(types.RequestTicketReviewerUpdate)
	if err := ctx.ShouldBindJSON(&r); err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	tid, exists := ctx.Get("TenantID")
	if !exists {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("租户ID不存在")
		})
		return
	}
	r.TenantId = tid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.TicketReviewService.DeleteReviewer(r)
	})
}

// GetReviewer 获取评委详情
func (trc ticketReviewController) GetReviewer(ctx *gin.Context) {
	r := new(types.RequestTicketReviewerQuery)
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
		return services.TicketReviewService.GetReviewer(r)
	})
}

// ListReviewers 获取评委列表
func (trc ticketReviewController) ListReviewers(ctx *gin.Context) {
	r := new(types.RequestTicketReviewerQuery)
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
		return services.TicketReviewService.ListReviewers(r)
	})
}
