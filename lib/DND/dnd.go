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

const (
	Spermtank         Gender   = "спермобак"
	Vaginacapitallist Gender   = "вагинокапиталист"
	Moongender        Gender   = "мунгендер"
	Agender           Gender   = "агендер"
	Gendervoy         Gender   = "гендервой"
	Gendervoid        Gender   = "гендервойд"
	Nonbinary         Gender   = "нонбайори"
	Xenogender        Gender   = "ксеногендер"
	Dwarf             Race     = "дворф"
	Halfling          Race     = "халфлинг"
	Human             Race     = "хуманс"
	Elf               Race     = "эльф"
	Drow              Race     = "драу"
	Gnome             Race     = "гнум"
	Dragonborn        Race     = "драгонборн"
	Halforc           Race     = "полуорк"
	Halfelf           Race     = "полуэльф"
	Tiefling          Race     = "тифлинг"
	Artificer         Class    = "изобретатель"
	Barbarian         Class    = "барбариан"
	Bard              Class    = "бард"
	Cleric            Class    = "жрец"
	Druid             Class    = "друль"
	Fighter           Class    = "солдат"
	Monk              Class    = "монк"
	Paladin           Class    = "паллАдин"
	Ranger            Class    = "егерь"
	Rogue             Class    = "шельма"
	Sorcerer          Class    = "колдун"
	Warlock           Class    = "военный замок"
	Wizard            Class    = "визард"
	Bar               Location = "Бар"
	Temple            Location = "Храм"
	Tavern            Location = "Таверна"
)

type Char struct {
	Name       string
	Gender     Gender
	Race       Race
	Class      Class
	AC         int
	Str        int
	Dex        int
	Con        int
	Intl       int
	Wis        int
	Cha        int
	Level      int
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
	rand.Seed(time.Now().UnixNano())
	var chel Char
	tmp := ""
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
	chel.Str, tmp = dice3of4()
	if chel.Race == Dwarf || chel.Race == Dragonborn || chel.Race == Halforc {
		chel.Generation += "Сила: " + tmp + " +2\n"
		chel.Str += 2
	} else if bon1 == 0 || bon2 == 0 || chel.Race == Human {
		chel.Generation += "Сила: " + tmp + " +1\n"
		chel.Str += 1
	} else {
		chel.Generation += "Сила: " + tmp + "\n"
	}
	chel.Dex, tmp = dice3of4()
	if chel.Race == Halfling || chel.Race == Elf {
		chel.Generation += "Ловкость: " + tmp + " +2\n"
		chel.Dex += 2
	} else if bon1 == 1 || bon2 == 1 || chel.Race == Human {
		chel.Generation += "Ловкость: " + tmp + " +1\n"
		chel.Dex += 1
	} else {
		chel.Generation += "Ловкость: " + tmp + "\n"
	}
	chel.Con, tmp = dice3of4()
	if chel.Race == Dwarf {
		chel.Generation += "Телосложение: " + tmp + " +2\n"
		chel.Con += 2
	} else if chel.Race == Halforc || bon1 == 2 || bon2 == 2 || chel.Race == Human {
		chel.Generation += "Телосложение: " + tmp + " +1\n"
		chel.Con += 1
	} else {
		chel.Generation += "Телосложение: " + tmp + "\n"
	}
	chel.Intl, tmp = dice3of4()
	if chel.Race == Gnome {
		chel.Generation += "Интеллект: " + tmp + " +2\n"
		chel.Intl += 2
	} else if bon1 == 3 || bon2 == 3 || chel.Race == Tiefling || chel.Race == Human {
		chel.Generation += "Интеллект: " + tmp + " +1\n"
		chel.Intl += 1
	} else {
		chel.Generation += "Интеллект: " + tmp + "\n"
	}
	chel.Wis, tmp = dice3of4()
	if bon1 == 4 || bon2 == 4 || chel.Race == Human {
		chel.Generation += "Мудрость: " + tmp + " +1\n"
		chel.Wis += 1
	} else {
		chel.Generation += "Мудрость: " + tmp + "\n"
	}
	chel.Cha, tmp = dice3of4()
	if chel.Race == Halfelf || chel.Race == Tiefling {
		chel.Generation += "Харя: " + tmp + " +2\n"
		chel.Cha += 2
	} else if chel.Race == Dragonborn || chel.Race == Human {
		chel.Generation += "Харя: " + tmp + " +1\n"
		chel.Cha += 1
	} else {
		chel.Generation += "Харя: " + tmp + "\n"
	}
	chel.Level = 1

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

func ifDicesStat(stat int) bool {
	if stat < 3 || stat > 18 {
		return false
	}
	return true
}
