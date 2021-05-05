package util

import (
	"strconv"
)


func Encode32(number int64) string {
	return strconv.FormatInt(number,32)
}