package dnd

type Equipment struct {
	Name       string
	CostGold   int
	CostSilver int
	Stack      int
	Quantity   int
	Weight     int
}

func CreateArrowsStack() *Equipment {
	var e Equipment
	e.Name = "Стак стрел"
	e.CostGold = 1
	e.Stack = 20
	e.Quantity = 20
	e.Weight = 1

	return &e
}
