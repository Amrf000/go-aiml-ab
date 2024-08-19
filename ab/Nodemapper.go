package ab

type Nodemapper struct {
	Category     *Category
	Height       int
	StarBindings *StarBindings
	Map          map[string]*Nodemapper
	Key          string
	Value        *Nodemapper
	ShortCut     bool
	Sets         []string
}

func NewNodemapper() *Nodemapper {
	this := &Nodemapper{}
	this.Height = MaxGraphHeight
	this.Map = map[string]*Nodemapper{}
	this.ShortCut = false
	return this
}
