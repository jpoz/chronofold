package main

import (
	"fmt"

	"github.com/jpoz/chronofold"
)

func main() {
	cf := chronofold.FromString("", "a")

	cf.Insert(chronofold.Op{cf.Last(), chronofold.Timestamp{"b", 1}, chronofold.Symbol{'h'}})
	cf.Insert(chronofold.Op{cf.Last(), chronofold.Timestamp{"b", 2}, chronofold.Symbol{'e'}})
	cf.Insert(chronofold.Op{cf.Last(), chronofold.Timestamp{"b", 3}, chronofold.Symbol{'l'}})
	cf.Insert(chronofold.Op{cf.Last(), chronofold.Timestamp{"b", 4}, chronofold.Symbol{'l'}})
	cf.Insert(chronofold.Op{cf.Last(), chronofold.Timestamp{"b", 5}, chronofold.Symbol{'o'}})
	cf.Insert(chronofold.Op{cf.Timestamp(1), chronofold.Timestamp{"b", 6}, chronofold.Tombstone{}}) // Delete h
	cf.Insert(chronofold.Op{cf.Last(), chronofold.Timestamp{"b", 1}, chronofold.Symbol{'H'}})

	fmt.Println(cf.String()) // Prints Hello

	fmt.Println(cf.Inspect())
}
