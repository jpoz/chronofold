package chronofold_test

import (
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
		&c.Node{Timestamp: c.Timestamp{"a", 0}, Value: c.Root{}, Next: c.Increment{}},
		&c.Node{Timestamp: c.Timestamp{"a", 1}, Value: c.Symbol{'H'}, Next: c.Increment{}},
		&c.Node{Timestamp: c.Timestamp{"a", 2}, Value: c.Symbol{'E'}, Next: c.Increment{}},
		&c.Node{Timestamp: c.Timestamp{"a", 3}, Value: c.Symbol{'L'}, Next: c.Increment{}},
		&c.Node{Timestamp: c.Timestamp{"a", 4}, Value: c.Symbol{'L'}, Next: c.Increment{}},
		&c.Node{Timestamp: c.Timestamp{"a", 5}, Value: c.Symbol{'O'}, Next: c.End{}},
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

func TestChronoFold_Insert(t *testing.T) {
	cf := c.FromString("Hello", "a")

	ts := cf.Timestamp(5)

	cf.Insert(c.Op{
		Target:    ts,
		Timestamp: c.Timestamp{"b", 5},
		Value:     c.Tombstone{},
	})

	assert.Equal(t, "Hell", cf.String())
}

func TestChronoFold_ManyInserts(t *testing.T) {
	cf := c.FromString("Hello", "a")

	start := 5
	for i, rune := range " how are you?" {
		target := cf.Timestamp(len(cf.Log) - 1)
		cf.Insert(
			c.Op{
				Target:    target,
				Timestamp: c.Timestamp{"b", start + i},
				Value:     c.Symbol{rune},
			},
		)
	}

	assert.Equal(t, "Hello how are you?", cf.String())
}

func TestChronoFold_PaperExample(t *testing.T) {
	cf := c.FromString("PINSK", "a") // ɑ6

	cf.Insert(c.Op{cf.Timestamp(1), c.Timestamp{"b", 7}, c.Tombstone{}})
	cf.Insert(c.Op{cf.Last(), c.Timestamp{"b", 8}, c.Symbol{'M'}})

	assert.Equal(t, "MINSK", cf.String()) // ɑ6β8

	cf = c.FromString("PINSK", "a")
	cf.Insert(c.Op{cf.Timestamp(2), c.Timestamp{"b", 7}, c.Tombstone{}})
	cf.Insert(c.Op{cf.Timestamp(3), c.Timestamp{"b", 8}, c.Tombstone{}})
	cf.Insert(c.Op{cf.Timestamp(4), c.Timestamp{"b", 8}, c.Tombstone{}})
	cf.Insert(c.Op{cf.Timestamp(5), c.Timestamp{"b", 10}, c.Tombstone{}})
	cf.Insert(c.Op{cf.Last(), c.Timestamp{"b", 11}, c.Symbol{'i'}})
	cf.Insert(c.Op{cf.Last(), c.Timestamp{"b", 12}, c.Symbol{'n'}})
	cf.Insert(c.Op{cf.Last(), c.Timestamp{"b", 13}, c.Symbol{'s'}})
	cf.Insert(c.Op{cf.Last(), c.Timestamp{"b", 14}, c.Symbol{'k'}})

	assert.Equal(t, "Pinsk", cf.String()) // ɑ6y14

	cf.Insert(c.Op{cf.Timestamp(1), c.Timestamp{"b", 7}, c.Tombstone{}})
	cf.Insert(c.Op{cf.Last(), c.Timestamp{"b", 8}, c.Symbol{'M'}})

	assert.Equal(t, "Minsk", cf.String()) // ɑ6y14β8

	cf = c.FromString("PINSK", "a")
	cf.Insert(c.Op{cf.Timestamp(1), c.Timestamp{"b", 7}, c.Tombstone{}})
	cf.Insert(c.Op{cf.Last(), c.Timestamp{"b", 8}, c.Symbol{'M'}})
	cf.Insert(c.Op{cf.Timestamp(2), c.Timestamp{"b", 7}, c.Tombstone{}})
	cf.Insert(c.Op{cf.Timestamp(3), c.Timestamp{"b", 8}, c.Tombstone{}})
	cf.Insert(c.Op{cf.Timestamp(4), c.Timestamp{"b", 8}, c.Tombstone{}})
	cf.Insert(c.Op{cf.Timestamp(5), c.Timestamp{"b", 10}, c.Tombstone{}})
	cf.Insert(c.Op{cf.Last(), c.Timestamp{"b", 11}, c.Symbol{'i'}})
	cf.Insert(c.Op{cf.Last(), c.Timestamp{"b", 12}, c.Symbol{'n'}})
	cf.Insert(c.Op{cf.Last(), c.Timestamp{"b", 13}, c.Symbol{'s'}})
	cf.Insert(c.Op{cf.Last(), c.Timestamp{"b", 14}, c.Symbol{'k'}})

	assert.Equal(t, "Minsk", cf.String()) // ɑ6β8y14
}
