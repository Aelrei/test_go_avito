package accessDB

import (
	"errors"
	"strconv"
)

func ValidateID(id string) (int, error) {
	parsedID, err := strconv.Atoi(id)
	if err != nil || parsedID <= 0 {
		return 0, errors.New("not correct one of parameters")
	}
	return parsedID, nil
}

func ValidateLimitOffset(id string) (int, error) {
	parsedID, err := strconv.Atoi(id)
	if err != nil || parsedID < 0 {
		return 0, errors.New("not correct one of parameters")
	}
	return parsedID, nil
}

func ValidateLastRevision(flag string) (bool, error) {
	if flag != "false" && flag != "true" {
		return false, errors.New("not correct one of parameters")
	}
	return true, nil
}
