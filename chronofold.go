package chronofold

import (
	"fmt"
	"strings"
)

type Node struct {
	Value     Value
	Next      Next
	Timestamp Timestamp
}

type ChronoFold struct {
	Log         []*Node
	currentTime int
}

func Empty() *ChronoFold {
	return &ChronoFold{
		Log: []*Node{
			{
				Value: Root{},
				Next:  End{},
				Timestamp: Timestamp{
					Author: "root",
					AuthorIdx:    0,
				},
			},
		},
	}
}

func New(nodes ...*Node) *ChronoFold {
	return &ChronoFold{
		Log: nodes,
	}
}

func FromString(text, author string) *ChronoFold {
	strlen := len(text)
	nodes := make([]*Node, strlen+1)
	nodes[0] = &Node{
		Value: Root{},
		Next:  incrementOrEnd(0, strlen),
	}

	for i, rune := range text {
		nodes[i+1] = &Node{
			Value: Symbol{rune},
			Next:  incrementOrEnd(i+1, strlen),
		}
	}

	return &ChronoFold{Log: nodes}
}

func (c *ChronoFold) Add(op *Op, ct *CausalTree) error {
	ref := op.Ref
	j, prevOp := ct.Ndx(ref)
	if prevOp == nil {
		return fmt.Errorf("invalid timestamp")
	}
	if ref.AuthorIdx > op.Timestamp.AuthorIdx {
		return fmt.Errorf("invalid ref (ref.AuthorIdx > op.Timestamp.AuthorIdx)")
	}
	if ref.AuthorIdx > j {
		return fmt.Errorf("invalid ref (ref.AuthorIdx > Reference index)")
	}
	if j >= len(c.Log) {
		return fmt.Errorf("invalid ref (j >= len(c.Log))")
	}

	prevNode := c.Log[j]
	addedIndex := len(c.Log)
	next := calculateNext(j, prevNode.Next, addedIndex)

	newNode := Node{
		Value:     op.Value,
		Next:      next,
		Timestamp: op.Timestamp,
	}
	c.Log = append(c.Log, &newNode)

	prevNode.Next = normalizeNext(j, addedIndex)

	return nil
}

func (c *ChronoFold) String() string {
	result := []rune{}

	c.iterateByNext(
		func(_ Timestamp) {},
		func(r rune, _ int, _ Timestamp) { result = append(result, r) },
		func(_ Timestamp) { result = result[:len(result)-1] },
	)

	return string(result)
}

func (c *ChronoFold) TimestampAt(tidx int) Timestamp {
	idx := 0
	var result Timestamp

	c.iterateByNext(
		func(t Timestamp) {
			if idx == tidx {
				result = t
			}
			idx++
		},
		func(r rune, _ int, t Timestamp) {
			if idx == tidx {
				result = t
			}
			idx++
		},
		func(_ Timestamp) {
			idx--
		},
	)

	return result
}

func (c *ChronoFold) ValueByTimestamp(tt Timestamp) Value {
	var result Value

	c.iterateByNext(
		func(t Timestamp) {
			if t == tt {
				result = Root{}
			}
		},
		func(r rune, _ int, t Timestamp) {
			if t == tt {
				result = Symbol{r}
			}
		},
		func(t Timestamp) {
			result = Tombstone{}
		},
	)

	return result
}

func (c *ChronoFold) Inspect() string {
	var builder strings.Builder
	builder.WriteString("ChronoFold{\n")
	for _, n := range c.Log {
		builder.WriteString(
			fmt.Sprintf(
				"\tNode{Value: %s, Next: %s, Timestamp: %s},\n",
				n.Value.String(),
				n.Next.String(),
				n.Timestamp.String(),
			),
		)
	}
	builder.WriteString("}")
	return builder.String()
}

func (c *ChronoFold) iterateByNext(
	onRoot func(Timestamp),
	onSymbol func(rune, int, Timestamp),
	onTombstone func(Timestamp),
) error {
	logIdx := int(0)

Loop:
	for {
		n := c.Log[logIdx]
		switch v := n.Value.(type) {
		case Root:
			onRoot(n.Timestamp)
		case Symbol:
			onSymbol(v.Char, logIdx, n.Timestamp)
		case Tombstone:
			onTombstone(n.Timestamp)
		}
		next := n.Next
		switch next.(type) {
		case Increment:
			logIdx++
		case End:
			break Loop
		case Index:
			logIdx = next.(Index).Idx
		}
	}

	return nil
}

func calculateNext(prevIdx int, prevNext Next, sourceIndex int) Next {
	var targetIdx int
	switch prevNext.(type) {
	case Index:
		targetIdx = prevNext.(Index).Idx
	case Increment:
		targetIdx = prevIdx + 1
	case End:
		return prevNext
	}

	return normalizeNext(sourceIndex, targetIdx)
}

func normalizeNext(sourceIndex, targetIndex int) Next {
	if (targetIndex - sourceIndex) == 1 {
		return Increment{}
	}

	return Index{Idx: targetIndex}
}

func incrementOrEnd(curr, size int) Next {
	if curr == size {
		return End{}
	}

	return Increment{}
}
