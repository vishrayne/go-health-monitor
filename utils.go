package monit

import (
	"strconv"
	"strings"
)

func asInteger(data []byte) int {
	i, err := strconv.Atoi(asString(data))
	if err != nil {
		return -1
	}

	return i
}

func asString(data []byte) string {
	return asSafeString(string(data))
}

func asSafeString(data string) string {
	stringData := strings.TrimSpace(string(data))
	return strings.TrimRight(stringData, "\n")
}

func asFloat(data string) float64 {
	safeString := asSafeString(data)
	floatVal, err := strconv.ParseFloat(safeString, 64)
	if err != nil {
		return -1
	}

	return floatVal
}
