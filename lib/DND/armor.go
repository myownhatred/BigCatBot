package dnd

import "math/rand"

type Armor struct {
	Name           string
	CostGold       int
	CostSilver     int
	AC             int
	NeededStrength int
	BadStealth     bool
	Weight         int
}

func CreateFakeArmor() *Armor {
	switch rand.Intn(10) + 1 {
	case 1:
		return CreateArmorBadRobe()
	case 2, 3, 4, 5, 6, 7, 8, 9:
		return CreateArmorRobe()
	case 10:
		return CreateArmorNiceRobe()
	}
	return CreateArmorRobe()
}

func CreateLightArmor() *Armor {
	switch rand.Intn(3) + 1 {
	case 1:
		return CrateArmorPadded()
	case 2:
		return CrateArmorLeather()
	case 3:
		return CrateArmorStuddedLeather()
	}
	return nil
}

func CreateMediumArmor() *Armor {
	switch rand.Intn(2) + 1 {
	case 1:
		return CrateArmorHide()
	case 2:
		return CrateArmorChainShirt()
	}
	return nil
}

func CreateArmorRobe() *Armor {
	var armor Armor
	armor.Name = "Халат"
	armor.CostSilver = 5
	armor.AC = 10
	armor.Weight = 1

	return &armor
}

func CreateArmorNiceRobe() *Armor {
	var armor Armor
	armor.Name = "Красивый халат"
	armor.CostSilver = 10
	armor.AC = 10
	armor.Weight = 1

	return &armor
}

func CreateArmorBadRobe() *Armor {
	var armor Armor
	armor.Name = "Обоссаный халат"
	armor.CostSilver = 1
	armor.AC = 10
	armor.Weight = 1

	return &armor
}

func CrateArmorPadded() *Armor {
	var quil Armor
	quil.Name = "Стёганый доспех"
	quil.CostGold = 5
	quil.AC = 11
	quil.BadStealth = true
	quil.Weight = 8

	return &quil
}

func CrateArmorLeather() *Armor {
	var leather Armor
	leather.Name = "Кожанный доспех"
	leather.CostGold = 10
	leather.AC = 11
	leather.BadStealth = false
	leather.Weight = 10

	return &leather
}

func CrateArmorStuddedLeather() *Armor {
	var studdedLeather Armor
	studdedLeather.Name = "Клёпанный кожанный доспех"
	studdedLeather.CostGold = 45
	studdedLeather.AC = 12
	studdedLeather.BadStealth = false
	studdedLeather.Weight = 13

	return &studdedLeather
}

func CrateArmorHide() *Armor {
	var armor Armor
	armor.Name = "Шкура"
	armor.CostGold = 10
	armor.AC = 12
	armor.BadStealth = false
	armor.Weight = 12

	return &armor
}

func CrateArmorChainShirt() *Armor {
	var armor Armor
	armor.Name = "Кольчужная рубаха"
	armor.CostGold = 50
	armor.AC = 13
	armor.BadStealth = false
	armor.Weight = 20

	return &armor
}

func CrateArmorScaleMail() *Armor {
	var armor Armor
	armor.Name = "Чешуйчатый доспех"
	armor.CostGold = 50
	armor.AC = 14
	armor.BadStealth = true
	armor.Weight = 45

	return &armor
}

func CreateArmorBreastplate() *Armor {
	var armor Armor
	armor.Name = "Кираса"
	armor.CostGold = 400
	armor.AC = 14
	armor.BadStealth = false
	armor.Weight = 20

	return &armor
}

func CreateArmorHalfplate() *Armor {
	var armor Armor
	armor.Name = "Полулаты"
	armor.CostGold = 750
	armor.AC = 15
	armor.BadStealth = true
	armor.Weight = 40

	return &armor
}

func CreateArmorShield() *Armor {
	var armor Armor
	armor.Name = "Шыт"
	armor.CostGold = 10
	armor.AC = 2
	armor.BadStealth = false
	armor.Weight = 6

	return &armor
}
