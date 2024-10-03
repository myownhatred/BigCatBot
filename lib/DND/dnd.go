package dnd

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Gender string
type Race string
type Class string
type DamageType string
type Stat string
type ButtonMode string
type ActionType string

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
	Artificer         Class  = "изобретатель"
	Barbarian         Class  = "барбариан"
	Bard              Class  = "бард"
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

	DamageDubas  DamageType = "дубасящий"
	DamagePierce DamageType = "колющий"
	DamageSlash  DamageType = "режущий"
	DamageAcid   DamageType = "кислотный"

	Strenght     Stat = "сила"
	Dexterity    Stat = "ловкость"
	Constitution Stat = "телосложение"
	Intellegence Stat = "интеллект"
	Wisidom      Stat = "мудрость"
	Charisma     Stat = "харизма"

	BtnsActions     ButtonMode = "действие"
	BtnsAttackMelee ButtonMode = "атака рукопашная"
	BtnsAttackRange ButtonMode = "атака с лонгренджи"
	BtnsSpellcast   ButtonMode = "заклы"

	MAttack   ActionType = "мили атака"
	RAttack   ActionType = "ренж атака"
	SpellCast ActionType = "колдунство"
	ItemUse   ActionType = "вещение"
)

type Char struct {
	Name        string
	Title       string
	Gender      Gender
	Race        Race
	Class       Class
	Hitpoints   int
	AC          int
	Str         int
	Dex         int
	Con         int
	Intl        int
	Wis         int
	Cha         int
	Level       int
	Initiative  int
	Weapon      *Weapon
	Armor       *Armor
	Target      *Char
	IsNPC       bool
	Generation  string
	Spells      []Spell
	CastingStat int
	ButtonMode  ButtonMode
}

var genders = [...]Gender{Spermtank, Vaginacapitallist, Moongender, Agender, Gendervoy, Gendervoid, Nonbinary, Xenogender}
var races = [...]Race{Dwarf, Halfling, Human, Elf, Drow, Gnome, Dragonborn, Halfelf, Halforc, Tiefling}
var classes = [...]Class{Artificer, Barbarian, Bard, Cleric, Druid, Fighter, Monk, Paladin, Ranger, Rogue, Sorcerer, Warlock, Wizard}

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
	chel.Generation += "🦬Сила: " + strconv.Itoa(chel.Str) + "\n"
	chel.Dex, _ = dice3of4()
	if chel.Race == Halfling || chel.Race == Elf {
		chel.Dex += 2
	} else if bon1 == 1 || bon2 == 1 || chel.Race == Human {
		chel.Dex += 1
	}
	chel.Generation += "🐈Ловкость: " + strconv.Itoa(chel.Dex) + "\n"
	chel.Con, _ = dice3of4()
	if chel.Race == Dwarf {
		chel.Con += 2
	} else if chel.Race == Halforc || bon1 == 2 || bon2 == 2 || chel.Race == Human {
		chel.Con += 1
	}
	chel.Generation += "🐻Телосложение: " + strconv.Itoa(chel.Con) + "\n"
	chel.Intl, _ = dice3of4()
	if chel.Race == Gnome {
		chel.Intl += 2
	} else if bon1 == 3 || bon2 == 3 || chel.Race == Tiefling || chel.Race == Human {
		chel.Intl += 1
	}
	chel.Generation += "🦊Интеллект: " + strconv.Itoa(chel.Intl) + "\n"
	chel.Wis, _ = dice3of4()
	if bon1 == 4 || bon2 == 4 || chel.Race == Human {
		chel.Wis += 1
	}
	chel.Generation += "🦉Мудрость: " + strconv.Itoa(chel.Wis) + "\n"
	chel.Cha, _ = dice3of4()
	if chel.Race == Halfelf || chel.Race == Tiefling {
		chel.Cha += 2
	} else if chel.Race == Dragonborn || chel.Race == Human || chel.Race == Drow {
		chel.Cha += 1
	}
	chel.Generation += "🦅Харя: " + strconv.Itoa(chel.Cha) + "\n"

	chel.Level = 1

	weapon := CreateWeaponCommon()
	chel.Weapon = weapon
	armor := CreateLightArmor()
	chel.Armor = armor

	// hit points
	switch chel.Class {
	case Artificer:
		chel.Hitpoints = 8 + calculateBonus(chel.Con)
		chel.Spells = append(chel.Spells, *CreateSpellAcidSplash())
		chel.CastingStat = chel.Intl
	case Barbarian:
		chel.Hitpoints = 12 + calculateBonus(chel.Con)
	case Bard:
		chel.Hitpoints = 8 + calculateBonus(chel.Con)
		chel.CastingStat = chel.Cha
	case Cleric:
		chel.Hitpoints = 8 + calculateBonus(chel.Con)
		chel.Spells = append(chel.Spells, *CreateSpellAcidSplash())
		chel.CastingStat = chel.Wis
	case Druid:
		chel.Hitpoints = 8 + calculateBonus(chel.Con)
		chel.Spells = append(chel.Spells, *CreateSpellAcidSplash())
		chel.CastingStat = chel.Wis
	case Fighter:
		chel.Hitpoints = 10 + calculateBonus(chel.Con)
	case Monk:
		chel.Hitpoints = 8 + calculateBonus(chel.Con)
	case Paladin:
		chel.Hitpoints = 10 + calculateBonus(chel.Con)
		chel.CastingStat = chel.Cha
	case Ranger:
		chel.Hitpoints = 10 + calculateBonus(chel.Con)
		switch rand.Intn(2) + 1 {
		case 1:
			chel.Armor = CrateArmorScaleMail()
		case 2:
			chel.Armor = CrateArmorLeather()
		}
		chel.CastingStat = chel.Wis
	case Rogue:
		chel.Hitpoints = 8 + calculateBonus(chel.Con)
	case Sorcerer:
		chel.Hitpoints = 6 + calculateBonus(chel.Con)
		chel.Spells = append(chel.Spells, *CreateSpellAcidSplash())
		chel.CastingStat = chel.Cha
	case Warlock:
		chel.Hitpoints = 8 + calculateBonus(chel.Con)
		chel.Spells = append(chel.Spells, *CreateSpellAcidSplash())
		chel.CastingStat = chel.Cha
	case Wizard:
		chel.Hitpoints = 6 + calculateBonus(chel.Con)
		chel.Spells = append(chel.Spells, *CreateSpellAcidSplash())
		chel.Armor = CreateFakeArmor()
		chel.CastingStat = chel.Intl
	}

	chel.Generation += "❤️Хиты: " + strconv.Itoa(chel.Hitpoints) + "\n"
	// remake AC calc to use dex limitations for med armor
	chel.AC = chel.Armor.AC + calculateBonus(chel.Dex)
	chel.Generation += "🛡️Армор: " + strconv.Itoa(chel.AC) + "\n"
	chel.ButtonMode = BtnsActions
	return chel
}

