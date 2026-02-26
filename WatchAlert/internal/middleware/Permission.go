package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/response"
	utils2 "watchAlert/pkg/tools"
)

func getRoleWithInheritedPermissions(db *gorm.DB, roleID string) ([]models.UserPermissions, error) {
	var role models.UserRole
	err := db.Preload("ParentRoles").Where("id = ?", roleID).First(&role).Error
	if err != nil {
		return nil, err
	}

	permissionMap := make(map[string]models.UserPermissions)

	var collectPermissions func(r models.UserRole)
	collectPermissions = func(r models.UserRole) {
		permissions := r.Permissions
		for i := range permissions {
			perm := permissions[i]
			permissionMap[perm.API] = perm
		}
		for _, parent := range r.ParentRoles {
			collectPermissions(parent)
		}
	}

	collectPermissions(role)

	permissions := make([]models.UserPermissions, 0, len(permissionMap))
	for _, perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func Permission() gin.HandlerFunc {
	return func(context *gin.Context) {
		tid := context.Request.Header.Get(TenantIDHeaderKey)
		if tid == "null" || tid == "" {
			return
		}
		tokenStr := context.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.TokenFail(context)
			context.Abort()
			return
		}

		userId := utils2.GetUserID(tokenStr)
		c := ctx.DO()

		var user models.Member
		err := c.DB.DB().Model(&models.Member{}).Where("user_id = ?", userId).First(&user).Error
		if gorm.ErrRecordNotFound == err {
			logc.Errorf(c.Ctx, fmt.Sprintf("用户不存在, uid: %s", userId))
		}
		if err != nil {
			response.PermissionFail(context)
			context.Abort()
			return
		}

		logc.Infof(c.Ctx, fmt.Sprintf("权限检查: UserName=%s, Role=%s, CreateBy=%s", user.UserName, user.Role, user.CreateBy))

		if user.UserName == "admin" || user.Role == "admin" {
			logc.Infof(c.Ctx, fmt.Sprintf("用户 %s 是超级管理员，跳过权限检查", user.UserName))
			context.Set("UserID", user.UserId)
			context.Set("UserId", user.UserId)
			context.Set("UserEmail", user.Email)
			context.Next()
			return
		}

		context.Set("UserID", user.UserId)
		context.Set("UserId", user.UserId)
		context.Set("UserEmail", user.Email)

		tenantUserInfo, err := c.DB.Tenant().GetTenantLinkedUserInfo(tid, userId)
		if err != nil {
			logc.Errorf(c.Ctx, fmt.Sprintf("获取租户用户角色失败 %s", err.Error()))
			response.TokenFail(context)
			context.Abort()
			return
		}

		permission, err := getRoleWithInheritedPermissions(c.DB.DB(), tenantUserInfo.UserRole)
		if err != nil {
			response.Fail(context, fmt.Sprintf("获取用户 %s 的角色权限失败, %s", user.UserName, err.Error()), "failed")
			logc.Errorf(c.Ctx, fmt.Sprintf("获取用户 %s 的角色权限失败 %s", user.UserName, err.Error()))
			context.Abort()
			return
		}

		urlPath := context.Request.URL.Path

		var pass bool
		for _, v := range permission {
			if urlPath == v.API {
				pass = true
				break
			}
		}

		if !pass {
			response.PermissionFail(context)
			context.Abort()
			return
		}
	}
}
