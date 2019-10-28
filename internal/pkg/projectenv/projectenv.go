package projectenv

import "os"

var _envName string

var _envList = []string{
	Development,
	Staging,
	Production,
}

// environment list
const (
	Development = "development"
	Staging     = "staging"
	Production  = "production"
)

// Init env package
func Init(envName string) {
	_envName = envName
}

// IsDevelopment return true is environment variable value from _envName is development
func IsDevelopment() bool {
	return os.Getenv(_envName) == Development
}

// IsStaging return true is environment variable value from _envName is staging
func IsStaging() bool {
	return os.Getenv(_envName) == Staging
}

// IsProduction return true is environment variable value from _envName is production
func IsProduction() bool {
	return os.Getenv(_envName) == Production
}
