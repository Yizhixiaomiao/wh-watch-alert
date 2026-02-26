package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

type userService struct {
	ctx *ctx.Context
}

type InterUserService interface {
	List(req interface{}) (interface{}, interface{})
	Get(req interface{}) (interface{}, interface{})
	Login(req interface{}) (interface{}, interface{})
	Update(req interface{}) (interface{}, interface{})
	Register(req interface{}) (interface{}, interface{})
	Delete(req interface{}) (interface{}, interface{})
	ChangePass(req interface{}) (interface{}, interface{})
	UpdateStatus(req interface{}) (interface{}, interface{})
	GetStatusHistory(req interface{}) (interface{}, interface{})
	GetPermissions(req interface{}) (interface{}, interface{})
	BatchOperation(req interface{}) (interface{}, interface{})
	GetActivityLogs(req interface{}) (interface{}, interface{})
}

func newInterUserService(ctx *ctx.Context) InterUserService {
	return &userService{
		ctx: ctx,
	}
}

func (us userService) List(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestUserQuery)

	data, err := us.ctx.DB.User().List(r.Query, r.JoinDuty, r.TenantId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (us userService) Get(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestUserQuery)

	data, _, err := us.ctx.DB.User().Get(r.UserId, r.UserName, r.Query)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (us userService) Login(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestUserLogin)
	originalPassword := r.Password
	r.Password = tools.GenerateHashPassword(r.Password)

	data, _, err := us.ctx.DB.User().Get("", r.UserName, "")
	if err != nil {
		return nil, err
	}

	setting, err := us.ctx.DB.Setting().Get()
	if err != nil {
		return nil, err
	}
	switch data.CreateBy {
	case "LDAP":
		if *setting.AuthType == models.SettingLdapAuth {
			err := LdapService.Login(r.UserName, originalPassword)
			if err != nil {
				logc.Error(us.ctx.Ctx, fmt.Sprintf("LDAP 用户登陆失败, err: %s", err.Error()))
				return nil, fmt.Errorf("LDAP 用户登陆失败, err: %s", err.Error())
			}
		} else {
			logc.Error(us.ctx.Ctx, "请先开启 LDAP 功能!")
			return nil, fmt.Errorf("请先开启 LDAP 功能!")
		}
	case "OIDC":
		logc.Error(us.ctx.Ctx, "请使用 OIDC 登录!")
		return nil, fmt.Errorf("请使用 OIDC 登录!")
	default:
		if data.Password != r.Password {
			return nil, fmt.Errorf("密码错误")
		}
	}

	tokenData, err := tools.GenerateToken(data.UserId, r.UserName, r.Password)
	if err != nil {
		return nil, err
	}

	duration := time.Duration(global.Config.Jwt.Expire) * time.Second
	us.ctx.Redis.Redis().Set("uid-"+data.UserId, tools.JsonMarshalToString(r), duration)

	return models.ResponseLoginInfo{
		Token:    tokenData,
		Username: r.UserName,
		UserId:   data.UserId,
	}, nil
}

