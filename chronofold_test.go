package chronofold_test

import (
	"fmt"
	"math/rand"
	"testing"

	c "github.com/jpoz/chronofold"
	"github.com/stretchr/testify/assert"
)

func TestChronoFold_empty(t *testing.T) {
	cf := c.Empty("a")
	assert.NotNil(t, cf)
	assert.Equal(t, cf.String(), "")
}

func TestChronoFold_hello(t *testing.T) {
	cf := c.New(
		&c.Node{Value: c.Root{}, Next: c.Increment{}},
		&c.Node{Value: c.Symbol{'H'}, Next: c.Increment{}},
		&c.Node{Value: c.Symbol{'E'}, Next: c.Increment{}},
		&c.Node{Value: c.Symbol{'L'}, Next: c.Increment{}},
		&c.Node{Value: c.Symbol{'L'}, Next: c.Increment{}},
		&c.Node{Value: c.Symbol{'O'}, Next: c.End{}},
	)
	assert.NotNil(t, cf)
	assert.Equal(t, cf.String(), "HELLO")
}

func TestChronoFold_jello(t *testing.T) {
	cf := c.New(
		&c.Node{Value: c.Root{}, Next: c.Increment{}},
		&c.Node{Value: c.Symbol{'H'}, Next: c.Index{6}},
		&c.Node{Value: c.Symbol{'E'}, Next: c.Increment{}},
		&c.Node{Value: c.Symbol{'L'}, Next: c.Increment{}},
		&c.Node{Value: c.Symbol{'L'}, Next: c.Increment{}},
		&c.Node{Value: c.Symbol{'O'}, Next: c.End{}},
		&c.Node{Value: c.Tombstone{}, Next: c.Increment{}},
		&c.Node{Value: c.Symbol{'J'}, Next: c.Index{2}},
	)
	assert.NotNil(t, cf)
	assert.Equal(t, cf.String(), "JELLO")
}

func TestChronoFold_NewChronologFromString(t *testing.T) {
	cf := c.FromString("Hello how are you?", "a")
	assert.NotNil(t, cf)
	assert.Equal(t, cf.String(), "Hello how are you?")
}

func TestChronoFold_Add(t *testing.T) {
	cf := c.FromString("Hello", "a")
	ct := c.NewCT(
		&c.Op{c.Timestamp{"a", 0}, c.Timestamp{"a", 0}, c.Root{}},
		&c.Op{c.Timestamp{"a", 1}, c.Timestamp{"a", 0}, c.Symbol{'H'}},
		&c.Op{c.Timestamp{"a", 2}, c.Timestamp{"a", 1}, c.Symbol{'e'}},
		&c.Op{c.Timestamp{"a", 3}, c.Timestamp{"a", 2}, c.Symbol{'l'}},
		&c.Op{c.Timestamp{"a", 4}, c.Timestamp{"a", 3}, c.Symbol{'l'}},
		&c.Op{c.Timestamp{"a", 5}, c.Timestamp{"a", 4}, c.Symbol{'o'}},
	)

	op := &c.Op{
		Timestamp: c.Timestamp{"b", 6},
		Ref:       c.Timestamp{"a", 5},
		Value:     c.Tombstone{},
	}

	ct.Add(op)
	err := cf.Add(op, ct)
	assert.NoError(t, err)

	assert.Equal(t, "Hell", cf.String())
}

