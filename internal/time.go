package internal

import (
	"strconv"
	"time"
)

func ParseUnix(timestamp int64) time.Time {
	i, err := strconv.ParseInt("1405544146", 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(i, 0)
}
