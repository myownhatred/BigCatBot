package dnd

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Gender string
type Race string
type Class string
type Location string
type DamageType string
type WeaponProperty string

const (
	Spermtank         Gender         = "спермобак"
	Vaginacapitallist Gender         = "вагинокапиталист"
	Moongender        Gender         = "мунгендер"
	Agender           Gender         = "агендер"
	Gendervoy         Gender         = "гендервой"
	Gendervoid        Gender         = "гендервойд"
	Nonbinary         Gender         = "нонбайори"
	Xenogender        Gender         = "ксеногендер"
	Dwarf             Race           = "дворф"
	Halfling          Race           = "халфлинг"
	Human             Race           = "хуманс"
	Elf               Race           = "эльф"
	Drow              Race           = "драу"
	Gnome             Race           = "гнум"
	Dragonborn        Race           = "драгонборн"
	Halforc           Race           = "полуорк"
	Halfelf           Race           = "полуэльф"
	Tiefling          Race           = "тифлинг"
	Artificer         Class          = "изобретатель"
	Barbarian         Class          = "барбариан"
	Bard              Class          = "бард"
	Cleric            Class          = "жрец"
	Druid             Class          = "друль"
	Fighter           Class          = "солдат"
	Monk              Class          = "монк"
	Paladin           Class          = "паллАдин"
	Ranger            Class          = "егерь"
	Rogue             Class          = "шельма"
	Sorcerer          Class          = "колдун"
	Warlock           Class          = "военный замок"
	Wizard            Class          = "визард"
	Bar               Location       = "Бар"
	Temple            Location       = "Храм"
	Tavern            Location       = "Таверна"
	DamageDubas       DamageType     = "дубасящий"
	DamagePierce      DamageType     = "колющий"
	DamageSlash       DamageType     = "режущий"
	WPVersatile       WeaponProperty = "универсальное"
	WPFencing         WeaponProperty = "Фехтовальное"
)

type Weapon struct {
	Name        string
	CostGold    int
	CostSilver  int
	DamType     DamageType
	Weight      int
	DamageRolls int
	DamageDice  int
}

type Char struct {
	Name       string
	Gender     Gender
	Race       Race
	Class      Class
	Hitpoints  int
	AC         int
	Str        int
	Dex        int
	Con        int
	Intl       int
	Wis        int
	Cha        int
	Level      int
	Weapon     *Weapon
	Generation string
}

var genders = [...]Gender{Spermtank, Vaginacapitallist, Moongender, Agender, Gendervoy, Gendervoid, Nonbinary, Xenogender}
var races = [...]Race{Dwarf, Halfling, Human, Elf, Drow, Gnome, Dragonborn, Halfelf, Halforc, Tiefling}
var classes = [...]Class{Artificer, Barbarian, Bard, Cleric, Druid, Fighter, Monk, Paladin, Ranger, Rogue, Sorcerer, Warlock, Wizard}

type Gaem struct {
	Party []Char
	Loc   Location
}

