package dnd

type Spell struct {
	Name              string
	NameRu            string
	ID                int
	DamageRolls       int
	DamageDice        int
	ComponentVerbal   bool
	ComponentSomatic  bool
	ComponentMaterial bool
	Level             int
	DamageType        DamageType
	SavingStat        Stat
}

func CreateSpellAcidSplash() *Spell {
	var spell Spell
	spell.Name = "Acid Splash"
	spell.NameRu = "Кислый сплеш"
	spell.ID = 1
	spell.DamageRolls = 1
	spell.DamageDice = 6
	spell.ComponentSomatic = true
	spell.ComponentVerbal = true
	spell.Level = 0
	spell.DamageType = DamageAcid
	spell.SavingStat = Dexterity

	return &spell
}

func CreateSpellSacredFlame() *Spell {
	var spell Spell
	spell.Name = "Sacred Flame"
	spell.NameRu = "Вечный огонь"
	spell.ID = 2
	spell.DamageRolls = 1
	spell.DamageDice = 8
	spell.ComponentSomatic = true
	spell.ComponentVerbal = true
	spell.Level = 0
	spell.DamageType = DamageRadian
	spell.SavingStat = Dexterity

	return &spell
}
