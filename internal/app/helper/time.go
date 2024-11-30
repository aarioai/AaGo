package helper

import (
	"time"

	"github.com/aarioai/AaGo/internal/app/atype"
)

type Time struct {
	Location *time.Location
}

func NewTime(loc *time.Location) *Time {
	if loc == nil {
		loc = time.Local
	}
	return &Time{Location: loc}
}
func (t *Time) Now() time.Time {
	return time.Now().In(t.Location)
}

func (t *Time) FmtNow() string {
	return t.Now().Format("2006-01-02 15:04:05")
}
func (t *Time) FmtDate() string {
	return t.Now().Format("2006-01-02")
}
func (t *Time) FmtTime() string {
	return t.Now().Format("15:04:05")
}

func (t *Time) ParseDatetime(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", s, t.Location)
}
func (t *Time) ParseDate(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", s, t.Location)
}
func (t *Time) Datetime() atype.Datetime {
	return atype.ToDatetime(t.Now())
}
func (t *Time) Date() atype.Date {
	return atype.ToDate(t.Now())
}
