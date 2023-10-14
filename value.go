package chronofold

type Value interface {
	String() string
}

type Root struct{}

func (v Root) String() string {
	return "∅"
}

type Tombstone struct{}

func (v Tombstone) String() string {
	return "⌫"
}

type Symbol struct {
	Char rune
}

func (v Symbol) String() string {
	return string(v.Char)
}
