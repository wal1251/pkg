package images

var _ Condition = (Predicate)(nil)

type (
	Predicate func(Image) bool
)

func MakeCondition(predicates ...Predicate) Condition {
	return Predicate(func(i Image) bool {
		for _, predicate := range predicates {
			if !predicate(i) {
				return false
			}
		}

		return true
	})
}

func (p Predicate) MakeTransformer(t Transformer) Transformer {
	return Transform(func(i Image) (Image, error) {
		if p.Check(i) {
			return t.Perform(i)
		}

		return i, nil
	})
}

func (p Predicate) Check(i Image) bool {
	if p == nil {
		return false
	}

	return p(i)
}
