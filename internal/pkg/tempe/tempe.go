// tempe is template replacer for replacing regex matching string
// it has Replacer function to return a KV map for matching string

package tempe

import (
	"bytes"
	"os"
	"regexp"
)

// ReplaceFunc for tempe
type ReplaceFunc func(matches [][]byte) (map[string]string, error)

// Tempe struct
type Tempe struct {
	regex   *regexp.Regexp
	replcer ReplaceFunc
}

// EnvVarPattern define environment variable pattern with ${MY_ENV_VAR}
const EnvVarPattern = "\\${[a-zA-Z0-9/-_--]+}"

// EnvVarReplacerFunc for replacing environment variable with the regex
var EnvVarReplacerFunc = func(matches [][]byte) (map[string]string, error) {
	kv := make(map[string]string)
	for _, m := range matches {
		k := string(m)
		v := os.Getenv(k[2 : len(k)-1])
		kv[k] = v
	}
	return kv, nil
}

// New tempe object
func New(regex string, replacer ReplaceFunc) (*Tempe, error) {
	t := Tempe{
		replcer: replacer,
	}
	rxp, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	t.regex = rxp

	return &t, nil
}

// ReplaceBytes the string
func (t *Tempe) ReplaceBytes(in []byte) ([]byte, error) {
	matches := t.regex.FindAll(in, -1)
	if len(matches) == 0 {
		return in, nil
	}

	kv, err := t.replcer(matches)
	if err != nil {
		return nil, err
	}

	for k, v := range kv {
		in = bytes.Replace(in, []byte(k), []byte(v), -1)
	}
	return in, nil
}
