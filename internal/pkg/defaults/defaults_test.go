package defaults

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type (
	Def struct {
		Str string        `default:"abcde"`
		I64 int64         `default:"123"`
		I32 int32         `default:"123"`
		I16 int16         `default:"123"`
		I8  int8          `default:"123"`
		I   int           `default:"123"`
		F32 float32       `default:"123.32"`
		F64 float64       `default:"123.32"`
		Dur time.Duration `default:"1m23s"`
	}

	Def2 struct {
		I16 int16
		F32 float32
		Str string
		I64 int32
	}
)

var (
	allSet = Def{
		Str: "abcde",
		I64: 123,
		I32: 123,
		I16: 123,
		I8:  123,
		I:   123,
		F32: 123.32,
		F64: 123.32,
		Dur: 1*time.Minute + 23*time.Second,
	}
	allSet2   = allSet
	notStruct = "not struct"
)

func TestSetDefault(t *testing.T) {
	allSet2.Str = "abc"
	notStruct := "not struct"

	testCases := []struct {
		name string
		inp  interface{}
		out  interface{}
		err  error
	}{
		{
			name: "from all empty",
			inp:  &Def{},
			out:  &allSet,
			err:  nil,
		},
		{
			name: "all set",
			inp:  &allSet,
			out:  &allSet,
			err:  nil,
		},
		{
			name: "pass value",
			inp:  Def{},
			out:  Def{},
			err:  ErrPassValue,
		},
		{
			name: "pass nil",
			inp:  nil,
			out:  nil,
			err:  ErrPassValue,
		},
		{
			name: "partially set",
			inp:  &Def{Str: "abc"},
			out:  &allSet2,
			err:  nil,
		},
		{
			name: "not struct",
			inp:  &notStruct,
			out:  &notStruct,
			err:  ErrNotStruct,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SetDefault(tc.inp)
			require.Equal(t, tc.err, err)

			require.Equal(t, tc.out, tc.inp)
		})
	}
}

func TestSetDefaultFrom(t *testing.T) {
	allSet2.Str = "abc"

	testCases := []struct {
		name string
		src  interface{}
		rep  interface{}
		exp  interface{}
		err  error
	}{
		{
			name: "from all empty",
			src:  &Def{},
			rep:  allSet,
			exp:  &allSet,
			err:  nil,
		},
		{
			name: "pass value",
			src:  Def{},
			rep:  allSet,
			exp:  Def{},
			err:  ErrPassValue,
		},
		{
			name: "pass nil",
			src:  nil,
			rep:  allSet,
			exp:  nil,
			err:  ErrPassValue,
		},
		{
			name: "not struct",
			src:  &notStruct,
			rep:  notStruct,
			exp:  &notStruct,
			err:  ErrNotStruct,
		},
		{
			name: "partially set",
			src:  &Def{Str: "abc"},
			rep:  allSet,
			exp:  &allSet2,
			err:  nil,
		},
		{
			name: "different type",
			src:  &Def{},
			rep: Def2{
				Str: "abcde",
				I16: 123,
				F32: 123.32,
				I64: 123,
			},
			exp: &Def{
				Str: "abcde",
				I16: 123,
				F32: 123.32,
				I64: 0, // not set due to mismatch type
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ReplaceDefaultFrom(tc.src, tc.rep)
			require.Equal(t, tc.err, err)

			require.Equal(t, tc.exp, tc.src)
		})
	}
}
