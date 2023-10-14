# Chronofold

A Go implementation of the ChronoFold algorithm from https://arxiv.org/abs/2002.09511

## Example

```go
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
	cf.Insert(chronofold.Op{cf.Last(), chronofold.Timestamp{"b", 7}, chronofold.Symbol{'H'}})

	fmt.Println(cf.String()) // Prints "Hello"
}
```

The underlying data structure looks like:

```
ChronoFold{
        Node{Timestamp:{a, 0}, Value: ∅, Next: Index(1)},
        Node{Timestamp:{b, 1}, Value: h, Next: Index(6)},
        Node{Timestamp:{b, 2}, Value: e, Next: Index(3)},
        Node{Timestamp:{b, 3}, Value: l, Next: Index(4)},
        Node{Timestamp:{b, 4}, Value: l, Next: Index(5)},
        Node{Timestamp:{b, 5}, Value: o, Next: End},
        Node{Timestamp:{b, 6}, Value: ⌫, Next: Index(7)},
        Node{Timestamp:{b, 7}, Value: H, Next: Index(2)},
}
```
