package auth

import (
	"errors"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"github.com/boof/umg/db"
	"github.com/boof/umg/rbac/users"
	"github.com/boof/umg/services"
)

const (
	AdminUser = "ADMIN"
)

// AdminHandler forces that the next handler function to has an admin role
func AdminHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !Let(c, AdminUser) {
			return echo.ErrUnauthorized
		}

		return next(c)
	}
}

// UserHandler forces that the next handler function to be called by a valid user
func UserHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := GetUser(c)
		if err != nil {
			return echo.ErrUnauthorized
		}

		if ok, _ := services.Expired(user.ID); ok {
			return echo.ErrUnauthorized
		}

		return next(c)
	}
}

// Let checks permission for doing an action
func Let(c echo.Context, args ...string) bool {
	user, err := GetUser(c)
	if err != nil {
		return false
	}

	if len(args) == 1 && args[0] == AdminUser {
		return user.IsAdmin()
	} else if len(args) == 2 {
		return services.HasDomPerm(user, args[0], args[1])
	} else if len(args) == 3 {
		return services.HasProdPerm(user, args[0], args[1], args[2])
	} else {
		return false
	}
}

// LetUserStuffs indicates that current user has any permission to carry out
// some user stuffs, like viewing profile, ...
func LetUserStuffs(c echo.Context, id int64) bool {
	user, err := GetUser(c)
	if err != nil {
		return false
	}

	return user.IsAdmin() || user.ID == id
}

func IsAdmin(c echo.Context) bool {
	user, err := GetUser(c)
	if err != nil {
		return false
	}

	return user.IsAdmin()
}

// GetUser get user from JWT token
func GetUser(c echo.Context) (*users.User, error) {
	ctx := c.Get("user")
	if ctx == nil {
		return nil, errors.New("invalid jwt token")
	}

	token := ctx.(*jwt.Token)
	claims := token.Claims.(*tokenClaims)

	if !claims.isAuthToken() {
		return nil, errors.New("invalid token")
	}

	user, err := claims.getUser()
	if err != nil {
		return nil, err
	}

	err = db.SetOnline(user.ID)
	if err != nil {
		log.Println("error while setting user status online: ", err)
	}

	return user, nil
}
