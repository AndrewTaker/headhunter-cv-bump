package headhunter

import (
	"strings"
	"time"
)

type HHTime time.Time

const timeLayout = "2006-01-02 15:04:05-07:00"

func (hht *HHTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02T15:04:05-0700", s)
	if err != nil {
		return err
	}
	*hht = HHTime(t)
	return nil
}

func (hht HHTime) MarshalJSON() ([]byte, error) {
	t := time.Time(hht)
	return []byte(`"` + t.Format(timeLayout) + `"`), nil
}

func (t HHTime) Format(layout string) string {
	return time.Time(t).Format(time.RFC1123)
}

type User struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

type Resume struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	CreatedAt    HHTime `json:"created_at"`
	UpdatedAt    HHTime `json:"updated_at"`
	AlternateURL string `json:"alternate_url"`
}
