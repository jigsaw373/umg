package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/boof/umg/auth"
	"github.com/boof/umg/db"
	"github.com/boof/umg/email"
	"github.com/boof/umg/rbac/users"
	"github.com/boof/umg/services"
	"github.com/boof/umg/util/datetime"
	"github.com/boof/umg/util/response"
)

func Login(c echo.Context) error {
	user := new(users.User)
	if err := c.Bind(user); err != nil {
		return response.BadReq(c, "bad request")
	}

	user, logErr := services.Login(user.Username, user.Password)
	if logErr != nil {
		return logErr.Echo(c)
	}

	token, refresh, err := auth.CreateTokens(user)
	if err != nil {
		return response.InternalErr(c, "unable to create token")
	}

	user.LastLogin = datetime.NowInEasternCanada()
	err = user.UpdateLastLogin()
	if err != nil {
		fmt.Println("Unable to save last login: ", err)
	}

	policies, _ := services.GetPolices(user)
	domains, _ := services.GetUserDomains(user.ID)

	return c.JSON(http.StatusOK, echo.Map{"result": echo.Map{
		"token":         token,
		"refresh_token": refresh,
		"username":      user.Username,
		"user_id":       user.ID,
		"admin":         user.IsAdmin(),
		"policies":      policies,
		"domains":       domains,
	}})
}

func ResetPasswordByEmail(c echo.Context) error {
	type Req struct {
		Email string `json:"email"`
	}

	req := new(Req)
	if err := c.Bind(req); err != nil || req.Email == "" {
		return response.BadReq(c, "bad request")
	}

	user, err := (&users.User{Email: req.Email}).GetByEmail()
	if err != nil {
		return response.NotFound(c, "user not found")
	}

	if ok, _ := services.Expired(user.ID); ok {
		return errors.New("Your access time is expired!")
	}

	token, err := db.GenResetPassToken(user.ID, 30)
	if err != nil {
		log.Printf("unable to generate reset password token: %v \n", err)
		return response.InternalErr(c, "internal server error")
	}

	url := "https://portal.edgecomenergy.ca/reset-password/" + token
	err = email.SendResetEmail(user.ID, user.Name, url, user.Email)

	return response.Done(c)
}

// SendWelcomeAndReset sends a welcome email that
// contains a reset password link
func SendWelcomeAndReset(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadReq(c, "bad request")
	}

	user, err := (&users.User{ID: id}).GetByID()
	if err != nil {
		return response.NotFound(c, "user not found")
	}

	token, err := db.GenResetPassToken(user.ID, 48*60)
	if err != nil {
		log.Printf("unable to generate reset password token: %v \n", err)
		return response.InternalErr(c, "internal server error")
	}

	url := "https://portal.edgecomenergy.ca/reset-password/" + token
	err = email.SendWelcomeAndResetEmail(user.ID, user.Name, user.Username, url, user.Email)

	return response.Done(c)
}

func GetUserEmailHistory(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadReq(c, "bad request")
	}

	history, getErr := email.GetUserEmailHistory(id)
	if getErr != nil {
		return getErr.Echo(c)
	}

	if history == nil || len(history) == 0 {
		return response.NotContent(c)
	}

	return response.OK(c, history)
}

func ValidateResetPassToken(c echo.Context) error {
	type Req struct {
		Token string `json:"token"`
	}

	req := new(Req)
	if err := c.Bind(req); err != nil || req.Token == "" {
		return response.BadReq(c, "bad request")
	}

	userID, err := db.ValidateResetPassToken(req.Token)
	if err != nil || userID < 0 {
		return response.NotFound(c, "invalid token")
	}

	return c.JSON(http.StatusOK, echo.Map{"result": echo.Map{"valid": true}})
}

func ChangePassword(c echo.Context) error {
	type Req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	req := new(Req)
	if err := c.Bind(req); err != nil || req.Token == "" {
		return response.BadReq(c, "bad request")
	}

	userID, err := db.ValidateResetPassToken(req.Token)
	if err != nil || userID < 0 {
		return response.NotFound(c, "invalid token")
	}

	// revoke token
	db.RevokeResetPassToken(req.Token)

	err = (&users.User{ID: userID}).ChangePassword(req.Password)
	if err != nil {
		return response.NotAcceptable(c, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{})
}

// ChangeUserPassword used by admin for changing users password
func ChangeUserPassword(c echo.Context) error {
	type Req struct {
		UserID   int64  `json:"user_id"`
		Password string `json:"password"`
	}

	req := new(Req)
	if err := c.Bind(req); err != nil {
		return response.BadReq(c, "bad request")
	}

	user := &users.User{ID: req.UserID}
	if _, err := user.GetByID(); err != nil {
		return response.BadReq(c, "invalid user")
	}

	err := user.ChangePassword(req.Password)
	if err != nil {
		return response.NotAcceptable(c, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{})
}
