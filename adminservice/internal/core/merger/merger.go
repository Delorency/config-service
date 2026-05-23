package merger

type Merger struct {
	priorities map[string]int
}

func NewMerger() *Merger {
	return &Merger{
		priorities: map[string]int{
			"global": 1,
			"app":    2,
			"env":    3,
		},
	}
}