func CharFromData(str, dex, con, intl, wis, cha, AC, hitpoints, gender, race, class int) (Char, error) {
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
	chel.AC = AC
	chel.Hitpoints = hitpoints
	chel.Level = 1
	return chel, nil
}

func (char *Char) CharStats() (str, dex, con, intl, wis, cha int) {
	return char.Str, char.Dex, char.Con, char.Intl, char.Wis, char.Cha
}

func dice3of4() (val int, scrib string) {
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

func dice6() (val int) {
	return rand.Intn(6) + 1
}

func dice8() (val int) {
	return rand.Intn(8) + 1
}

func dice12() (val int) {
	return rand.Intn(12) + 1
}

func dice20() (val int) {
	return rand.Intn(20) + 1
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

func (c *Char) GetInitiative() int {
	return calculateBonus(c.Dex) + dice20()
}

func (c *Char) GetAttackDamage(target int) (int, string) {
	//hit or miss
	roll := dice20()
	switch roll {
	case 1:
		return 0, "ролл 1 - кретинический промох"
	case 20:
		dmg := 0
		for d := 0; d < 2; d++ {
			for i := 0; i < c.Weapon.DamageRolls; i++ {
				dmg += rand.Intn(c.Weapon.DamageDice) + 1 + c.DnDCharGetWeaponBonus()
			}
		}
		return dmg, fmt.Sprintf("ролл 20 - кретический крит попал на %d уроны", dmg)
	default:
		dmg := 0
		mod := calculateBonus(c.Str)
		if c.Weapon.ifFencing() {
			mod = calculateBonus(c.Dex)
		}
		if roll+mod >= target {
			for i := 0; i < c.Weapon.DamageRolls; i++ {
				dmg += rand.Intn(c.Weapon.DamageDice) + 1 + c.DnDCharGetWeaponBonus()
			}
			return dmg, fmt.Sprintf("ролл %d + мод %d попал на %d уроны", roll, mod, dmg)
		} else {
			return 0, fmt.Sprintf("ролл %d + мод %d промах ибучий", roll, mod)
		}
	}
}

func (c *Char) GetSpellDamage(target *Char, spellIndex int) (int, string) {
	// hit or miss
	// TODO take it from class table
	masteryBonus := 2
	savingBonus := 0
	if c.Spells[spellIndex].SavingStat == target.SavingBonus() {
		savingBonus += 2
	}
	SL := masteryBonus + calculateBonus(c.CastingStat) + 8
	ST := dice20() + savingBonus + target.CalculateBonusByStat(c.Spells[spellIndex].SavingStat)
	if ST > SL {
		return dice6(), "спелл ебашит"
	}
	return 0, "спелл мима"
}

func (c *Char) DnDCharGetWeaponBonus() int {
	for _, p := range c.Weapon.WeaponProperties {
		if p == WPFencing {
			return calculateBonus(c.Dex)
		}
	}
	return calculateBonus(c.Str)
}

func (c *Char) DnDCharIfRangedAttack() bool {
	for _, p := range c.Weapon.WeaponProperties {
		if p == WPThrown || p == WPRange {
			return true
		}
	}
	return false
}

func (c *Char) CalculateBonusByStat(stat Stat) int {
	return calculateBonus(c.StatAtoi(stat))
}

func (c *Char) StatAtoi(stat Stat) int {
	switch stat {
	case Strenght:
		return c.Str
	case Dexterity:
		return c.Dex
	case Constitution:
		return c.Con
	case Intellegence:
		return c.Intl
	case Wisidom:
		return c.Wis
	case Charisma:
		return c.Cha
	}
	return 0
}

// TODO return array, or make checker if given stat present in saving bonus array
func (c *Char) SavingBonus() Stat {
	// TODO take it from class
	return Dexterity
}

// simple placeholder for saving throws
func (c *Char) SavingThrow(stat Stat) int {
	return dice20() + c.CalculateBonusByStat(stat)
}
