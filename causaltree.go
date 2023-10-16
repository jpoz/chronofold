package chronofold

import (
	"fmt"
	"strings"
)

type Timestamp struct {
	Author    string
	AuthorIdx int
}

func (t Timestamp) String() string {
	return fmt.Sprintf("%s(%d)", t.Author, t.AuthorIdx)
}

type Op struct {
	Timestamp Timestamp
	Ref       Timestamp
	Value     Value
}

type CausalTree struct {
	Log []*Op
}

func NewCT(ops ...*Op) *CausalTree {
	return &CausalTree{
		Log: ops,
	}
}

func (ct *CausalTree) Ndx(t Timestamp) (int, *Op) {
	for i := t.AuthorIdx; i < len(ct.Log); i++ {
		if ct.Log[i].Timestamp == t {
			return i, ct.Log[i]
		}
	}

	return -1, nil
}

func (ct *CausalTree) Add(op *Op) {
	ct.Log = append(ct.Log, op)
}

func (ct *CausalTree) Inspect() string {
	var builder strings.Builder
	builder.WriteString("CausalTree{\n")
	for _, n := range ct.Log {
		builder.WriteString(
			fmt.Sprintf(
				"\tOp{Timestamp: %s, Ref: %s, Value: %s},\n",
				n.Timestamp.String(),
				n.Ref.String(),
				n.Value.String(),
			),
		)
	}
	builder.WriteString("}")
	return builder.String()
}
