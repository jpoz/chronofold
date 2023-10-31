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

	for i, r := range "The ability to write well remains a coveted skill. *Weather* it's crafting a compelling narrative or composing a technical document." {
		op := &chronofold.Op{chronofold.Timestamp{"a", i + 1}, chronofold.Timestamp{"a", i}, chronofold.Symbol{r}}
		ct.Add(op)
		cf.Add(op, ct)
	}

	fmt.Println(cf.String())

	fmt.Println(cf.Inspect())
	fmt.Println(ct.Inspect())
}
