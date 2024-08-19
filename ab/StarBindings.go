package ab

type StarBindings struct {
	InputStars *Stars
	ThatStars  *Stars
	TopicStars *Stars
}

func NewStarBindings() *StarBindings {
	return &StarBindings{
		InputStars: NewStars(),
		ThatStars:  NewStars(),
		TopicStars: NewStars(),
	}
}
