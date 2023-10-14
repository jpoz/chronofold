package chronofold

import (
	"fmt"
	"strings"
)

type Timestamp struct {
	Author    string
	AuthorIdx int
}

type Node struct {
	Timestamp Timestamp
	Value     Value
	Next      Next
}

type Op struct {
	Target    Timestamp
	Timestamp Timestamp
	Value     Value
}

type ChronoFold struct {
	Log         []*Node
	currentTime int
}

func Empty(author string) *ChronoFold {
	return &ChronoFold{
		Log: []*Node{
			{
				Timestamp: Timestamp{
					Author:    author,
					AuthorIdx: 0,
				},
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
		Timestamp: Timestamp{author, 0},
		Value:     Root{},
		Next:      incrementOrEnd(0, strlen),
	}

	for i, rune := range text {
		nodes[i+1] = &Node{
			Timestamp: Timestamp{author, int(i + 1)},
			Value:     Symbol{rune},
			Next:      incrementOrEnd(i+1, strlen),
		}
	}

	return &ChronoFold{Log: nodes}
}

func (c *ChronoFold) Timestamp(idx int) Timestamp {
	return c.Log[idx].Timestamp
}

func (c *ChronoFold) Last() Timestamp {
	return c.Log[len(c.Log)-1].Timestamp
}

func (c *ChronoFold) Insert(op Op) error {
	j, prevNode := c.ndx(op.Target)
	if prevNode == nil {
		return fmt.Errorf("invalid timestamp")
	}

	addedIndex := len(c.Log)
	next := calculateNext(j, prevNode.Next, addedIndex)

	newNode := Node{
		Timestamp: op.Timestamp,
		Value:     op.Value,
		Next:      next,
	}
	c.Log = append(c.Log, &newNode)

	prevNode.Next = Index{Idx: addedIndex}
	c.Log[j] = prevNode

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
				"\tNode{Timestamp:{%s, %d}, Value: %s, Next: %s},\n",
				n.Timestamp.Author,
				n.Timestamp.AuthorIdx,
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

func (c *ChronoFold) ndx(t Timestamp) (int, *Node) {
	for i, n := range c.Log {
		if n.Timestamp == t {
			return int(i), n
		}
	}

	return 0, nil // idk what todo here
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
