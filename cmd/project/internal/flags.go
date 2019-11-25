package project

import (
	"flag"
	"regexp"
)

// debugFlag holds the debug flag structure
// spec:
// debug flag is a flag.Parse value
// for example: -debug=-server=1-testconfig=1
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
	return fs.Parse(regex.FindAllString(value, -1))
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
