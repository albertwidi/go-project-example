package project

import (
	"fmt"
	"strconv"
	"strings"
)

// parseFlags with comma seperated value spesification
// for example: --log=level=info,file=./somefile.log
func parseFlags(value string) (map[string]interface{}, error) {
	fkv := make(map[string]interface{})
	flags := strings.Split(value, ",")
	for _, flag := range flags {
		kv := strings.Split(flag, "=")
		if len(kv) < 2 {
			return nil, fmt.Errorf("debugflag: flag is not in kv format: %s", flag)
		}

		fkv[kv[0]] = kv[1]
	}

	return fkv, nil
}

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

// String return the value of the flag
func (df *debugFlag) String() string {
	return df.flag
}

// Set string value to debug flag
func (df *debugFlag) Set(value string) error {
	if value == "" {
		return nil
	}

	df.flag = value
	kv, err := parseFlags(value)
	if err != nil {
		return err
	}

	for k, v := range kv {
		switch k {
		case "devserver":
			bint, err := strconv.Atoi(v.(string))
			if err != nil {
				return err
			}
			df.DevServer = bint == 1
		case "testconfig":
			bint, err := strconv.Atoi(v.(string))
			if err != nil {
				return err
			}
			df.TestConfig = bint == 1
		}
	}
	return nil
}

// logFlag holds the log flag structure
// logFlag comma seperated value spec
// for example: --log=level=debug,file=./somefile.log
type logFlag struct {
	flag  string
	File  string
	Level string
}

// String return the value of the flag
func (lf *logFlag) String() string {
	return lf.flag
}

func (lf *logFlag) Set(value string) error {
	if value == "" {
		return nil
	}

	lf.flag = value
	kv, err := parseFlags(value)
	if err != nil {
		return err
	}

	var ok bool
	for k, v := range kv {
		switch k {
		case "level":
			lf.Level, ok = v.(string)
			if !ok {
				return fmt.Errorf("logflag: expcect string value for level flag, got %v", v)
			}
		case "file":
			lf.File, ok = v.(string)
			if !ok {
				return fmt.Errorf("logflag: expcect string value for file flag, got %v", v)
			}
		}
	}
	return nil
}
