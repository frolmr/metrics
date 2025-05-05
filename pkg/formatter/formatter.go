package formatter

import (
	"errors"
	"strconv"
	"strings"
)

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

func CheckSchemeFormat(scheme string) error {
	if scheme != "http" && scheme != "https" && scheme != "grpc" {
		return errors.New("bad scheme must be http/https or grpc")
	}
	return nil
}

func CheckAddrFormat(addr string) error {
	hp := strings.Split(addr, ":")

	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}

	_, err := strconv.Atoi(hp[1])

	if err != nil {
		return err
	}

	return nil
}
