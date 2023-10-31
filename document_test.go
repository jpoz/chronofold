package chronofold

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	doc := NewDocument()
	doc.Insert("a", 0, 'H')
	doc.Insert("a", 1, 'e')
	doc.Insert("a", 2, 'l')
	doc.Insert("a", 3, 'l')
	doc.Insert("a", 4, 'o')
	doc.Insert("a", 5, '!')

	assert.Equal(t, doc.String(), "Hello!")
}

func TestDelete(t *testing.T) {
	doc := NewDocument()
	doc.Insert("a", 0, 'H')
	doc.Insert("a", 1, 'e')
	doc.Insert("a", 2, 'l')
	doc.Insert("a", 3, 'l')
	doc.Insert("a", 4, 'o')
	doc.Insert("a", 5, '!')

	doc.Delete("a", 5)
	doc.Delete("a", 6)

	assert.Equal(t, doc.String(), "Hell!")
}

func TestMerge(t *testing.T) {
	doc := NewDocument()
	doc.Insert("a", 0, 'H')
	doc.Insert("a", 1, 'e')
	doc.Insert("a", 2, 'l')
	doc.Insert("a", 3, 'l')
	doc.Insert("a", 4, 'o')
	doc.Insert("a", 5, '!')

	doc2 := NewDocument()
	doc2.Insert("b", 0, 'H')
	doc2.Insert("b", 1, 'e')
	doc2.Insert("b", 2, 'l')
	doc2.Insert("b", 3, 'l')
	doc2.Insert("b", 4, 'o')
	doc2.Insert("b", 5, '!')

	doc.Merge(doc2)

	assert.Equal(t, doc.String(), "Hello!Hello!")
}
