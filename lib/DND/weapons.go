package dnd

import "math/rand"

type Weapon struct {
	Name             string
	CostGold         int
	CostSilver       int
	DamType          DamageType
	Range            int
	LongRange        int
	Weight           int
	DamageRolls      int
	DamageDice       int
	VersRolls        int
	VersDice         int
	WeaponProperties []WeaponProperty
	AC               int
}

type WeaponProperty string

const (
	WPVersatile WeaponProperty = "универсальное"
	WPFencing   WeaponProperty = "фехтовальное"
	WPLight     WeaponProperty = "лёгкое"
	WPThrown    WeaponProperty = "метательное"
	WPTwohanded WeaponProperty = "двуручное"
	WPRange     WeaponProperty = "дальнее"
	WPAmmo      WeaponProperty = "боеприпас"
)

func CreateWeaponCommon() *Weapon {
	switch rand.Intn(8) + 1 {
	case 1:
		return CreateWeaponBattleStaff()
	case 2:
		return CreateWeaponMace()
	case 3:
		return CreateWeaponClub()
	case 4:
		return CreateWeaponDagger()
	case 5:
		return CreateWeaponHandaxe()
	case 6:
		return CreateWeaponJavelin()
	case 7:
		return CreateWeaponLightHammer()
	case 8:
		return CreateWeaponGreatclub()
	}
	return nil
}

func CreateWeaponClub() *Weapon {
	var club Weapon
	club.Name = "Дубинка"
	club.CostSilver = 1
	club.DamType = DamageDubas
	club.DamageRolls = 1
	club.DamageDice = 4
	club.Weight = 2
	club.WeaponProperties = []WeaponProperty{WPLight}
	return &club
}

func CreateWeaponMace() *Weapon {
	var mace Weapon
	mace.Name = "Булава"
	mace.CostGold = 5
	mace.DamType = DamageDubas
	mace.DamageRolls = 1
	mace.DamageDice = 6
	mace.Weight = 4
	mace.WeaponProperties = []WeaponProperty{}

	return &mace
}

func CreateWeaponBattleStaff() *Weapon {
	var bs Weapon
	bs.Name = "Боевой Посох"
	bs.CostSilver = 2
	bs.DamType = DamageDubas
	bs.DamageRolls = 1
	bs.DamageDice = 6
	bs.Weight = 4
	bs.WeaponProperties = []WeaponProperty{WPVersatile}

	return &bs
}

func CreateWeaponDagger() *Weapon {
	var dagger Weapon
	dagger.Name = "Кенджал"
	dagger.CostGold = 2
	dagger.DamType = DamagePierce
	dagger.DamageRolls = 1
	dagger.DamageDice = 4
	dagger.Weight = 2
	dagger.Range = 20
	dagger.LongRange = 60
	dagger.WeaponProperties = []WeaponProperty{WPLight, WPFencing, WPThrown}

	return &dagger
}

func CreateWeaponHandaxe() *Weapon {
	var handaxe Weapon
	handaxe.Name = "Топорик"
	handaxe.CostGold = 5
	handaxe.DamType = DamageSlash
	handaxe.DamageRolls = 1
	handaxe.DamageDice = 6
	handaxe.Weight = 2
	handaxe.Range = 20
	handaxe.LongRange = 60
	handaxe.WeaponProperties = []WeaponProperty{WPLight, WPThrown}

	return &handaxe
}

func CreateWeaponJavelin() *Weapon {
	var javelin Weapon
	javelin.Name = "Дротик"
	javelin.CostSilver = 5
	javelin.DamType = DamagePierce
	javelin.DamageRolls = 1
	javelin.DamageDice = 6
	javelin.Weight = 2
	javelin.Range = 30
	javelin.LongRange = 120
	javelin.WeaponProperties = []WeaponProperty{WPThrown}

	return &javelin
}

func CreateWeaponLightHammer() *Weapon {
	var lighthammer Weapon
	lighthammer.Name = "Молоточек"
	lighthammer.CostGold = 2
	lighthammer.DamType = DamageDubas
	lighthammer.DamageRolls = 1
	lighthammer.DamageDice = 4
	lighthammer.Weight = 2
	lighthammer.Range = 20
	lighthammer.LongRange = 60
	lighthammer.WeaponProperties = []WeaponProperty{WPLight, WPThrown}

	return &lighthammer
}

func CreateWeaponGreatclub() *Weapon {
	var greatclub Weapon
	greatclub.Name = "Дубина"
	greatclub.CostSilver = 2
	greatclub.DamType = DamageDubas
	greatclub.DamageRolls = 1
	greatclub.DamageDice = 8
	greatclub.Weight = 10
	greatclub.WeaponProperties = []WeaponProperty{WPTwohanded}

	return &greatclub
}

// Martial Melee Weapons

func CreateBattleaxe() *Weapon {
	var w Weapon
	w.Name = "Боевой топор"
	w.CostGold = 10
	w.DamType = DamageSlash
	w.DamageRolls = 1
	w.DamageDice = 8
	w.VersRolls = 1
	w.VersDice = 10
	w.Weight = 4
	w.WeaponProperties = []WeaponProperty{WPVersatile}

	return &w
}

func CreateShortSword() *Weapon {
	var w Weapon
	w.Name = "Короткий меч"
	w.CostGold = 10
	w.DamType = DamageSlash
	w.DamageRolls = 1
	w.DamageDice = 6
	w.Weight = 2
	w.WeaponProperties = []WeaponProperty{WPLight, WPFencing}

	return &w
}

// ranged martial weapons

func CreateWeaponLongbow() *Weapon {
	var w Weapon
	w.Name = "Длинный лук"
	w.CostGold = 50
	w.DamType = DamagePierce
	w.Range = 150
	w.LongRange = 600
	w.DamageRolls = 1
	w.DamageDice = 8
	w.Weight = 2
	w.WeaponProperties = []WeaponProperty{}

	return &w
}

func CreateWeaponShield() *Weapon {
	var armor Weapon
	armor.Name = "Шыт"
	armor.CostGold = 10
	armor.AC = 2
	armor.Weight = 6

	return &armor
}

// misc

func (w *Weapon) ifFencing() bool {
	for _, p := range w.WeaponProperties {
		if p == WPFencing {
			return true
		}
	}
	return false
}
