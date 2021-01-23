package password

import (
	"errors"
	"fmt"

	"github.com/alexandrevicenzi/unchained"
	"github.com/go-playground/validator/v10"

	"github.com/boof/umg/settings"
)

var (
	// use a single instance of Validate, it caches struct info
	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

func HashPassword(pass string) (string, error) {
	return unchained.MakePassword(pass, settings.PasswordSalt, "default")
}

func IsValidPass(pass, hash string) bool {
	valid, _ := unchained.CheckPassword(pass, hash)
	return valid
}

func ValidateUsername(username string) error {
	if errs := validate.Var(username, "min=1"); errs != nil {
		return errors.New("username must contain at least one character")
	}

	if errs := validate.Var(username, "max=64"); errs != nil {
		return errors.New("username must be shorter than or equal 64 characters")
	}

	hasUnderline := false
	hasHyphen := false

	for i, ch := range username {
		if i == 0 {
			if (ch < '0' || ch > '9') && (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') {
				return errors.New("username should start with a letter or number")
			}

			continue
		}

		if (ch < '0' || ch > '9') && (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && ch != '.' &&
			ch != '@' && ch != '_' && ch != '-' {
			return errors.New("username can only contain letters, numbers and one of the \"_ . @ -\"")
		}

		if ch == '@' {
			if err := validate.Var(username, "email,required"); err != nil {
				return errors.New("username can contain @ only when it's a valid email address")
			}
		}

		if ch == '-' {
			if hasUnderline {
				return errors.New("username can't contain both '_' and '-'")
			}

			if username[i-1] == '.' {
				return errors.New("'.-' is not valid in the username")
			}

			hasHyphen = true
		}

		if ch == '_' {
			if hasHyphen {
				return errors.New("username can't contain both '-' and '_'")
			}

			if username[i-1] == '.' {
				return errors.New("'._' is not valid in the username")
			}

			hasUnderline = true
		}

		if ch == '-' && username[i-1] == '-' {
			return errors.New("username can't contain two consecutive '-'")
		}

		if ch == '.' {
			if username[i-1] == '.' {
				return errors.New("username can't contain two consecutive '.'")
			}

			if username[i-1] == '-' {
				return errors.New("'-.' is not valid in the username")
			}

			if username[i-1] == '_' {
				return errors.New("'_.' is not valid in the username")
			}
		}

		if i == len(username)-1 {
			if ch == '.' || ch == '-' || ch == '@' {
				return fmt.Errorf("last character of the username can't be '%c'", ch)
			}
		}
	}

	return nil
}
