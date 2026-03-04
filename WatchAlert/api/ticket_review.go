package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
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
}

// AssignReviewers 分配评委
func (trc ticketReviewController) AssignReviewers(ctx *gin.Context) {
	r := new(types.RequestTicketReviewAssign)
	if err := BindJson(ctx, r); err != nil {
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
	if err := BindJson(ctx, r); err != nil {
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
	if err := BindQuery(ctx, r); err != nil {
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
		return services.TicketReviewService.GetReview(r)
	})
}

// ListReviews 获取评审列表
func (trc ticketReviewController) ListReviews(ctx *gin.Context) {
	r := new(types.RequestTicketReviewQuery)
	if err := BindQuery(ctx, r); err != nil {
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
		return services.TicketReviewService.ListReviews(r)
	})
}
