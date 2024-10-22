package dnd

func NewCleric(chel *Char) {
	chel.Armor = CrateArmorChainShirt()
	chel.Weapon = CreateWeaponMace()
	chel.WeaponOffhand = CreateWeaponShield()
	chel.Spells = append(chel.Spells, *CreateSpellSacredFlame())
}
