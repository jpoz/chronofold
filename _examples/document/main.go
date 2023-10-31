package main

import (
	"fmt"

	"github.com/jpoz/chronofold"
)

func main() {
	var err error
	doc := chronofold.NewDocument()

	err = doc.AddOp(&chronofold.Op{
		chronofold.Timestamp{"a", 0},
		chronofold.Timestamp{"root", 0},
		chronofold.Symbol{'>'},
	})
	if err != nil {
		panic(err)
	}

	err = doc.AddOp(&chronofold.Op{
		chronofold.Timestamp{"b", 0},
		chronofold.Timestamp{"root", 0},
		chronofold.Symbol{'>'},
	})
	if err != nil {
		panic(err)
	}

	for i, r := range "Write" {
		op := &chronofold.Op{chronofold.Timestamp{"b", i + 1}, chronofold.Timestamp{"b", i}, chronofold.Symbol{r}}
		err := doc.AddOp(op)
		if err != nil {
			panic(err)
		}
	}
	for i, r := range "Read" {
		op := &chronofold.Op{chronofold.Timestamp{"a", i + 1}, chronofold.Timestamp{"a", i}, chronofold.Symbol{r}}
		err := doc.AddOp(op)
		if err != nil {
			panic(err)
		}
	}


	fmt.Println(doc.Inspect())
	fmt.Printf("%#v", doc.String())
}
