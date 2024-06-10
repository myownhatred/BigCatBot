package dnd

import (
	"math/rand"
	"strconv"
	"time"
)

type Gender string
type Race string
type Class string

const (
	Spermtank         Gender = "спермобак"
	Vaginacapitallist Gender = "вагинокапиталист"
	Moongender        Gender = "мунгендер"
	Agender           Gender = "агендер"
	Gendervoy         Gender = "гендервой"
	Gendervoid        Gender = "гендервойд"
	Nonbinary         Gender = "нонбайори"
	Xenogender        Gender = "ксеногендер"
	Dwarf             Race   = "дворф"
	Halfling          Race   = "халфлинг"
	Human             Race   = "хуманс"
	Elf               Race   = "эльф"
	Drow              Race   = "драу"
	Gnome             Race   = "гнум"
	Dragonborn        Race   = "драгонборн"
	Halforc           Race   = "полуорк"
	Halfelf           Race   = "полуэльф"
	Tiefling          Race   = "тифлинг"
	Artificer         Class  = "артифисер"
	Barbarian         Class  = "барбариан"
	Bard              Class  = "бесполезный"
	Cleric            Class  = "жрец"
	Druid             Class  = "друль"
	Fighter           Class  = "солдат"
	Monk              Class  = "монк"
	Paladin           Class  = "паллАдин"
	Ranger            Class  = "егерь"
	Rogue             Class  = "шельма"
	Sorcerer          Class  = "колдун"
	Warlock           Class  = "военный замок"
	Wizard            Class  = "визард"
)

type Char struct {
	Name       string
	Gender     Gender
	Race       Race
	Class      Class
	Str        int
	Dex        int
	Con        int
	Intl       int
	Wis        int
	Cha        int
	Level      int
	Generation string
}

func RollChar() Char {
	genders := []Gender{Spermtank, Vaginacapitallist, Moongender, Agender, Gendervoy, Gendervoid, Nonbinary, Xenogender}
	races := []Race{Dwarf, Halfling, Human, Elf, Drow, Gnome, Dragonborn, Halfelf, Halforc, Tiefling}
	classes := []Class{Artificer, Barbarian, Bard, Cleric, Druid, Fighter, Monk, Paladin, Ranger, Rogue, Sorcerer, Warlock, Wizard}
	rand.Seed(time.Now().UnixNano())
	var chel Char
	chel.Gender = genders[rand.Intn(len(genders))]
	chel.Class = classes[rand.Intn(len(classes))]
	chel.Race = races[rand.Intn(len(races))]
	chel.Str, tmp := dice3of4()
	chel.Generation += "Сила: " + tmp + "\n"
	chel.Dex, tmp = dice3of4()
	chel.Generation += "Ловкость: " + tmp + "\n"
	chel.Con, tmp = dice3of4()
	chel.Generation += "Телосложение: " + tmp + "\n"
	chel.Intl, tmp = dice3of4()
	chel.Generation += "Интеллект: " + tmp + "\n"
	chel.Wis, tmp = dice3of4()
	chel.Generation += "Мудрость: " + tmp + "\n"
	chel.Cha, tmp = dice3of4()
	chel.Generation += "Харя: " + tmp + "\n"
	chel.Level = 1

	return chel
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
