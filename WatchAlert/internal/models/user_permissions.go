package models

type UserPermissions struct {
	Key           string `json:"key" gorm:"primaryKey;column:key;type:longtext"`
	PermissionKey string `json:"permissionKey" gorm:"column:permission_key;type:varchar(100)"`
	Title         string `json:"title" gorm:"column:title;type:varchar(255)"` // 权限描述（中文）
	API           string `json:"api" gorm:"column:api;type:longtext"`
	Category      string `json:"category" gorm:"column:category;type:varchar(50)"`        // 一级分类
	SubCategory   string `json:"subCategory" gorm:"column:sub_category;type:varchar(50)"` // 二级分类
	Order         int    `json:"order" gorm:"column:order;type:int"`                      // 排序
}

func TableName() string {
	return "user_permissions"
}

