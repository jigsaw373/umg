package auth

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4/middleware"
	"github.com/boof/umg/db"
	"github.com/boof/umg/rbac/users"
	"github.com/boof/umg/settings"
)

type tokenClaims struct {
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

type refreshTokenClaims struct {
	jwt.StandardClaims
}

func (claim *tokenClaims) getUser() (*users.User, error) {
	return getUserFromSubject(claim.Subject)
}

func (claim *tokenClaims) isAuthToken() bool {
	return claim.Audience == ""
}

func (claim *refreshTokenClaims) getUser() (*users.User, error) {
	return getUserFromSubject(claim.Subject)
}

func (claim *refreshTokenClaims) isRefreshToken() bool {
	return claim.Audience == settings.RefreshTokenAudience
}

func CreateTokens(user *users.User) (string, string, error) {
	// Set token claims
	tc := &tokenClaims{
		Admin: user.IsAdmin(),
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.FormatInt(user.ID, 10),
			ExpiresAt: time.Now().Add(settings.JWTExpiry * time.Minute).Unix(),
		},
	}

	// set refresh token claims
	rc := &refreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.FormatInt(user.ID, 10),
			Audience:  settings.RefreshTokenAudience,
			ExpiresAt: time.Now().Add(settings.JWTRefreshExpiry * time.Minute).Unix(),
		},
	}

	// Create tokens
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tc)
	rawRefresh := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)

	// Generate encoded tokens
	token, err := rawToken.SignedString([]byte(settings.JWTSecret))
	if err != nil {
		return "", "", err
	}

	refresh, err := rawRefresh.SignedString([]byte(settings.JWTSecret))
	if err != nil {
		return "", "", err
	}

	return token, refresh, nil
}

func GetUserFromToken(token string) (*users.User, error) {
	claims := &tokenClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(settings.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

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

func GetMiddlewareConfig() middleware.JWTConfig {
	return middleware.JWTConfig{
		Claims:     &tokenClaims{},
		SigningKey: []byte(settings.JWTSecret),
	}
}

func getUserFromSubject(subject string) (*users.User, error) {
	id, err := strconv.ParseInt(subject, 10, 64)
	if err != nil {
		return nil, err
	}

	return (&users.User{ID: id}).GetByID()
}
