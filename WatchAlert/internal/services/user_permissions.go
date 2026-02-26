package services

import (
	"watchAlert/internal/ctx"
)

type (
	userPermissionService struct {
		ctx *ctx.Context
	}

	InterUserPermissionService interface {
		List() (interface{}, interface{})
	}
)

func newInterUserPermissionService(ctx *ctx.Context) InterUserPermissionService {
	return &userPermissionService{
		ctx: ctx,
	}
}

func (up userPermissionService) List() (interface{}, interface{}) {
	permissions, err := up.ctx.DB.UserPermissions().List()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(permissions))

	for _, perm := range permissions {
		permMap := map[string]interface{}{
			"key":           perm.Key,
			"permissionKey": perm.PermissionKey,
			"api":           perm.API,
			"category":      perm.Category,
			"subCategory":   perm.SubCategory,
			"order":         perm.Order,
			"title":         perm.Title,
		}
		result = append(result, permMap)
	}

	return result, nil
}
