package model

import (
	"database/sql/driver"
	"fmt"
	"pkg/utils"
	"time"
)

type HHTime time.Time

const TimeLayout = "2006-01-02 15:04:05-07:00"

func (hht *HHTime) Scan(value any) error {
	switch v := value.(type) {
	case string:
		t, err := time.Parse(TimeLayout, v)
		if err != nil {
			return err
		}
		*hht = HHTime(t)
		return nil
	case []byte:
		t, err := time.Parse(TimeLayout, string(v))
		if err != nil {
			return err
		}
		*hht = HHTime(t)
		return nil
	case time.Time:
		*hht = HHTime(v)
		return nil
	default:
		return fmt.Errorf("HHTime.Scan: cannot scan type %T into HHTime", v)
	}
}

func (hht HHTime) Value() (driver.Value, error) {
	return time.Time(hht), nil
}

type Resume struct {
	ID           string
	AlternateUrl string
	Title        string
	CreatedAt    HHTime
	UpdatedAt    HHTime
	AlternateURL string
	IsScheduled  int
}

type Scheduler struct {
	UserID      string
	ResumeID    string
	ResumeTitle string
	Timestamp   string
	Error       string
}

type JoinedScheduler struct {
	UserID       string
	AccessToken  string
	RefreshToken string
	ResumeID     string
	ResumeTitle  string
}

type Token struct {
	AccessToken  string
	RefreshToken string
}

func (t *Token) Encrypt() error {
	var err error

	if t.AccessToken, err = utils.Encrypt(t.AccessToken); err != nil {
		return err
	}

	if t.RefreshToken, err = utils.Encrypt(t.RefreshToken); err != nil {
		return err
	}

	return nil
}

func (t *Token) Decrypt() error {
	var err error

	if t.AccessToken, err = utils.Decrypt(t.AccessToken); err != nil {
		return err
	}

	if t.RefreshToken, err = utils.Decrypt(t.RefreshToken); err != nil {
		return err
	}

	return nil
}

type User struct {
	ID         string
	FirstName  string
	LastName   string
	MiddleName string
}

type Session struct {
	ID        string
	ExpiresAt time.Time
	UserID    string
}
