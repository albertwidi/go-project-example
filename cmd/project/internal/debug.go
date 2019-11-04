package project

import (
	"fmt"
	"strconv"
	"strings"
)

// debugFlag holds the debug flag structure
// spec:
// debug flag is comma seprated value with 'equal/=' as operator
// the value can be identified as '1' as 'true' and '0' as 'false'
// for example: debug=devserver=1,testconfig=1
type debugFlag struct {
	flag       string
	DevServer  bool
	TestConfig bool
}

// String return the value of flag
func (df *debugFlag) String() string {
	return df.flag
}

// Set string value to debug flag
func (df *debugFlag) Set(value string) error {
	if value == "" {
		return nil
	}

	df.flag = value
	flags := strings.Split(value, ",")
	for _, flag := range flags {
		kv := strings.Split(flag, "=")

		if len(kv) < 2 {
			return fmt.Errorf("debugflag: flag is not in kv format: %s", flag)
		}

		var err error
		switch kv[0] {
		case "devserver":
			df.DevServer, err = strconv.ParseBool(kv[1])
		case "testconfig":
			df.TestConfig, err = strconv.ParseBool(kv[1])
		}

		if err != nil {
			return err
		}
	}
	return nil
}
