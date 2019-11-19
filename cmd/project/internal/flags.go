package project

import (
	"flag"
	"fmt"
	"regexp"
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

	// find pattern -{flag_name}={flag_value}
	regex, err := regexp.Compile("-[a-zA-Z0-9]+=[a-zA-Z0-9]+")
	if err != nil {
		return err
	}

	df.flag = value
	fs := flag.CommandLine
	fs.BoolVar(&df.DevServer, "devserver", false, "for activating dev server")
	fs.BoolVar(&df.TestConfig, "testconfig", false, "for testing the project configuration")
	fs.Parse(regex.FindAllString(value, -1))
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
