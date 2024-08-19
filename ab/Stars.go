package ab

type Stars struct {
	Items []string
}

func NewStars() *Stars {
	return &Stars{
		Items: make([]string, 0),
	}
}

func (s *Stars) Add(item string) {
	s.Items = append(s.Items, item)
}

func (s *Stars) Star(i int) string {
	if i < len(s.Items) {
		return s.Items[i]
	}
	return ""
}
