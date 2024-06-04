package config

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var durationUnits = map[string]time.Duration{
	"":             time.Second,
	"n":            time.Nanosecond,
	"ns":           time.Nanosecond,
	"u":            time.Microsecond,
	"us":           time.Microsecond,
	"µ":            time.Microsecond,
	"µs":           time.Microsecond,
	"ms":           time.Millisecond,
	"milli":        time.Millisecond,
	"millis":       time.Millisecond,
	"millisecond":  time.Millisecond,
	"milliseconds": time.Millisecond,
	"s":            time.Second,
	"sec":          time.Second,
	"secs":         time.Second,
	"second":       time.Second,
	"seconds":      time.Second,
	"m":            time.Minute,
	"min":          time.Minute,
	"mins":         time.Minute,
	"minute":       time.Minute,
	"minutes":      time.Minute,
	"h":            time.Hour,
	"hr":           time.Hour,
	"hrs":          time.Hour,
	"hour":         time.Hour,
	"hours":        time.Hour,
	"d":            24 * time.Hour,
	"day":          24 * time.Hour,
	"days":         24 * time.Hour,
	"w":            7 * 24 * time.Hour,
	"wk":           7 * 24 * time.Hour,
	"wks":          7 * 24 * time.Hour,
	"week":         7 * 24 * time.Hour,
	"weeks":        7 * 24 * time.Hour,
}

var fileSizes = map[string]int64{
	"":    1,
	"b":   1,
	"kb":  1000,
	"mb":  1000 * 1000,
	"gb":  1000 * 1000 * 1000,
	"tb":  1000 * 1000 * 1000 * 1000,
	"kib": 1024,
	"mib": 1024 * 1024,
	"gib": 1024 * 1024 * 1024,
	"tib": 1024 * 1024 * 1024 * 1024,
}

func splitNumUnit(s string) (string, string) {
	// read until the first non-numeric digit (default is the full string)
	index := len(s)
	for i, char := range s {
		if char != '.' && char != '-' && !unicode.IsDigit(char) {
			index = i
			break
		}
	}

	num, unit := s[:index], s[index:]
	unit = strings.TrimSpace(unit)
	unit = strings.ToLower(unit)

	return num, unit
}

func parseFileSize(s string) (int64, error) {
	num, unit := splitNumUnit(s)

	// parse the number
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0, err
	}

	unitValue, ok := fileSizes[unit]
	if !ok {
		return 0, errors.New("invalid unit")
	}

	return int64(f * float64(unitValue)), nil
}
