package services

import (
	"github.com/boof/umg/rbac/access"
	"github.com/boof/umg/rest_errors"
)

func AddExpire(expire *access.Expire) rest_errors.Error {
	if err := expire.Save(); err != nil {
		return rest_errors.NewNotAcceptableError(err.Error())
	}

	return nil
}

func EditExpire(expire *access.Expire) rest_errors.Error {
	if err := expire.Update(); err != nil {
		return rest_errors.NewNotAcceptableError(err.Error())
	}

	return nil
}

func DelExpire(userID int64) rest_errors.Error {
	if err := (&access.Expire{UserID: userID}).RemoveByUserID(); err != nil {
		return rest_errors.NewNotAcceptableError(err.Error())
	}

	return nil
}

// AccessExpired indicates that given user access to portal is expired or not
func Expired(userID int64) (bool, rest_errors.Error) {
	return access.Expired(userID)
}
