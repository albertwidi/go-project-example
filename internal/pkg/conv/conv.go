package conv

import "strconv"

//Float64ToString return string form of float64
func Float64ToString(f float64, precission int) string {
	return strconv.FormatFloat(f, 'f', precission, 64)
}

//Int64ToString return string form of int64
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// StringToInt64 convert string to int64
func StringToInt64(i string) (int64, error) {
	n, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		return 0, err
	}

	return n, nil
}
