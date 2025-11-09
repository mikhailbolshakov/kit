package vault

import (
	"fmt"
	"time"
)

func ToDateFormat(t time.Time) string {
	return t.Format("2006-01-02T15:04:05Z")
}

func ToTTLHours(hours int) string {
	return fmt.Sprintf("%dh", hours)
}
