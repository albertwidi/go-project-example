package project

type arrayFlags []string

// String return string implementation of array flags
func (af arrayFlags) String() string {
	return ""
}

// Set for append the value of arrayFlags
func (af arrayFlags) Set(value string) error {
	af = append(af, value)
	return nil
}

// Flags of project
type Flags struct {
	DebugMode         bool
	TimeZone          string
	ConfigurationFile string
	EnvironmentFiles  arrayFlags
}

// Run the project
func Run(f Flags) error {
	return nil
}
