package chronofold

import "fmt"

type Document struct {
	cf  *ChronoFold
	ct  *CausalTree
	max int
}

func NewDocument() *Document {
	cf := Empty()
	ct := NewCT(
		&Op{Timestamp{"root", 0}, Timestamp{"root", 0}, Root{}},
	)

	return &Document{
		cf,
		ct,
		0,
	}
}

func (d *Document) Insert(a string, i int, v rune) error {
	value := Symbol{v}
	t, r, err := d.TimestampPairFor(a, i)
	if err != nil {
		return err
	}

	op := &Op{t, r, value}
	return d.AddOp(op)
}

func (d *Document) Delete(a string, i int) error {
	value := Tombstone{}
	t, r, err := d.TimestampPairFor(a, i)
	if err != nil {
		return err
	}

	op := &Op{t, r, value}
	return d.AddOp(op)
}

func (d *Document) Merge(d2 *Document) error {
	for _, op := range d2.ct.Log {
		if err := d.AddOp(op); err != nil {
			return err
		}
	}

	return nil
}

func (d *Document) AddOp(op *Op) error {
	d.ct.Add(op)
	if err := d.cf.Add(op, d.ct); err != nil {
		return err
	}

	if op.Timestamp.AuthorIdx > d.max {
		d.max = op.Timestamp.AuthorIdx
	}

	return nil
}

func (d *Document) String() string {
	return d.cf.String()
}

func (d *Document) TimestampAt(idx int) Timestamp {
	return d.cf.TimestampAt(idx)
}

func (d *Document) ValueByTimestamp(t Timestamp) Value {
	return d.cf.ValueByTimestamp(t)
}

func (d *Document) Inspect() string {
	return d.cf.Inspect() + "\n" + d.ct.Inspect()
}

func (d *Document) TimestampPairFor(a string, i int) (Timestamp, Timestamp, error) {
	t := Timestamp{a, d.max + 1}
	r := d.TimestampAt(i)
	if r.Author == "" {
		return t, r, fmt.Errorf("invalid index %d", i)
	}

	return t, r, nil
}
