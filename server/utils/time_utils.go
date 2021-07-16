package utils

import "time"

func GetTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
