package chronofold

import (
	"fmt"
	"strings"
)

type Node struct {
	Value Value
	Next  Next
}

type ChronoFold struct {
	Log         []*Node
	currentTime int
}

func Empty(author string) *ChronoFold {
	return &ChronoFold{
		Log: []*Node{
			{
				Value: Root{},
				Next:  End{},
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
		Value: op.Value,
		Next:  next,
	}
	c.Log = append(c.Log, &newNode)

	prevNode.Next = normalizeNext(j, addedIndex)

	return nil
}

func (c *ChronoFold) String() string {
	result := []rune{}

	c.iterateByNext(
		func(r rune, idx int) { result = append(result, r) },
		func() { result = result[:len(result)-1] },
	)

	return string(result)
}

func (c *ChronoFold) Inspect() string {
	var builder strings.Builder
	builder.WriteString("ChronoFold{\n")
	for _, n := range c.Log {
		builder.WriteString(
			fmt.Sprintf(
				"\tNode{Value: %s, Next: %s},\n",
				n.Value.String(),
				n.Next.String(),
			),
		)
	}
	builder.WriteString("}")
	return builder.String()
}

func (c *ChronoFold) iterateByNext(
	onSymbol func(rune, int),
	onTombstone func(),
) error {
	logIdx := int(0)

Loop:
	for {
		n := c.Log[logIdx]
		switch v := n.Value.(type) {
		case Root:
			// Do nothing
		case Symbol:
			onSymbol(v.Char, logIdx)
		case Tombstone:
			onTombstone()
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
