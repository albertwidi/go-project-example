package tempe

import "testing"

// a simple test
func TestReplace(t *testing.T) {
	replacefunc := func(matches [][]byte) (map[string]string, error) {
		kv := make(map[string]string)
		for _, m := range matches {
			k := string(m)
			v := "asd"
			kv[k] = v
		}
		return kv, nil
	}

	s := "${this} is ${a} string"
	sbyte := []byte(s)

	te, err := New("\\${[a-zA-Z0-9/-_--]+}", replacefunc)
	if err != nil {
		t.Error(err)
	}
	out, err := te.ReplaceBytes(sbyte)
	if err != nil {
		t.Error(err)
		return
	}

	sout := string(out)
	if sout != "asd is asd string" {
		t.Error("string unmatch")
		return
	}
}
