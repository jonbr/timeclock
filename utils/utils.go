package utils

import (
	"strconv"
)

func CastStringToUint(ID string) (uint, error) {
	uintId, err := strconv.ParseUint(ID, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(uintId), nil
}