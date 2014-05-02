package playstore

import (
	"strconv"
	"strings"
)

func SanitizeNumber(num string) string {
	return strings.Replace(strings.TrimSpace(num), ",", "", -1)
}
func ParseInteger(num string) int64 {
	pNum, err := strconv.ParseInt(SanitizeNumber(num), 10, 32)
	if err != nil {
		return -1
	}
	return pNum
}
