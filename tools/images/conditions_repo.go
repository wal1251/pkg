package images

var _ ConditionRepo = (ConditionsRepo)(nil)

type (
	ConditionsRepo map[string]Condition
)

func (c ConditionsRepo) Predicate(name string) Predicate {
	if cond, ok := c[name]; ok {
		return cond.Check
	}

	return nil
}

func DefaultConditions() ConditionsRepo {
	conditions := make(ConditionsRepo)
	conditions["OrientationHorizontal"] = Predicate(func(i Image) bool {
		size := i.Bounds().Size()

		return size.X >= size.Y
	})

	conditions["OrientationVertical"] = Predicate(func(i Image) bool {
		size := i.Bounds().Size()

		return size.Y > size.X
	})

	return conditions
}
