package email

import (
	"encoding/json"
	"time"

	"github.com/boof/umg/db"
	"github.com/boof/umg/rbac/users"
	"github.com/boof/umg/rest_errors"
	"github.com/boof/umg/settings"
	"github.com/boof/umg/util/datetime"
)

type History struct {
	ID     int64     `xorm:"pk not null autoincr 'id'"`
	UserID int64     `xorm:"not null 'user_id'"`
	Status string    `xorm:"not null 'status'"`
	Date   time.Time `xorm:"not null 'date'"`
}

func init() {
	db.Sync(new(History))
}

func (h *History) Save() rest_errors.Error {
	if err := h.InsertValidate(); err != nil {
		return err
	}

	// remove id
	h.ID = 0

	_, err := db.Engine.Insert(h)
	if err != nil {
		return rest_errors.NewInternalServerError("Database error", err)
	}

	return nil
}

func (h *History) InsertValidate() rest_errors.Error {
	_, err := (&users.User{ID: h.UserID}).GetByID()
	if err != nil {
		return rest_errors.NewBadRequestError("Invalid user")
	}

	if h.Status == "" {
		return rest_errors.NewBadRequestError("Empty status")
	}

	return nil
}

func (h *History) MarshalJSON() ([]byte, error) {
	type HistoryJSON struct {
		Date   string `json:"date"`
		Status string `json:"status"`
	}

	history := &HistoryJSON{
		Date:   h.Date.Format(settings.DTLayout),
		Status: h.Status,
	}

	return json.Marshal(history)
}

func AddHistory(userId int64, status string) rest_errors.Error {
	history := &History{
		UserID: userId,
		Status: status,
		Date:   datetime.NowInEasternCanada(),
	}

	return history.Save()
}

func GetUserEmailHistory(userID int64) ([]History, rest_errors.Error) {
	var history []History

	if err := db.Engine.Where("user_id=?", userID).Find(&history); err != nil {
		return nil, rest_errors.NewInternalServerError("Database error", err)
	}

	return history, nil
}
