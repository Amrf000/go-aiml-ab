package dom

type DocumentType interface {
	Node

	GetName() string
	GetPublicID() string
	GetSystemID() string
	GetDefinition() string
}
