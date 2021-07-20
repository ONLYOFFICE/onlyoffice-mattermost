package utils

import "time"

func GetTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetTime(timestamp int64) time.Time {
	micro := timestamp / int64(time.Microsecond)
	remainder := (timestamp % int64(time.Microsecond)) * int64(time.Millisecond)
	dateTime := time.Unix(micro, remainder)

	return dateTime
}
