package mergesf_test

import (
	"encoding/json"
	"fmt"

	"github.com/vvakame/mergesf"
)

type Cat struct {
	Name string
	Age  int
}

type AnimalAttr struct {
	Kind string
}

func Example() {
	merged, err := mergesf.Merge(
		&Cat{Name: "yukari", Age: 7},
		&AnimalAttr{Kind: "Norwegian Forest Cat & Ragdoll"},
	)
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(merged)
	if err != nil {
		panic(err)
	}

	// Output:
	// {"Name":"yukari","Age":7,"Kind":"Norwegian Forest Cat \u0026 Ragdoll"}
	fmt.Println(string(b))
}

func (cat *Cat) Hello() {
	fmt.Printf("%s, %d", cat.Name, cat.Age)
}

func (attr *AnimalAttr) Hi() {
	fmt.Printf("%s", attr.Kind)
}
