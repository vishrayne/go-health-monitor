package monit

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	// Underline constant
	Underline = "\n================================\n"
	// NotAvailable constant
	NotAvailable = "N.A"
	// Byte size
	Byte = 1.0
	// KiloByte size
	KiloByte = 1024 * Byte
	// MegaByte size
	MegaByte = 1024 * KiloByte
	// GigaByte size
	GigaByte = 1024 * MegaByte
	// TeraByte size
	TeraByte = 1024 * GigaByte
)

// ==================
// byte converters
// ==================

func parseString(data []byte, err error) string {
	if err != nil {
		return NotAvailable
	}

	return asString(data)
}

func parseInt(data []byte, err error) int {
	if err != nil {
		return -1
	}

	return asInteger(data)
}

// ==================
// string converters
// ==================

func parseInt64(data string, err error) uint64 {
	if err != nil {
		return 0
	}

	return asUInt64(data)
}

func parseFloat(data string, err error) float64 {
	if err != nil {
		return -1
	}

	return asFloat(data)
}

// ==================
// generic converters
// ==================

func asInteger(data []byte) int {
	i, err := strconv.Atoi(asString(data))
	if err != nil {
		return -1
	}

	return i
}

func asUInt64(data string) uint64 {
	i64, err := strconv.ParseInt(data, 10, 0)
	if err != nil {
		return 0
	}

	return uint64(i64)
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

// ==================
// json converter
// ==================

func asPrettyJSON(data interface{}) string {
	jsonD, err := json.MarshalIndent(data, "", "  ")
	dealWithError("json", err)
	return string(jsonD)
}

func asJSON(data interface{}) string {
	jsonD, err := json.Marshal(data)
	dealWithError("json", err)

	fmt.Println(string(jsonD))
	fmt.Println("====")

	rawJSON := json.RawMessage(jsonD)
	bytes, err := rawJSON.MarshalJSON()
	dealWithError("json", err)

	fmt.Println(string(bytes))

	return string(bytes)
}

// ==================
// size converter
// ==================

func asHumanBytes(bytes uint64) string {
	unit := ""
	value := float32(bytes)
	switch {
	case bytes >= TeraByte:
		unit = "T"
		value = value / TeraByte
	case bytes >= GigaByte:
		unit = "G"
		value = value / GigaByte
	case bytes >= MegaByte:
		unit = "M"
		value = value / MegaByte
	case bytes >= KiloByte:
		unit = "K"
		value = value / KiloByte
	case bytes >= Byte:
		unit = "B"
	case bytes == 0:
		return "0"
	}
	stringValue := fmt.Sprintf("%.1f", value)
	stringValue = strings.TrimSuffix(stringValue, ".0")
	return fmt.Sprintf("%s%s", stringValue, unit)
}
