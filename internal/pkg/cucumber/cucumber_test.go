// this test is intended for testing the godog itself and working as integration test

package cucumber_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/albertwidi/go-project-example/internal/pkg/cucumber"
	"github.com/cucumber/godog"
)

func runWebServer() error {
	type BookRequest struct {
		BookID int64 `json:"book_id"`
	}
	type BookResponse struct {
		BookID int64  `json:"book_id"`
		Name   string `json:"name"`
		Author string `json:"author"`
	}
	http.HandleFunc("/v1/book/detail", func(w http.ResponseWriter, r *http.Request) {
		b := BookRequest{}
		out, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if err := json.Unmarshal(out, &b); err != nil {
			err = fmt.Errorf("BookDetailEndpoint: %w", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		resp := BookResponse{
			BookID: b.BookID,
			Name:   "testing",
			Author: "what?",
		}
		out, err = json.Marshal(resp)
		if err != nil {
			err = fmt.Errorf("BookDetailEndpoint: %w", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(out)
		return
	})
	if err := http.ListenAndServe(":9863", nil); err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	var (
		options = godog.Options{
			Output: os.Stdout,
			Format: "pretty",
		}
		errChan = make(chan error)
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go func() {
		errChan <- runWebServer()
	}()

	select {
	case err := <-errChan:
		log.Fatal(err)
	case <-ctx.Done():
		// this means the webserver is alive
		break
	}

	flag.Parse()
	options.Paths = flag.Args()

	c, err := cucumber.New(&cucumber.Options{
		Debug: cucumber.Debug{
			LogFile: "test.log",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	apiFeature := &cucumber.APIFeature{
		Options: cucumber.APIFeatureOptions{
			EndpointsMapping: map[string]string{
				"book": "http://127.0.0.1:9863",
			},
		},
	}
	c.RegisterFeatures(apiFeature)

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		c.FeatureContext(s)
	}, options)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
