// the test scenario is using local storage options as its backend
// it will trigger the bucket via init()
package objectstorage_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/albertwidi/go-project-example/internal/pkg/objectstorage"
	"github.com/albertwidi/go-project-example/internal/pkg/objectstorage/local"
)

// File struct
type File struct {
	FileName string
}

// Content struct
type Content struct {
	Body   []byte
	ToFile string
}

var (
	localStorage *objectstorage.Storage
)

func init() {
	l, err := local.New(context.Background(), "./testbucket/", nil)
	if err != nil {
		log.Fatal(err)
	}
	localStorage = objectstorage.New(l)

	// prepare data for download testing
	data := []byte("haloha")
	if err := ioutil.WriteFile("./testbucket/testdownload.txt", data, 0744); err != nil {
		log.Fatal(err)
	}
}

func TestUpload(t *testing.T) {
	cases := []struct {
		Data        interface{}
		ExpectError error
	}{
		// upload content
		{
			Data: Content{
				Body:   []byte("haloha"),
				ToFile: "haloha.txt",
			},
			ExpectError: nil,
		},
		// upload file
		{
			Data: File{
				FileName: "objectstorage.go",
			},
			ExpectError: nil,
		},
	}

	for _, c := range cases {
		var (
			filepath string
			err      error
		)

		switch c.Data.(type) {
		case File:
			d := c.Data.(File)

			f, err := os.Open(d.FileName)
			if err != c.ExpectError {
				t.Error(err)
				return
			}

			filepath, err = localStorage.Upload(context.TODO(), f, d.FileName, nil)
			if err != nil {
				t.Error(err)
				return
			}
		case Content:
			d := c.Data.(Content)
			buff := bytes.NewBuffer(d.Body)
			filepath, err = localStorage.Upload(context.TODO(), buff, d.ToFile, nil)
			if err != c.ExpectError {
				t.Error(err)
				return
			}
		}

		if filepath != "" {
			_, err = os.Stat(filepath)
			if err != nil {
				t.Error(err)
				return
			}

			os.Remove(filepath)
			os.Remove(fmt.Sprintf("%s.%s", filepath, ".attrs"))
		}
	}
}

func TestUploadFile(t *testing.T) {
	cases := []struct {
		FileName    string
		ExpectError error
	}{
		{
			FileName:    "objectstorage.go",
			ExpectError: nil,
		},
		// file not found
		{
			FileName:    "haloha.go",
			ExpectError: os.ErrNotExist,
		},
	}

	for _, c := range cases {
		var (
			filepath string
			err      error
		)

		filepath, err = localStorage.UploadFile(context.TODO(), c.FileName, c.FileName, nil)
		if err != c.ExpectError {
			if !os.IsNotExist(err) {
				t.Error(err)
				return
			}
		}

		if filepath != "" {
			_, err = os.Stat(filepath)
			if err != nil {
				t.Error(err)
				return
			}

			os.Remove(filepath)
			log.Println(filepath)
			os.Remove(fmt.Sprintf("%s.%s", filepath, ".attrs"))
		}
	}
}

func TestUploadByte(t *testing.T) {
	cases := []struct {
		Content     []byte
		ToFile      string
		ExpectError error
	}{
		{
			Content:     []byte("objectstorage.go"),
			ToFile:      "haloha.txt",
			ExpectError: nil,
		},
		// byte is empty
		{
			Content:     nil,
			ToFile:      "nil.txt",
			ExpectError: objectstorage.ErrByteEmpty,
		},
	}

	for _, c := range cases {
		var (
			filepath string
			err      error
		)

		filepath, err = localStorage.UploadByte(context.TODO(), c.Content, c.ToFile, nil)
		if err != c.ExpectError {
			t.Error(err)
			return
		}

		if filepath != "" {
			_, err = os.Stat(filepath)
			if err != nil {
				t.Error(err)
				return
			}

			os.Remove(filepath)
			log.Println(filepath)
			os.Remove(fmt.Sprintf("%s.%s", filepath, ".attrs"))
		}
	}
}

func TestDownloadByte(t *testing.T) {
	cases := []struct {
		Key           string
		ExpectContent []byte
		ExpectError   error
	}{
		{
			Key:           "testdownload.txt",
			ExpectContent: []byte("haloha"),
			ExpectError:   nil,
		},
	}

	for _, c := range cases {
		b, err := localStorage.DownloadByte(context.TODO(), c.Key, nil)
		if err != c.ExpectError {
			t.Error(err)
			return
		}

		if string(b) != string(c.ExpectContent) {
			t.Error("content doesn't match")
			return
		}
	}
}

func TestDownloadFile(t *testing.T) {
	cases := []struct {
		Key         string
		ExpectError error
	}{
		{
			Key:         "testdownload.txt",
			ExpectError: nil,
		},
	}

	for _, c := range cases {
		err := localStorage.DownloadFile(context.TODO(), c.Key, c.Key, nil)
		if err != c.ExpectError {
			t.Error(err)
			return
		}

		_, err = os.Stat(c.Key)
		if err != nil {
			if !os.IsExist(err) {
				t.Error(err)
				return
			}
		}

		os.Remove(c.Key)
	}
}
