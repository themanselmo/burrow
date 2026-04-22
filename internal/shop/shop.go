package shop

type FoodItem struct {
	ID            string
	Name          string
	Description   string
	Cost          int
	EnergyRestore float64
}

var FoodItems = []FoodItem{
	{
		ID:            "small_fish",
		Name:          "Small Fish",
		Description:   "A quick snack.",
		Cost:          5,
		EnergyRestore: 30,
	},
	{
		ID:            "herring",
		Name:          "Herring",
		Description:   "A satisfying meal.",
		Cost:          10,
		EnergyRestore: 55,
	},
	{
		ID:            "big_mackerel",
		Name:          "Big Mackerel",
		Description:   "Fully restores energy.",
		Cost:          18,
		EnergyRestore: 100,
	},
}
