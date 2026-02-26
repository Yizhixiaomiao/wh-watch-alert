package initialization

import (
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"

	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
)

var perms []models.UserPermissions

func InitPermissionsSQL(ctx *ctx.Context) {
	// 从数据库读取权限数据（如果已存在）
	var psData []models.UserPermissions
	err := ctx.DB.DB().Model(&models.UserPermissions{}).Find(&psData).Error

	if err != nil {
		logc.Errorf(ctx.Ctx, "读取权限数据失败: %v", err)
		return
	}

	// 如果数据库中没有权限数据，则不执行初始化
	if len(psData) == 0 {
		logc.Infof(ctx.Ctx, "数据库中无权限数据，跳过初始化")
		return
	}

	perms = psData
	logc.Infof(ctx.Ctx, "从数据库加载了 %d 个权限", len(psData))
}

func InitUserRolesSQL(ctx *ctx.Context) {
	var adminRole models.UserRole
	var db = ctx.DB.DB().Model(&models.UserRole{})

	roles := models.UserRole{
		ID:          "admin",
		Name:        "admin",
		Description: "system",
		Permissions: perms,
		UpdateAt:    time.Now().Unix(),
	}

	err := db.Where("name = ?", "admin").First(&adminRole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ctx.DB.DB().Create(&roles).Error
		}
	} else {
		err = db.Where("name = ?", "admin").Updates(models.UserRole{Permissions: perms}).Error
	}

	if err != nil {
		logc.Errorf(ctx.Ctx, err.Error())
		panic(err)
	}
	logc.Infof(ctx.Ctx, "系统初始化完成，admin 角色包含 %d 项权限", len(perms))
}
