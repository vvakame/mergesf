package mergesf_test

import (
	"encoding/json"
	"testing"

	"github.com/vvakame/mergesf"
)

func Test_motivation(t *testing.T) {
	{ // merge 2 struct into 1 flat JSON.
		type A struct {
			A string
		}
		type B struct {
			B string
		}

		type C struct {
			A
			B
		}

		b, err := json.Marshal(&C{A{A: "a"}, B{B: "b"}})
		if err != nil {
			t.Fatal(err)
		}

		if v := string(b); v != `{"A":"a","B":"b"}` {
			t.Errorf("unexpected: %v", v)
		}
	}
	{ // merge 2 struct, but one is dynamically determined... Can we?
		type A struct {
			A string
		}
		type unknown struct {
			B string
		}
		type B2 interface{}

		type C struct {
			A
			B2
		}

		b, err := json.Marshal(&C{A{A: "a"}, unknown{B: "b"}})
		if err != nil {
			t.Fatal(err)
		}

		// below is the unexpected result! not {"A":"a","B":"b"}
		if v := string(b); v != `{"A":"a","B2":{"B":"b"}}` {
			t.Errorf("unexpected: %v", v)
		}
	}
	{ // use mergesf :+1:
		type A struct {
			A string
		}
		type unknown struct {
			B string
		}

		a := &A{A: "a"}
		var i interface{} = &unknown{B: "b"}

		merged, err := mergesf.Merge(a, i)
		if err != nil {
			t.Fatal(err)
		}

		b, err := json.Marshal(merged)
		if err != nil {
			t.Fatal(err)
		}

		// nice! :+1:
		if v := string(b); v != `{"A":"a","B":"b"}` {
			t.Errorf("unexpected: %v", v)
		}
	}
}
