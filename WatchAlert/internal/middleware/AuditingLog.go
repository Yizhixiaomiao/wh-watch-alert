package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logc"
	"io"
	"io/ioutil"
	"strings"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/response"
	"watchAlert/pkg/tools"
)

func AuditingLog() gin.HandlerFunc {
	return func(context *gin.Context) {
		// Operation user
		var username string
		createBy := tools.GetUser(context.Request.Header.Get("Authorization"))
		if createBy != "" {
			username = createBy
		} else {
			username = "用户未登录"
		}

		// Response log
		body := context.Request.Body
		readBody, err := io.ReadAll(body)
		if err != nil {
			logc.Error(ctx.DO().Ctx, err)
			return
		}
		// 将 body 数据放回请求中
		context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(readBody))

		tid := context.Request.Header.Get(TenantIDHeaderKey)
		if tid == "" {
			response.Fail(context, "租户ID不能为空", "failed")
			context.Abort()
			return
		}

		// 当请求处理完成后才会执行 Next() 后面的代码
		context.Next()

		// 从数据库获取权限信息
		var permissions []models.UserPermissions
		err = ctx.DO().DB.DB().Model(&models.UserPermissions{}).Find(&permissions).Error
		if err != nil {
			logc.Error(ctx.DO().Ctx, err)
		}

		// 根据 API 路径查找对应的权限
		var auditType string
		for _, perm := range permissions {
			if strings.Contains(context.Request.URL.Path, perm.API) {
				auditType = perm.Key
				break
			}
		}

		auditLog := models.AuditLog{
			TenantId:   tid,
			ID:         "Trace" + tools.RandId(),
			Username:   username,
			IPAddress:  context.ClientIP(),
			Method:     context.Request.Method,
			Path:       context.Request.URL.Path,
			CreatedAt:  time.Now().Unix(),
			StatusCode: context.Writer.Status(),
			Body:       string(readBody),
			AuditType:  auditType,
		}

		c := ctx.DO()
		err = c.DB.AuditLog().Create(auditLog)
		if err != nil {
			response.Fail(context, "审计日志写入数据库失败, "+err.Error(), "failed")
			context.Abort()
			return
		}
	}
}
