package ulid

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

// UlidIface interface, best to use instead of Ulid struct
type UlidIface interface {
	Ulid() string
}

// Ulid struct
type Ulid struct {
	ulidchan chan string
}

// Worker of ulid
type Worker struct {
	ulidchan chan string
}

func (w *Worker) Work(source rand.Source) {
	r := rand.New(source)
	for {
		w.ulidchan <- ulid.MustNew(ulid.Now(), r).String()
	}
}

// New ulid
func New(workernumber int) *Ulid {
	ch := make(chan string, workernumber*10)
	u := Ulid{
		ulidchan: ch,
	}

	for i := 0; i < workernumber; i++ {
		s := rand.NewSource(time.Now().UnixNano())
		w := Worker{
			ulidchan: ch,
		}
		go w.Work(s)
	}
	return &u
}

// Ulid return ulid string
func (u *Ulid) Ulid() string {
	return <-u.ulidchan
}

// UlidMock for mocking ulid
type UlidMock struct {
	ulidID     []string
	defaultVal string
}

// NewMock return new mock object
func NewMock(ulid ...string) *UlidMock {
	return &UlidMock{
		ulidID:     ulid,
		defaultVal: "loremipsumdolorsitamet",
	}
}

// DefaultValue return default value of mock
func (u *UlidMock) DefaultValue() string {
	return u.defaultVal
}

// Ulid return mock ulid value
func (u *UlidMock) Ulid() string {
	if u == nil {
		return ""
	}

	if len(u.ulidID) == 0 {
		return u.defaultVal
	}

	temp := u.ulidID[0]
	u.ulidID = u.ulidID[1:]

	return temp
}
