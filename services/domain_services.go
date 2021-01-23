package services

import (
	"github.com/boof/umg/rbac/domains"
	"github.com/boof/umg/rest_errors"
)

func GetAllDomains() ([]domains.Domain, rest_errors.Error) {
	res, err := domains.GetAllDomains()
	if err != nil {
		return nil, rest_errors.NewNotFoundError(err.Error())
	}

	return res, nil
}
