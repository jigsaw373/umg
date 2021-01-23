package controller

import (
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/boof/umg/rbac/access"
	"github.com/boof/umg/services"
	"github.com/boof/umg/util/response"
)

func AddAccessExpire(c echo.Context) error {
	req := new(access.ExpireReq)
	if err := c.Bind(req); err != nil {
		return response.BadReq(c, "bad request")
	}

	date, err := req.GetDate()
	if err != nil {
		return response.BadReq(c, "invalid date")
	}

	expire := &access.Expire{UserID: req.UserID, ExpireAt: date}

	if err := services.AddExpire(expire); err != nil {
		return err.Echo(c)
	}

	return response.Created(c, echo.Map{"id": expire.ID})
}

func EditAccessExpire(c echo.Context) error {
	req := new(access.ExpireReq)
	if err := c.Bind(req); err != nil {
		return response.BadReq(c, "bad request")
	}

	date, err := req.GetDate()
	if err != nil {
		return response.BadReq(c, "invalid date")
	}

	expire := &access.Expire{UserID: req.UserID, ExpireAt: date}

	if err := services.EditExpire(expire); err != nil {
		return err.Echo(c)
	}

	return response.Done(c)
}

func DelAccessExpire(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadReq(c, "bad request")
	}

	if err := services.DelExpire(id); err != nil {
		return err.Echo(c)
	}

	return response.Done(c)
}
