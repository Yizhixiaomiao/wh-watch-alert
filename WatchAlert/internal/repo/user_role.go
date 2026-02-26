package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
)

type (
	UserRoleRepo struct {
		entryRepo
	}

InterUserRoleRepo interface {
		List() ([]models.UserRole, error)
		Create(r models.UserRole) error
		Update(r models.UserRole) error
		Delete(id string) error
		GetByName(name string) (models.UserRole, error)
		GetWithPermissions(id string) (models.UserRole, error)
		GetWithInheritedPermissions(id string) (models.UserRole, error)
		ValidateInheritance(parentID, childID string) (bool, error)
	}
)

func newUserRoleInterface(db *gorm.DB, g InterGormDBCli) InterUserRoleRepo {
	return &UserRoleRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (ur UserRoleRepo) List() ([]models.UserRole, error) {
	var (
		data []models.UserRole
		db   = ur.DB().Model(&models.UserRole{})
	)

	err := db.Where("id != ?", "admin").Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (ur UserRoleRepo) Create(r models.UserRole) error {
	err := ur.g.Create(&models.UserRole{}, &r)
	if err != nil {
		return err
	}

	return nil
}

func (ur UserRoleRepo) Update(r models.UserRole) error {
	u := Updates{
		Table: models.UserRole{},
		Where: map[string]interface{}{
			"id = ?": r.ID,
		},
		Updates: r,
	}

	err := ur.g.Updates(u)
	if err != nil {
		return err
	}

	return nil
}

func (ur UserRoleRepo) Delete(id string) error {
	d := Delete{
		Table: models.UserRole{},
		Where: map[string]interface{}{
			"id = ?": id,
		},
	}

	err := ur.g.Delete(d)
	if err != nil {
		return err
	}

	return nil
}

func (ur UserRoleRepo) GetByName(name string) (models.UserRole, error) {
	var role models.UserRole
	err := ur.DB().Model(&models.UserRole{}).Where("name = ?", name).First(&role).Error
	if err != nil {
		return role, err
	}
	return role, nil
}

func (ur UserRoleRepo) GetWithPermissions(id string) (models.UserRole, error) {
	var role models.UserRole
	err := ur.DB().Preload("Permissions").Where("id = ?", id).First(&role).Error
	return role, err
}

func (ur UserRoleRepo) GetWithInheritedPermissions(id string) (models.UserRole, error) {
	var role models.UserRole
	err := ur.DB().Preload("Permissions").Preload("ParentRoles").Where("id = ?", id).First(&role).Error
	return role, err
}

func (ur UserRoleRepo) ValidateInheritance(parentID, childID string) (bool, error) {
	var count int64
	err := ur.DB().Raw(`
		WITH RECURSIVE inheritance_chain AS (
			SELECT parent_role_id FROM role_inheritance WHERE role_id = ?
			UNION
			SELECT ri.parent_role_id FROM role_inheritance ri
			INNER JOIN inheritance_chain ic ON ri.role_id = ic.parent_role_id
		)
		SELECT COUNT(*) FROM inheritance_chain WHERE parent_role_id = ?
	`, childID, parentID).Scan(&count).Error

	if err != nil {
		return false, err
	}

	return count == 0, nil
}
