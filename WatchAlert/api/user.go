package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	middleware "watchAlert/internal/middleware"
	"watchAlert/internal/services"
	"watchAlert/internal/types"
	jwtUtils "watchAlert/pkg/tools"
)

type userController struct{}

var UserController = new(userController)

/*
用户 API
/api/w8t/user
*/
func (userController userController) API(gin *gin.RouterGroup) {

	a := gin.Group("user")
	a.Use(
		middleware.Auth(),
		middleware.Permission(),
		middleware.ParseTenant(),
		middleware.AuditingLog(),
	)
	{
		a.GET("userList", userController.List)
		a.POST("userUpdate", userController.Update)
		a.POST("userDelete", userController.Delete)
		a.POST("userChangePass", userController.ChangePass)
		a.POST("userUpdateStatus", userController.UpdateStatus)
		a.GET("userStatusHistory", userController.GetStatusHistory)
		a.POST("batchOperation", userController.BatchOperation)
		a.GET("activityLogs", userController.GetActivityLogs)
	}

	// 单独处理 userPermissions，允许用户访问自己的权限信息
	c := gin.Group("user")
	c.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		c.GET("userPermissions", userController.GetPermissions)
	}

}

func (userController userController) List(ctx *gin.Context) {
	r := new(types.RequestUserQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.List(r)
	})
}

func (userController userController) GetUserInfo(ctx *gin.Context) {
	r := new(types.RequestUserQuery)
	BindQuery(ctx, r)

	token := ctx.Request.Header.Get("Authorization")
	username := jwtUtils.GetUser(token)
	r.UserName = username

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Get(r)
	})
}

func (userController userController) Login(ctx *gin.Context) {
	r := new(types.RequestUserLogin)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Login(r)
	})
}

func (userController userController) Register(ctx *gin.Context) {
	r := new(types.RequestUserCreate)
	BindJson(ctx, r)

	createUser := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
	r.CreateBy = createUser

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Register(r)
	})
}

func (userController userController) Update(ctx *gin.Context) {
	r := new(types.RequestUserUpdate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Update(r)
	})
}

func (userController userController) Delete(ctx *gin.Context) {
	r := new(types.RequestUserQuery)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Delete(r)
	})
}

func (userController userController) CheckUser(ctx *gin.Context) {
	r := new(types.RequestUserQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.Get(r)
	})
}

func (userController userController) ChangePass(ctx *gin.Context) {
	r := new(types.RequestUserChangePassword)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.ChangePass(r)
	})
}

func (userController userController) UpdateStatus(ctx *gin.Context) {
	r := new(types.RequestUserStatusUpdate)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.UpdateStatus(r)
	})
}

func (userController userController) GetStatusHistory(ctx *gin.Context) {
	r := new(types.RequestUserStatusQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.GetStatusHistory(r)
	})
}

func (userController userController) GetPermissions(ctx *gin.Context) {
	r := new(types.RequestUserPermissionsQuery)
	BindQuery(ctx, r)

	// 获取 TenantID
	tenantId := ctx.GetHeader("TenantID")
	if tenantId == "" {
		tenantId = "default"
	}
	r.TenantId = tenantId

	// 从上下文获取当前用户ID
	currentUserID, exists := ctx.Get("UserID")
	if !exists {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("无法获取当前用户信息")
		})
		return
	}

	// 如果未指定用户ID，则使用当前用户ID（只允许查看自己的权限）
	if r.UserId == "" {
		// 尝试从userId参数获取（兼容前端可能的驼峰命名）
		userId := ctx.Query("userId")
		if userId != "" {
			r.UserId = userId
		} else {
			r.UserId = currentUserID.(string)
		}
	}

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.GetPermissions(r)
	})
}

func (userController userController) BatchOperation(ctx *gin.Context) {
	r := new(types.RequestUserBatchOperation)
	BindJson(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.BatchOperation(r)
	})
}

func (userController userController) GetActivityLogs(ctx *gin.Context) {
	r := new(types.RequestUserActivityLogQuery)
	BindQuery(ctx, r)

	Service(ctx, func() (interface{}, interface{}) {
		return services.UserService.GetActivityLogs(r)
	})
}
