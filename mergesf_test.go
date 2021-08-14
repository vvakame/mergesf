package mergesf_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/vvakame/mergesf"
)

func TestMerge(t *testing.T) {
	mergesf.RecoverPanic = true

	tests := []struct {
		name    string
		args    []interface{}
		want    string
		wantErr bool
	}{
		{
			name: "basic usage",
			args: []interface{}{
				&A{Name: "yukari", Age: 7},
				&B{Kind: "Norwegian Forest Cat & Ragdoll"},
			},
			want: `{"Name":"yukari","Age":7,"Kind":"Norwegian Forest Cat \u0026 Ragdoll"}`,
		},
		{
			name: "with json tag",
			args: []interface{}{
				&A{Name: "yukari", Age: 7},
				&B{Kind: "Norwegian Forest Cat & Ragdoll"},
				&C{Favorite: "Foxtail"},
			},
			want: `{"Name":"yukari","Age":7,"Kind":"Norwegian Forest Cat \u0026 Ragdoll","fav":"Foxtail"}`,
		},
		{
			name: "len 0",
			args: nil,
			want: `null`,
		},
		{
			name: "len 1",
			args: []interface{}{
				&A{Name: "yukari", Age: 7},
			},
			want: `{"Name":"yukari","Age":7}`,
		},
		{
			name: "error - duplicated field",
			args: []interface{}{
				&A{Name: "yukari", Age: 7},
				&DupA{Name: "vvakame"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergesf.Merge(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Merge() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				return
			}

			b, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}

			if got := string(b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type A struct {
	Name string
	Age  int
}

func (a *A) Hello() {
	fmt.Printf("%s, %d", a.Name, a.Age)
}

type B struct {
	Kind string
}

func (b *B) Hi() {
	fmt.Printf("%s", b.Kind)
}

type C struct {
	Favorite string `json:"fav"`
}

type DupA struct {
	Name string
}
