package cucumber

import (
	"os"
	"testing"
)

func TestFileShouldExist(t *testing.T) {
	file := "file_test.go"
	ff := FileFeature{}
	if err := ff.fileShouldExist(file); err != nil {
		t.Error(err)
		return
	}
}

func TestFileShouldNotExist(t *testing.T) {
	file := "file_should_never_exist"
	ff := FileFeature{}
	if err := ff.fileShouldNotExist(file); err != nil {
		t.Error(err)
		return
	}
}

func TestDeleteFile(t *testing.T) {
	fileName := "test.txt"
	if _, err := os.Create(fileName); err != nil {
		t.Error(err)
		return
	}

	ff := FileFeature{}
	if err := ff.fileShouldExist(fileName); err != nil {
		t.Error(err)
		return
	}

	if err := ff.deleteFile(fileName); err != nil {
		t.Error(err)
		return
	}

	if err := ff.fileShouldNotExist(fileName); err != nil {
		t.Error(err)
		return
	}
}
