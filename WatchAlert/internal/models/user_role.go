package models

type UserRole struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Permissions []UserPermissions `json:"permissions" gorm:"permissions;serializer:json"`
	ParentRoles []UserRole        `json:"parentRoles" gorm:"many2many:role_inheritance;joinForeignKey:user_role_id;joinReferences:parent_role_id"`
	UpdateAt    int64             `json:"updateAt"`
}
