package accessDB

import (
	"errors"
	"strconv"
)

func ValidateID(id string) error {
	parsedID, err := strconv.Atoi(id)
	if err != nil || parsedID <= 0 {
		return errors.New("not correct one of parameters")
	}
	return nil
}

func ValidateLimitOffset(id string) error {
	parsedID, err := strconv.Atoi(id)
	if err != nil || parsedID < 0 {
		return errors.New("not correct one of parameters")
	}
	return nil
}

func ValidateLastRevision(flag string) error {
	if flag != "false" && flag != "true" {
		return errors.New("not correct one of parameters")
	}
	return nil
}
