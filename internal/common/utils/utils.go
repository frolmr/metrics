package utils

import "strconv"

func IntToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func FloatToString(value float64) string {
	return strconv.FormatFloat(value, 'g', -1, 64)
}

func StringToInt(stringValue string) (int64, error) {
	if value, err := strconv.ParseInt(stringValue, 10, 64); err != nil {
		return 0, err
	} else {
		return value, nil
	}
}

func StringToFloat(stringValue string) (float64, error) {
	if value, err := strconv.ParseFloat(stringValue, 64); err != nil {
		return 0, err
	} else {
		return value, nil
	}
}
