package helper

import "time"

func TimeNowUTC() time.Time {
	return time.Now().UTC()
}

func TimeNowEpochUtc() int64 {
	return TimeNowUTC().Unix()
}
