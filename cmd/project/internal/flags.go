package project

import (
	"fmt"
	"strconv"
	"strings"
)

// parseFlags with comma seperated value spec
// for example: --log=level=info,file=./somefile.log
func parseFlags(value string) (map[string]string, error) {
	fkv := make(map[string]string)
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
// for example: debug=server=1,testconfig=1
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
		case "server":
			bint, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			df.DevServer = bint == 1
		case "testconfig":
			bint, err := strconv.Atoi(v)
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
	Color bool
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

	for k, v := range kv {
		switch k {
		case "level":
			lf.Level = v
		case "file":
			lf.File = v
		case "color":
			bint, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			lf.Color = bint == 1
		}
	}
	return nil
}

type envFileFlag struct {
	flag     string
	envFiles []string
}

func (vf *envFileFlag) String() string {
	return vf.flag
}

func (vf *envFileFlag) Set(value string) error {
	vf.envFiles = append(vf.envFiles, value)
	return nil
}
