package date

import "time"

// SelectTime 凌晨至早8点时间确认
func SelectTime() bool {
	now := time.Now().UTC()
	mornDate, _ := time.Parse("2006-01-02", now.Format("2006-01-02"))
	startTime := mornDate.Add(time.Hour * 16)
	stopTime := mornDate.Add(time.Hour * 24)
	if now.After(startTime) && now.Before(stopTime) {
		return true
	}
	return false
}