func RollChar() Char {
	genders := []Gender{Spermtank, Vaginacapitallist, Moongender, Agender, Gendervoy, Gendervoid, Nonbinary, Xenogender}
	races := []Race{Dwarf, Halfling, Human, Elf, Drow, Gnome, Dragonborn, Halfelf, Halforc, Tiefling}
	classes := []Class{Artificer, Barbarian, Bard, Cleric, Druid, Fighter, Monk, Paladin, Ranger, Rogue, Sorcerer, Warlock, Wizard}
	var chel Char
	chel.Gender = genders[rand.Intn(len(genders))]
	chel.Class = classes[rand.Intn(len(classes))]
	chel.Race = races[rand.Intn(len(races))]
	chel.Generation += "Гендир: " + string(chel.Gender) + "\n"
	chel.Generation += "Расса: " + string(chel.Race) + "\n"
	chel.Generation += "Клас: " + string(chel.Class) + "\n"
	// picking two stats to be +1
	bon1 := -1
	bon2 := -1
	if chel.Race == Halfelf {
		bon1 = rand.Intn(5)
		bon2 = rand.Intn(5)
		for bon2 == bon1 {
			bon2 = rand.Intn(5)
		}
	}
	chel.Str, _ = dice3of4()
	if chel.Race == Dwarf || chel.Race == Dragonborn || chel.Race == Halforc {
		chel.Str += 2
	} else if bon1 == 0 || bon2 == 0 || chel.Race == Human {
		chel.Str += 1
	}
	chel.Generation += "Сила: " + strconv.Itoa(chel.Str) + "\n"
	chel.Dex, _ = dice3of4()
	if chel.Race == Halfling || chel.Race == Elf {
		chel.Dex += 2
	} else if bon1 == 1 || bon2 == 1 || chel.Race == Human {
		chel.Dex += 1
	}
	chel.Generation += "Ловкость: " + strconv.Itoa(chel.Dex) + "\n"
	chel.Con, _ = dice3of4()
	if chel.Race == Dwarf {
		chel.Con += 2
	} else if chel.Race == Halforc || bon1 == 2 || bon2 == 2 || chel.Race == Human {
		chel.Con += 1
	}
	chel.Generation += "Телосложение: " + strconv.Itoa(chel.Con) + "\n"
	chel.Intl, _ = dice3of4()
	if chel.Race == Gnome {
		chel.Intl += 2
	} else if bon1 == 3 || bon2 == 3 || chel.Race == Tiefling || chel.Race == Human {
		chel.Intl += 1
	}
	chel.Generation += "Интеллект: " + strconv.Itoa(chel.Intl) + "\n"
	chel.Wis, _ = dice3of4()
	if bon1 == 4 || bon2 == 4 || chel.Race == Human {
		chel.Wis += 1
	}
	chel.Generation += "Мудрость: " + strconv.Itoa(chel.Wis) + "\n"
	chel.Cha, _ = dice3of4()
	if chel.Race == Halfelf || chel.Race == Tiefling {
		chel.Cha += 2
	} else if chel.Race == Dragonborn || chel.Race == Human || chel.Race == Drow {
		chel.Cha += 1
	}
	chel.Generation += "Харя: " + strconv.Itoa(chel.Cha) + "\n"

	chel.Level = 1

	// hit points
	switch chel.Class {
	case Artificer:
		chel.Hitpoints = 8 + calculateBonus(chel.Con)
	case Barbarian:
		chel.Hitpoints = 12 + calculateBonus(chel.Con)
	case Bard:
		chel.Hitpoints = 8 + calculateBonus(chel.Con)
	case Cleric:
		chel.Hitpoints = 8 + calculateBonus(chel.Con)
	case Druid:
		chel.Hitpoints = 8 + calculateBonus(chel.Con)
	case Fighter:
		chel.Hitpoints = 10 + calculateBonus(chel.Con)
	case Monk:
		chel.Hitpoints = 8 + calculateBonus(chel.Con)
	case Paladin:
		chel.Hitpoints = 10 + calculateBonus(chel.Con)
	case Ranger:
		chel.Hitpoints = 10 + calculateBonus(chel.Con)
	case Rogue:
		chel.Hitpoints = 8 + calculateBonus(chel.Con)
	case Sorcerer:
		chel.Hitpoints = 6 + calculateBonus(chel.Con)
	case Warlock:
		chel.Hitpoints = 8 + calculateBonus(chel.Con)
	case Wizard:
		chel.Hitpoints = 6 + calculateBonus(chel.Con)
	}

	chel.Generation += "Хиты: " + strconv.Itoa(chel.Hitpoints) + "\n"

	var dub Weapon
	dub.Name = "дубинка"
	if chel.Cha > 15 && chel.Con > 15 && chel.Gender != Vaginacapitallist {
		dub.Name = "пенис большой"
	}
	dub.DamType = DamageDubas
	dub.DamageRolls = 1
	dub.DamageDice = 4
	chel.Weapon = &dub
	return chel
}

func CharFromData(str, dex, con, intl, wis, cha, gender, race, class int) (Char, error) {
	var chel Char
	// validity check if what
	valid := ifDicesStat(str) && ifDicesStat(dex) && ifDicesStat(con) &&
		ifDicesStat(intl) && ifDicesStat(wis) && ifDicesStat(cha)
	if !valid {
		return chel, fmt.Errorf("невалидные статы")
	}
	if gender > len(genders)-1 || gender < 0 {
		return chel, fmt.Errorf("неверный гендер")
	}
	if race > len(races)-1 || race < 0 {
		return chel, fmt.Errorf("неверная rase")
	}
	if class > len(classes)-1 || class < 0 {
		return chel, fmt.Errorf("неверный класс")
	}
	chel.Gender = genders[gender]
	chel.Race = races[race]
	chel.Class = classes[class]
	chel.Str = str
	chel.Dex = dex
	chel.Con = con
	chel.Intl = intl
	chel.Wis = wis
	chel.Cha = cha
	chel.Level = 1
	return chel, nil
}

func (char *Char) CharStats() (str, dex, con, intl, wis, cha int) {
	return char.Str, char.Dex, char.Con, char.Intl, char.Wis, char.Cha
}

func dice3of4() (val int, scrib string) {
	rand.Seed(time.Now().UnixNano())
	min := rand.Intn(6) + 1
	summ := min
	scrib += strconv.Itoa(min)
	for i := 0; i < 3; i++ {
		k := rand.Intn(6) + 1
		summ += k
		if k < min {
			min = k
		}
		scrib += " + " + strconv.Itoa(k)
	}
	scrib += " = " + strconv.Itoa(summ) + " => " + strconv.Itoa(summ-min)
	return summ - min, scrib
}

func dice8() (val int) {
	return rand.Intn(8) + 1
}

func dice12() (val int) {
	return rand.Intn(12) + 1
}

func ifDicesStat(stat int) bool {
	if stat < 3 || stat > 18 {
		return false
	}
	return true
}

func calculateBonus(value int) int {
	if value > 10 {
		return (value - 10) / 2
	} else {
		return (value - 11) / 2 // This will give -1 for 9 or 8
	}
}