func TestChronoFold_ManyAdds(t *testing.T) {
	cf := c.FromString("Hello", "a")
	ct := c.NewCT(
		&c.Op{c.Timestamp{"a", 0}, c.Timestamp{"a", 0}, c.Root{}},
		&c.Op{c.Timestamp{"a", 1}, c.Timestamp{"a", 0}, c.Symbol{'H'}},
		&c.Op{c.Timestamp{"a", 2}, c.Timestamp{"a", 1}, c.Symbol{'e'}},
		&c.Op{c.Timestamp{"a", 3}, c.Timestamp{"a", 2}, c.Symbol{'l'}},
		&c.Op{c.Timestamp{"a", 4}, c.Timestamp{"a", 3}, c.Symbol{'l'}},
		&c.Op{c.Timestamp{"a", 5}, c.Timestamp{"a", 4}, c.Symbol{'o'}},
	)

	op := &c.Op{
		Timestamp: c.Timestamp{"b", 6},
		Ref:       c.Timestamp{"a", 5},
		Value:     c.Symbol{' '},
	}
	ct.Add(op)
	err := cf.Add(op, ct)
	assert.NoError(t, err)

	start := 6
	for i, rune := range "how are you?" {
		op := &c.Op{
			Timestamp: c.Timestamp{"b", start + i + 1},
			Ref:       c.Timestamp{"b", start + i},
			Value:     c.Symbol{rune},
		}

		ct.Add(op)
		err := cf.Add(op, ct)
		assert.NoError(t, err)
	}

	assert.Equal(t, "Hello how are you?", cf.String())
}

func TestChronoFold_PaperExample(t *testing.T) {
	base := func() (*c.ChronoFold, *c.CausalTree) {
		cf := c.FromString("PINSK", "a") // ɑ6
		ct := c.NewCT(
			&c.Op{c.Timestamp{"a", 0}, c.Timestamp{"a", 0}, c.Root{}},
			&c.Op{c.Timestamp{"a", 1}, c.Timestamp{"a", 0}, c.Symbol{'P'}},
			&c.Op{c.Timestamp{"a", 2}, c.Timestamp{"a", 1}, c.Symbol{'I'}},
			&c.Op{c.Timestamp{"a", 3}, c.Timestamp{"a", 2}, c.Symbol{'N'}},
			&c.Op{c.Timestamp{"a", 4}, c.Timestamp{"a", 3}, c.Symbol{'S'}},
			&c.Op{c.Timestamp{"a", 5}, c.Timestamp{"a", 4}, c.Symbol{'K'}},
		)
		return cf, ct
	}

	addOp := func(cf *c.ChronoFold, ct *c.CausalTree, op *c.Op) {
		ct.Add(op)
		err := cf.Add(op, ct)
		assert.NoError(t, err)
	}

	cf, ct := base()
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"b", 6}, Ref: c.Timestamp{"a", 1}, Value: c.Tombstone{}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"b", 7}, Ref: c.Timestamp{"b", 6}, Value: c.Symbol{'M'}})
	assert.Equal(t, "MINSK", cf.String()) // ɑ6β8

	cf, ct = base()
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"b", 6}, Ref: c.Timestamp{"a", 5}, Value: c.Tombstone{}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"b", 7}, Ref: c.Timestamp{"a", 4}, Value: c.Tombstone{}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"b", 8}, Ref: c.Timestamp{"a", 3}, Value: c.Tombstone{}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"b", 9}, Ref: c.Timestamp{"a", 2}, Value: c.Tombstone{}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"b", 10}, Ref: c.Timestamp{"a", 1}, Value: c.Symbol{'i'}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"b", 11}, Ref: c.Timestamp{"b", 10}, Value: c.Symbol{'n'}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"b", 12}, Ref: c.Timestamp{"b", 11}, Value: c.Symbol{'s'}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"b", 13}, Ref: c.Timestamp{"b", 12}, Value: c.Symbol{'k'}})
	assert.Equal(t, "Pinsk", cf.String()) // ɑ6y14

	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"c", 6}, Ref: c.Timestamp{"a", 1}, Value: c.Tombstone{}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"c", 7}, Ref: c.Timestamp{"c", 6}, Value: c.Symbol{'M'}})
	assert.Equal(t, "Minsk", cf.String()) // ɑ6y14β8

	cf, ct = base()
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"b", 6}, Ref: c.Timestamp{"a", 1}, Value: c.Tombstone{}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"b", 7}, Ref: c.Timestamp{"b", 6}, Value: c.Symbol{'M'}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"c", 6}, Ref: c.Timestamp{"a", 5}, Value: c.Tombstone{}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"c", 7}, Ref: c.Timestamp{"c", 6}, Value: c.Tombstone{}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"c", 8}, Ref: c.Timestamp{"c", 7}, Value: c.Tombstone{}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"c", 9}, Ref: c.Timestamp{"c", 8}, Value: c.Tombstone{}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"c", 10}, Ref: c.Timestamp{"c", 9}, Value: c.Symbol{'i'}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"c", 11}, Ref: c.Timestamp{"c", 10}, Value: c.Symbol{'n'}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"c", 12}, Ref: c.Timestamp{"c", 11}, Value: c.Symbol{'s'}})
	addOp(cf, ct, &c.Op{Timestamp: c.Timestamp{"c", 13}, Ref: c.Timestamp{"c", 12}, Value: c.Symbol{'k'}})
	assert.Equal(t, "Minsk", cf.String()) // ɑ6β8y14
	fmt.Println(cf.Inspect())
	fmt.Println(ct.Inspect())
}

