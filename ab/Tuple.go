package ab

import (
	"fmt"
	"strings"
)

type Tuple struct {
	VisibleVars map[string]bool
	Name        string
	Data        map[string]string
}

var index = 0
var TupleMap = make(map[string]*Tuple)

func NewTuple(varSet, VisibleVars map[string]bool, tuple *Tuple) *Tuple {
	this := &Tuple{
		VisibleVars: map[string]bool{},
		Data:        make(map[string]string),
	}

	if VisibleVars != nil {
		for key := range VisibleVars {
			this.VisibleVars[key] = true
		}
	}

	if varSet == nil && tuple != nil {
		for key, value := range tuple.Data {
			tuple.Data[key] = value
		}
		for key := range tuple.VisibleVars {
			this.VisibleVars[key] = true
		}
	}

	if varSet != nil {
		for key := range varSet {
			tuple.Data[key] = UnboundVariable
		}
	}

	tuple.Name = fmt.Sprintf("tuple%d", index)
	index++
	TupleMap[tuple.Name] = this

	return this
}
func NewTupleClone(tuple *Tuple) *Tuple {
	return NewTuple(nil, nil, tuple)
}
func NewTupleWith(varSet, VisibleVars map[string]bool) *Tuple {
	return NewTuple(varSet, VisibleVars, nil)
}

func (t *Tuple) Equals(other *Tuple) bool {
	if len(t.VisibleVars) != len(other.VisibleVars) {
		return false
	}

	for key := range t.VisibleVars {
		if _, ok := other.VisibleVars[key]; !ok {
			return false
		}
		if t.Data[key] != other.Data[key] {
			return false
		}
	}

	if ContainsUnboundVariable(t) || ContainsUnboundVariable(other) {
		return false
	}

	return true
}

func ContainsUnboundVariable(t *Tuple) bool {
	for _, value := range t.Data {
		if value == UnboundVariable {
			return true
		}
	}
	return false
}

func (t *Tuple) HashCode() int {
	result := 1
	for key := range t.VisibleVars {
		result = 31*result + HashString(key)
		if value, ok := t.Data[key]; ok {
			result = 31*result + HashString(value)
		}
	}
	return result
}

func HashString(s string) int {
	hash := 0
	for _, c := range s {
		hash = hash*31 + int(c)
	}
	return hash
}

func (t *Tuple) GetVars() map[string]bool {
	return t.VisibleVars
}

func (t *Tuple) PrintVars() string {
	var result strings.Builder
	for key := range t.Data {
		if _, ok := t.VisibleVars[key]; ok {
			result.WriteString(" " + key)
		} else {
			result.WriteString(" [" + key + "]")
		}
	}
	return result.String()
}

func (t *Tuple) GetValue(variable string) string {
	value, ok := t.Data[variable]
	if !ok {
		return DefaultGet
	}
	return value
}

func (t *Tuple) Bind(variable, value string) {
	if currentValue, ok := t.Data[variable]; ok && currentValue != UnboundVariable {
		fmt.Printf("%s already bound to %s\n", variable, currentValue)
	} else {
		t.Data[variable] = value
	}
}

func (t *Tuple) PrintTuple() string {
	var result strings.Builder
	for key, value := range t.Data {
		result.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}
	return result.String()
}
