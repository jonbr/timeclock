package utils

import (
	"strconv"

	"github.com/gookit/goutil/dump"
)

// CastStringToUint takes in map parameters as strings and returns
// there uint representation as a slice.
func CastStringToUint(pathParams map[string]string) ([]uint, error) {
	var paramsToUint []uint
	for _, value := range pathParams {
		u64, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return nil, err
		}

		paramsToUint = append(paramsToUint, uint(u64))
	}

	return paramsToUint, nil
}

func CastStringToInt(str string) (int, error) {
	i, err := strconv.Atoi(str)
	dump.P(err)
	if err != nil && str != "" {
		return 0, err
	}

	return i, nil
}

/*func CastStringToInt(models.Measurement, url.Values) error {
	return nil
}*/

func CastParamToUint(pathParam string) (uint, error) {
	//var paramsToUint []uint
	//for _, value := range pathParams {
	u64, err := strconv.ParseUint(pathParam, 10, 32)
	if err != nil {
		return 0, err
	}

	//paramsToUint = append(paramsToUint, uint(u64))
	//}

	return uint(u64), nil
}
