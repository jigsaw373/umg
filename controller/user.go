package controller

import (
	"github.com/boof/umg/util/response"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/boof/umg/auth"
	"github.com/boof/umg/rbac/users"
	"github.com/boof/umg/services"
	"github.com/boof/umg/settings"
	"github.com/boof/umg/util/request"
)

func GetUser(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadReq(c, "bad request")
	}

	if !auth.LetUserStuffs(c, id) {
		return echo.ErrUnauthorized
	}

	user, getErr := services.GetUserByID(id)
	if getErr != nil {
		return getErr.Echo(c)
	}

	return response.OK(c, user)
}

// GetUsers returns all users
func GetUsers(c echo.Context) error {
	format := c.QueryParam("format")

	// pagination parameters
	count, page := request.GetPagination(c)

	order := c.QueryParam("order")
	if order != services.Desc {
		order = services.Asc
	}

	sortBy := c.QueryParam("sort")
	if !services.IsValidSort(sortBy) {
		sortBy = services.ID
	}

	if format == "simple" {
		return GetSimpleUsers(c, count, page, order, sortBy)
	} else if format == "with-role" {
		return GetRoleUsers(c, count, page, order, sortBy)
	} else {
		return response.BadReq(c, "invalid format")
	}
}

func GetSimpleUsers(c echo.Context, count, page int64, order string, sortBy string) error {
	users, pages, err := services.GetSimpleUsers(count, page, order, sortBy)
	if err != nil {
		return err.Echo(c)
	}

	// set pagination header
	c.Response().Header().Add("X-Pagination-Page-Count", strconv.FormatInt(pages, 10))
	return response.OK(c, users)
}

func GetRoleUsers(c echo.Context, count, page int64, order string, sortBy string) error {
	users, pages, err := services.GetRoleUsers(count, page, order, sortBy)
	if err != nil {
		return err.Echo(c)
	}

	// set pagination header
	c.Response().Header().Add("X-Pagination-Page-Count", strconv.FormatInt(pages, 10))
	return response.OK(c, users)
}

func GetUserDomains(c echo.Context) error {
	user, err := auth.GetUser(c)
	if err != nil {
		return echo.ErrUnauthorized
	}

	domIDs, GetErr := services.GetUserDomains(user.ID)
	if GetErr != nil {
		return GetErr.Echo(c)
	}

	return response.OK(c, domIDs)
}

func GetUserProducts(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadReq(c, "bad request")
	}

	user, err := auth.GetUser(c)
	if err != nil {
		return echo.ErrUnauthorized
	}

	prods, GetErr := services.GetUserProducts(user.ID, id)
	if GetErr != nil {
		return GetErr.Echo(c)
	}

	return response.OK(c, prods)
}

func AddUser(c echo.Context) error {
	user := new(users.User)
	if err := c.Bind(user); err != nil {
		return response.BadReq(c, "bad request")
	}

	err := services.AddUser(user)
	if err != nil {
		return err.Echo(c)
	}

	return response.OK(c, echo.Map{"id": user.ID})
}

func AddUserWithRole(c echo.Context) error {
	sendEmail := c.QueryParam("send_email")
	expireAt := c.QueryParam("expire_at")

	user := new(users.User)
	if err := c.Bind(user); err != nil {
		return response.BadReq(c, "bad request")
	}

	var expireTime *time.Time
	if expireAt != "" {
		date, err := time.Parse(settings.DTLayout, expireAt)
		if err != nil {
			return response.BadReq(c, "invalid expiry date")
		}

		expireTime = &date
	}

	err := services.AddUserWithRole(user, sendEmail == "true", expireTime)
	if err != nil {
		return err.Echo(c)
	}

	return response.Created(c, echo.Map{"id": user.ID})
}

func UpdateUser(c echo.Context) error {
	user := new(users.User)
	if err := c.Bind(user); err != nil {
		return response.BadReq(c, "bad request")
	}

	if !auth.LetUserStuffs(c, user.ID) {
		return echo.ErrUnauthorized
	}

	err := services.UpdateUser(user)
	if err != nil {
		return err.Echo(c)
	}

	return response.Done(c)
}

func DelUser(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadReq(c, "bad request")
	}

	delErr := services.DelUser(id)
	if delErr != nil {
		return delErr.Echo(c)
	}

	return response.Done(c)
}
