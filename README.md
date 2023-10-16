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
	cf := chronofold.Empty()

	ct := chronofold.NewCT(
		&chronofold.Op{chronofold.Timestamp{"a", 0}, chronofold.Timestamp{"a", 0}, chronofold.Root{}},
	)

	for i, r := range "Hello" {
		op := &chronofold.Op{chronofold.Timestamp{"a", i + 1}, chronofold.Timestamp{"a", i}, chronofold.Symbol{r}}
		ct.Add(op)
		cf.Add(op, ct)
	}

	fmt.Println(cf.String()) // Prints "Hello"
}
```

The underlying data structures looks like:

```
ChronoFold{
        Node{Value: ∅, Next: +1},
        Node{Value: H, Next: +1},
        Node{Value: e, Next: +1},
        Node{Value: l, Next: +1},
        Node{Value: l, Next: +1},
        Node{Value: o, Next: End},
}
CausalTree{
        Op{Timestamp: a(0), Ref: a(0), Value: ∅},
        Op{Timestamp: a(1), Ref: a(0), Value: H},
        Op{Timestamp: a(2), Ref: a(1), Value: e},
        Op{Timestamp: a(3), Ref: a(2), Value: l},
        Op{Timestamp: a(4), Ref: a(3), Value: l},
        Op{Timestamp: a(5), Ref: a(4), Value: o},
}
```
