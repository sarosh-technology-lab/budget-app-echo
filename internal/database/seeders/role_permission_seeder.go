package seeders

import (
	"budget-backend/common"
	"budget-backend/internal/models"
)

func SeedRolePermissions() {
	db, err := common.Sql()
	if err != nil {
		panic(err)
	}

	var rolePermissions = []models.RolePermission{
		{RoleID: 1, PermissionID: 1},
		{RoleID: 1, PermissionID: 2},
		{RoleID: 1, PermissionID: 3},
		{RoleID: 1, PermissionID: 4},
		{RoleID: 1, PermissionID: 5},
		{RoleID: 1, PermissionID: 6},
		{RoleID: 1, PermissionID: 7},
		{RoleID: 1, PermissionID: 8},
		{RoleID: 2, PermissionID: 1},
		{RoleID: 2, PermissionID: 5},
	}
	db.Create(&rolePermissions)
}