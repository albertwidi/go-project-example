package cucumber

import (
	"fmt"
	"os"

	"github.com/cucumber/godog"
)

// FileFeature for cucumber
type FileFeature struct {
}

func (ff *FileFeature) fileShouldExist(filepath string) error {
	_, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	return nil
}

func (ff *FileFeature) fileShouldNotExist(filepath string) error {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return nil
	}
	if err == nil {
		return fmt.Errorf("file: file with filepath: %s still exists", filepath)
	}
	return err
}

func (ff *FileFeature) deleteFile(filepath string) error {
	return os.Remove(filepath)
}

// FeatureContext for file
func (ff *FileFeature) FeatureContext(s *godog.Suite) {
	s.Step(`^I delete file "([^"]*)"$`, ff.deleteFile)
	s.Step(`^the file "([^"]*)" should exist`, ff.fileShouldExist)
	s.Step(`^the file "([^"]*)" should not exist`, ff.fileShouldNotExist)
}
