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

// AnyToString convert type (int, int64, float32, float64, byte, and []bytes) to string
// BEWARE: do not use this function for a very spesific usecase
func AnyToString(n interface{}, p ...int) string {
	var t string
	switch n.(type) {
	case int:
		t = strconv.Itoa(n.(int))
	case int64:
		t = strconv.FormatInt(n.(int64), 10)
	case float32:
		if len(p) > 0 {
			t = strconv.FormatFloat(float64(n.(float32)), 'f', p[0], 64)
		} else {
			t = strconv.FormatFloat(float64(n.(float32)), 'f', -1, 64)
		}
	case float64:
		if len(p) > 0 {
			t = strconv.FormatFloat(n.(float64), 'f', p[0], 64)
		} else {
			t = strconv.FormatFloat(n.(float64), 'f', -1, 64)
		}
	case byte:
		t = string(n.(byte))
	case []byte:
		t = string(n.([]byte))
	case string:
		t = n.(string)
	case bool:
		t = strconv.FormatBool(n.(bool))
	}

	return t
}
