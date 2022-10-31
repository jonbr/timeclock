package utils

import (
	//"fmt"
	"strconv"
	//"github.com/gookit/goutil/dump"
)

func CastStringToUint(pathParams map[string]string) ([]uint, error) {
	var paramsToUint []uint
	for _, value := range pathParams {
		u64, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return nil, err
		}

		paramsToUint = append(paramsToUint, uint(u64))
	}

	//dump.P(paramsToUint)

	return paramsToUint, nil
}