func (us userService) Register(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestUserCreate)

	_, ok, _ := us.ctx.DB.User().Get("", r.UserName, "")
	if ok {
		return nil, fmt.Errorf("用户已存在")
	}

	// 在初始化admin用户时会固定一个userid，所以这里需要做一下判断；
	if r.UserId == "" {
		r.UserId = tools.RandUid()
	}

	if r.CreateBy == "" {
		r.CreateBy = "system"
	}

	err := us.ctx.DB.User().Create(models.Member{
		UserId:     r.UserId,
		UserName:   r.UserName,
		Email:      r.Email,
		Phone:      r.Phone,
		Password:   tools.GenerateHashPassword(r.Password),
		Role:       r.Role,
		CreateBy:   r.CreateBy,
		CreateAt:   time.Now().Unix(),
		JoinDuty:   r.JoinDuty,
		DutyUserId: r.DutyUserId,
		Tenants:    r.Tenants,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (us userService) Update(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestUserUpdate)
	var dbData models.Member

	db := us.ctx.DB.DB().Model(models.Member{})
	db.Where("user_id = ?", r.UserId).First(&dbData)

	oldRole := dbData.Role

	if r.Password == "" {
		r.Password = dbData.Password
	} else {
		r.Password = tools.GenerateHashPassword(r.Password)
	}
	err := us.ctx.DB.User().Update(models.Member{
		UserId:     r.UserId,
		UserName:   r.UserName,
		Email:      r.Email,
		Phone:      r.Phone,
		Password:   r.Password,
		Role:       r.Role,
		CreateBy:   r.CreateBy,
		CreateAt:   r.CreateAt,
		JoinDuty:   r.JoinDuty,
		DutyUserId: r.DutyUserId,
		Tenants:    r.Tenants,
	})
	if err != nil {
		return nil, err
	}

	if oldRole != r.Role {
		userData, _, err := us.ctx.DB.User().Get(r.UserId, "", "")
		if err != nil {
			logc.Error(us.ctx.Ctx, fmt.Sprintf("获取用户信息失败: %s", err.Error()))
			return nil, err
		}

		for _, tenantId := range userData.Tenants {
			err := us.ctx.DB.Tenant().ChangeTenantUserRole(tenantId, r.UserId, r.Role)
			if err != nil {
				logc.Error(us.ctx.Ctx, fmt.Sprintf("更新租户 %s 中用户 %s 的角色失败: %s", tenantId, r.UserId, err.Error()))
			} else {
				logc.Info(us.ctx.Ctx, fmt.Sprintf("成功更新租户 %s 中用户 %s 的角色从 %s 到 %s", tenantId, r.UserId, oldRole, r.Role))
			}
		}
	}

	us.ctx.DB.User().ChangeCache(r.UserId)

	return nil, nil
}

func (us userService) Delete(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestUserQuery)
	err := us.ctx.DB.User().Delete(r.UserId)
	if err != nil {
		return nil, err
	}

	us.ctx.DB.User().ChangeCache(r.UserId)

	return nil, nil
}

func (us userService) ChangePass(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestUserChangePassword)

	arr := md5.Sum([]byte(r.Password))
	hashPassword := hex.EncodeToString(arr[:])
	r.Password = hashPassword

	err := us.ctx.DB.User().ChangePass(r.UserId, r.Password)
	if err != nil {
		return nil, err
	}

	us.ctx.DB.User().ChangeCache(r.UserId)

	return nil, nil
}

func (us userService) UpdateStatus(req interface{}) (interface{}, interface{}) {
	return nil, nil
}

func (us userService) GetStatusHistory(req interface{}) (interface{}, interface{}) {
	return nil, nil
}

func (us userService) GetPermissions(req interface{}) (interface{}, interface{}) {
	r := req.(*types.RequestUserPermissionsQuery)

	tenantId := r.TenantId
	if tenantId == "" {
		tenantId = "default"
	}

	logc.Infof(us.ctx.Ctx, fmt.Sprintf("GetPermissions: UserId=%s, TenantId=%s", r.UserId, tenantId))

	tenantUserInfo, err := us.ctx.DB.Tenant().GetTenantLinkedUserInfo(tenantId, r.UserId)
	if err != nil {
		logc.Errorf(us.ctx.Ctx, fmt.Sprintf("GetTenantLinkedUserInfo failed: %v", err))
		return nil, err
	}

	logc.Infof(us.ctx.Ctx, fmt.Sprintf("GetTenantLinkedUserInfo: UserRole=%s", tenantUserInfo.UserRole))

	var role models.UserRole
	err = us.ctx.DB.DB().Model(&models.UserRole{}).Where("id = ?", tenantUserInfo.UserRole).First(&role).Error
	if err != nil {
		logc.Errorf(us.ctx.Ctx, fmt.Sprintf("GetUserRole failed: %v", err))
		return nil, err
	}

	logc.Infof(us.ctx.Ctx, fmt.Sprintf("GetUserRole: RoleId=%s, PermissionsCount=%d", role.ID, len(role.Permissions)))

	permKeyMap := make(map[string]string)
	userPerms, err := us.ctx.DB.UserPermissions().List()
	if err != nil {
		logc.Errorf(us.ctx.Ctx, fmt.Sprintf("GetUserPermissionsList failed: %v", err))
	} else {
		for _, perm := range userPerms {
			permKeyMap[perm.API] = perm.PermissionKey
		}
	}

	var permission []models.UserPermissions
	for _, perm := range role.Permissions {
		perm.PermissionKey = permKeyMap[perm.API]
		permission = append(permission, perm)
	}

	logc.Infof(us.ctx.Ctx, fmt.Sprintf("GetPermissions: Returning %d permissions", len(permission)))

	return permission, nil
}

func (us userService) BatchOperation(req interface{}) (interface{}, interface{}) {
	return nil, nil
}

func (us userService) GetActivityLogs(req interface{}) (interface{}, interface{}) {
	return nil, nil
}
