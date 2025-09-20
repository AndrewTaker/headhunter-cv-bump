package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type HHTime time.Time

const timeLayout = "2006-01-02 15:04:05-07:00"

func (hht *HHTime) Scan(value any) error {
	switch v := value.(type) {
	case string:
		t, err := time.Parse(timeLayout, v)
		if err != nil {
			return err
		}
		*hht = HHTime(t)
		return nil
	case []byte:
		t, err := time.Parse(timeLayout, string(v))
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

type SqliteUser struct {
	ID         string
	FirstName  string
	LastName   string
	MiddleName string
}

type SqliteResume struct {
	ID           string
	AlternateUrl string
	Title        string
	CreatedAt    HHTime
	UpdatedAt    HHTime
	AlternateURL string
	IsScheduled  int
}

type SqliteToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    uint
	Code         string
	TokenType    string
}

type SqliteScheduler struct {
	UserID      string
	ResumeID    string
	ResumeTitle string
	Timestamp   string
	Error       string
}
