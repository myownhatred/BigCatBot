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
	WeaponProperties []WeaponProperty
}

type WeaponProperty string

const (
	WPVersatile WeaponProperty = "универсальное"
	WPFencing   WeaponProperty = "фехтовальное"
	WPLight     WeaponProperty = "лёгкое"
	WPThrown    WeaponProperty = "метательное"
)

func CreateWeaponCommon() *Weapon {
	switch rand.Intn(4) + 1 {
	case 1:
		return CreateWeaponBattleStaff()
	case 2:
		return CreateWeaponMace()
	case 3:
		return CreateWeaponClub()
	case 4:
		return CreateWeaponDagger()
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
	dagger.Name = "кенджал"
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
