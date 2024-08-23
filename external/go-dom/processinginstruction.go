package dom

type ProcessingInstruction interface {
	CharacterData

	GetTarget() string
	SetTarget(string)
}
