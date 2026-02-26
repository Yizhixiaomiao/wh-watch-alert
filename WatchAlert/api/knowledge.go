package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
)

type knowledgeController struct{}

var KnowledgeController = new(knowledgeController)

/*
知识库 API
/api/w8t/knowledge
*/
func (kc knowledgeController) API(gin *gin.RouterGroup) {
	// 需要审计日志的操作
	a := gin.Group("knowledge")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.POST("create", KnowledgeController.CreateKnowledge)
		a.POST("update", KnowledgeController.UpdateKnowledge)
		a.POST("delete", KnowledgeController.DeleteKnowledge)
		a.POST("like", KnowledgeController.LikeKnowledge)
		a.POST("save-to-ticket", KnowledgeController.SaveToTicket)
		a.POST("category/create", KnowledgeController.CreateCategory)
		a.POST("category/update", KnowledgeController.UpdateCategory)
		a.POST("category/delete", KnowledgeController.DeleteCategory)
	}

	// 查询操作
	b := gin.Group("knowledge")
	b.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
	)
	{
		b.GET("list", KnowledgeController.ListKnowledges)
		b.GET("get", KnowledgeController.GetKnowledge)
		b.GET("category/list", KnowledgeController.ListCategories)
		b.GET("category/get", KnowledgeController.GetCategory)
	}
}

// CreateKnowledge 创建知识
func (kc knowledgeController) CreateKnowledge(ctx *gin.Context) {
	r := new(types.RequestKnowledgeCreate)
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
	r.AuthorId = uid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.KnowledgeService.CreateKnowledge(r)
	})
}

// UpdateKnowledge 更新知识
func (kc knowledgeController) UpdateKnowledge(ctx *gin.Context) {
	r := new(types.RequestKnowledgeUpdate)
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
		return services.KnowledgeService.UpdateKnowledge(r)
	})
}

// DeleteKnowledge 删除知识
func (kc knowledgeController) DeleteKnowledge(ctx *gin.Context) {
	r := new(types.RequestKnowledgeDelete)
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
		return services.KnowledgeService.DeleteKnowledge(r)
	})
}

// GetKnowledge 获取知识详情
func (kc knowledgeController) GetKnowledge(ctx *gin.Context) {
	r := new(types.RequestKnowledgeQuery)
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
		return services.KnowledgeService.GetKnowledge(r)
	})
}

// ListKnowledges 获取知识列表
func (kc knowledgeController) ListKnowledges(ctx *gin.Context) {
	r := new(types.RequestKnowledgeQuery)
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
		return services.KnowledgeService.ListKnowledges(r)
	})
}

// LikeKnowledge 点赞知识
func (kc knowledgeController) LikeKnowledge(ctx *gin.Context) {
	r := new(types.RequestKnowledgeLike)
	BindJson(ctx, r)

	uid, exists := ctx.Get("UserID")
	if !exists {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("用户ID不存在")
		})
		return
	}
	r.UserId = uid.(string)

	Service(ctx, func() (interface{}, interface{}) {
		return services.KnowledgeService.LikeKnowledge(r)
	})
}

// SaveToTicket 保存知识到工单
func (kc knowledgeController) SaveToTicket(ctx *gin.Context) {
	r := new(types.RequestKnowledgeSaveToTicket)
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
		return services.KnowledgeService.SaveToTicket(r)
	})
}

// CreateCategory 创建分类
func (kc knowledgeController) CreateCategory(ctx *gin.Context) {
	r := new(types.RequestKnowledgeCategoryCreate)
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
		return services.KnowledgeService.CreateCategory(r)
	})
}

// UpdateCategory 更新分类
func (kc knowledgeController) UpdateCategory(ctx *gin.Context) {
	r := new(types.RequestKnowledgeCategoryUpdate)
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
		return services.KnowledgeService.UpdateCategory(r)
	})
}

// DeleteCategory 删除分类
func (kc knowledgeController) DeleteCategory(ctx *gin.Context) {
	r := new(types.RequestKnowledgeCategoryDelete)
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
		return services.KnowledgeService.DeleteCategory(r)
	})
}

// GetCategory 获取分类详情
func (kc knowledgeController) GetCategory(ctx *gin.Context) {
	r := new(types.RequestKnowledgeCategoryQuery)
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
		return services.KnowledgeService.GetCategory(r)
	})
}

// ListCategories 获取分类列表
func (kc knowledgeController) ListCategories(ctx *gin.Context) {
	r := new(types.RequestKnowledgeCategoryQuery)
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
		return services.KnowledgeService.ListCategories(r)
	})
}
