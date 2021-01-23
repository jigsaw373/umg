package access

import (
	"errors"
	"log"
	"time"

	"github.com/boof/umg/db"
	"github.com/boof/umg/rbac/users"
	"github.com/boof/umg/rest_errors"
	"github.com/boof/umg/settings"
	"github.com/boof/umg/util/datetime"
)

func init() {
	db.Sync(new(Expire))
}

func (req *ExpireReq) GetDate() (time.Time, error) {
	return time.Parse(settings.DTLayout, req.Date)
}

// Save saves new access expire
func (a *Expire) Save() error {
	err := a.validateForInsert()
	if err != nil {
		return err
	}

	// remove id
	a.ID = 0

	_, err = db.Engine.Insert(a)
	return err
}

// Update updates expired time for an user
func (a *Expire) Update() error {
	access, err := a.GetByUserID()
	if err != nil {
		return a.Save()
	}

	if _, err := db.Engine.Id(access.ID).Cols("expire_at").Update(a); err != nil {
		log.Print("error while updating expired time: ", err)
		return errors.New("database error")
	}

	return nil
}

func (a *Expire) RemoveByID() error {
	_, err := db.Engine.Id(a.ID).Delete(&Expire{})
	return err
}

func (a *Expire) RemoveByUserID() error {
	_, err := db.Engine.Delete(&Expire{UserID: a.UserID})
	return err
}

func (a *Expire) validateForInsert() error {
	user, err := (&users.User{ID: a.UserID}).GetByID()
	if err != nil {
		return errors.New("user not found")
	}

	if user.IsAdmin() {
		return errors.New("you can't set access expiration time on the admin user")
	}

	if _, err := a.GetByUserID(); err == nil {
		return errors.New("there is expire time for the current user")
	}

	return nil
}

// GetByID returns an access expire with the given id
func (a *Expire) GetByID() (*Expire, error) {
	access := &Expire{ID: a.ID}
	if has, err := db.Engine.Get(access); !has || err != nil {
		return nil, errors.New("access expire not found")
	}

	return access, nil
}

// GetByID returns an access expire with the given userID
func (a *Expire) GetByUserID() (*Expire, error) {
	access := &Expire{UserID: a.UserID}
	if has, err := db.Engine.Get(access); !has || err != nil {
		return nil, errors.New("there is no access expire for the given user")
	}

	return access, nil
}

// Expired indicates that given user access to portal is expired or not
func Expired(userID int64) (bool, rest_errors.Error) {
	access := &Expire{UserID: userID}

	has, err := db.Engine.Get(access)
	if err != nil {
		return true, rest_errors.NewNotFoundError("There is no access expire for the given user")
	}

	return has && access.ExpireAt.Before(datetime.NowInEasternCanada()), nil
}
