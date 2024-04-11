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