func BenchmarkGoRuneAppend(b *testing.B) {
	s := []rune("Hello")
	for i := 0; i < b.N; i++ {
		s = append(s, 'o')
	}
}

func BenchmarkChronoFoldAppend(b *testing.B) {
	cf := c.FromString("Hello", "a") // ɑ6
	ct := c.NewCT(
		&c.Op{c.Timestamp{"a", 0}, c.Timestamp{"a", 0}, c.Root{}},
		&c.Op{c.Timestamp{"a", 1}, c.Timestamp{"a", 0}, c.Symbol{'H'}},
		&c.Op{c.Timestamp{"a", 2}, c.Timestamp{"a", 1}, c.Symbol{'e'}},
		&c.Op{c.Timestamp{"a", 3}, c.Timestamp{"a", 2}, c.Symbol{'l'}},
		&c.Op{c.Timestamp{"a", 4}, c.Timestamp{"a", 3}, c.Symbol{'l'}},
		&c.Op{c.Timestamp{"a", 5}, c.Timestamp{"a", 4}, c.Symbol{'0'}},
	)
	start := 5
	for i := 0; i < b.N; i++ {
		op := &c.Op{
			c.Timestamp{"a", start + i + 1},
			c.Timestamp{"a", start + i},
			c.Symbol{'r'},
		}
		ct.Add(op)
		err := cf.Add(op, ct)
		assert.NoError(b, err)
	}
}

func BenchmarkGoRuneInsert(b *testing.B) {
	s := []rune("Hello")
	for i := 0; i < b.N; i++ {
		index := rand.Intn(len(s))
		s = append(s[:index+1], s[index:]...)
	}
}

func BenchmarkChronoFoldInsert(b *testing.B) {
	cf := c.FromString("Hello", "a") // ɑ6
	ct := c.NewCT(
		&c.Op{c.Timestamp{"a", 0}, c.Timestamp{"a", 0}, c.Root{}},
		&c.Op{c.Timestamp{"a", 1}, c.Timestamp{"a", 0}, c.Symbol{'H'}},
		&c.Op{c.Timestamp{"a", 2}, c.Timestamp{"a", 1}, c.Symbol{'e'}},
		&c.Op{c.Timestamp{"a", 3}, c.Timestamp{"a", 2}, c.Symbol{'l'}},
		&c.Op{c.Timestamp{"a", 4}, c.Timestamp{"a", 3}, c.Symbol{'l'}},
		&c.Op{c.Timestamp{"a", 5}, c.Timestamp{"a", 4}, c.Symbol{'0'}},
	)
	start := 5
	for i := 0; i < b.N; i++ {
		index := rand.Intn(len(cf.Log))
		op := &c.Op{
			c.Timestamp{"a", index},
			c.Timestamp{"a", start + i},
			c.Symbol{'r'},
		}
		ct.Add(op)
		cf.Add(op, ct)
	}
}
