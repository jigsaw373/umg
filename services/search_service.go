package services

import (
	"fmt"
	"github.com/boof/umg/rbac/access"
	"time"

	"github.com/boof/umg/db"
	"github.com/boof/umg/rbac/roles"
	"github.com/boof/umg/rbac/users"
	"github.com/boof/umg/settings"
	"github.com/boof/umg/util/datetime"
)

func SearchUsers(text string) []*users.RoleUser {
	res := make([]*users.RoleUser, 0)

	if text == "" {
		return res
	}

	var search []users.User
	query := fmt.Sprintf("SELECT * FROM \"user\" "+
		"WHERE username ILIKE '%%%s%%' OR name ILIKE '%%%s%%' or email ILIKE '%%%s%%'", text, text, text)

	db.Engine.SQL(query).Limit(40).Find(&search)

	for _, user := range search {
		ru := &users.RoleUser{
			ID:        user.ID,
			Username:  user.Username,
			Name:      user.Name,
			Email:     user.Email,
			Online:    db.IsOnline(user.ID),
			CreatedAt: user.CreatedAt.Format("Jan 02, 2006"),
			Roles:     make([]*roles.Role, 0),
		}

		if user.LastLogin.Before(time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)) {
			ru.LastLogin = "-"
		} else {
			ru.LastLogin = datetime.DurationString(user.LastLogin)
		}

		for _, r := range user.RoleIDs {
			role, err := (&roles.Role{ID: r}).GetByID()
			if err == nil {
				ru.Roles = append(ru.Roles, role)
			}
		}

		expire, err := (&access.Expire{UserID: user.ID}).GetByUserID()
		if err == nil {
			ru.ExpireAt = expire.ExpireAt.Format(settings.UserDTLayout)
		} else {
			ru.ExpireAt = "undefined"
		}

		res = append(res, ru)
	}

	return res
}

func SearchRoles(text string) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)

	if text == "" {
		return res
	}

	var search []roles.Role
	query := fmt.Sprintf("SELECT * FROM role WHERE name ILIKE '%%%s%%'", text)
	db.Engine.SQL(query).Limit(20).Find(&search)

	for _, role := range search {
		out := make(map[string]interface{})
		out["id"] = role.ID
		out["name"] = role.Name

		policies, err := GetNamedPolicies(role.ID)
		if err == nil {
			out["policies"] = policies
		}

		res = append(res, out)
	}

	return res
}
