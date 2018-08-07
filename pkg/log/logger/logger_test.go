package logger

import "testing"

func TestNewDefault(t *testing.T) {
	l := Default()
	if l == nil {
		t.Error("default logger is nil")
	}
}
