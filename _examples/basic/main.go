package main

import (
	"fmt"

	"github.com/jpoz/chronofold"
)

func main() {
	cf := chronofold.Empty()
	ct := chronofold.NewCT(
		&chronofold.Op{chronofold.Timestamp{"a", 0}, chronofold.Timestamp{"a", 0}, chronofold.Root{}},
	)

	for i, r := range "Hello" {
		op := &chronofold.Op{chronofold.Timestamp{"a", i + 1}, chronofold.Timestamp{"a", i}, chronofold.Symbol{r}}
		ct.Add(op)
		cf.Add(op, ct)
	}

	fmt.Println(cf.String())

	fmt.Println(cf.Inspect())
	fmt.Println(ct.Inspect())
}
