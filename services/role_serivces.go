package services

import (
	"fmt"
	"github.com/boof/umg/db"
	"github.com/boof/umg/rbac/roles"
	"github.com/boof/umg/rbac/users"
	"github.com/boof/umg/rest_errors"
)

func GetAllRoles(count, page int64, order, sortBy string) ([]roles.Role, rest_errors.Error) {
	res, err := roles.GetAllRoles(int(count), int(page), order, sortBy)
	if err != nil {
		return nil, rest_errors.NewNotFoundError(err.Error())
	}

	return res, nil
}

func RemoveRoleByID(roleID int64) rest_errors.Error {
	removePoliciesByRole(roleID)

	var allUsers []users.User
	query := fmt.Sprintf("SELECT * FROM \"user\" WHERE role_ids LIKE '%%%d,%%' OR role_ids LIKE '%%%d]%%'", roleID, roleID)
	db.Engine.SQL(query).Find(&allUsers)
	for _, user := range allUsers {
		user.DisallowRole(roleID)
	}

	if _, err := db.Engine.ID(roleID).Delete(&roles.Role{}); err != nil {
		return rest_errors.NewInternalServerError("Unable to delete a role", err)
	}

	return nil
}
