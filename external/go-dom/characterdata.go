package dom

type CharacterData interface {
	Node

	GetValue() string
	SetValue(string)
}
