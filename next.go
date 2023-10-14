package chronofold

import "fmt"

// Next interface
type Next interface {
	String() string
}

// Increment
type Increment struct{}

func (Increment) String() string {
	return "+1"
}

type End struct{}

func (End) String() string {
	return "End"
}

type Index struct {
	Idx int
}

func (i Index) String() string {
	return fmt.Sprintf("Index(%d)", i.Idx)
}
