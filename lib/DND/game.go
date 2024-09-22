package dnd

type Location struct {
	Name        string
	Description string
	Host        *Char
	Locals      []Char
	Actions     []Action
}

type Action struct {
	Name  string
	Die   int
	Stat  Stat
	Level int
}

const (
	Bar    = "Бар"
	Temple = "Храм"
	Plaza  = "Площадь"
)

type Game struct {
	Party           map[int64]Char
	ActiveParty     map[int64]*Char
	Locations       []Location
	CurrentLocation *Location
	CombatOrder     []*Char
	CombatIndex     int
	CombatFlag      bool
	CombatFC        chan bool
	CombatCBMessage string
}

func NewGame() *Game {
	var barman Char
	barman, _ = CharFromData(10, 10, 10, 10, 10, 10, 14, 25, 0, 2, 2)
	barman.Name = "Васян"
	barman.Title = "бармен"
	barman.Class = ""
	barman.Weapon = CreateWeaponClub()
	barman.IsNPC = true

	var bomzh Char
	bomzh, _ = CharFromData(10, 10, 10, 10, 10, 10, 13, 1, 0, 2, 2)
	bomzh.Name = "Керилл"
	bomzh.Title = "бомж"
	bomzh.Class = ""
	bomzh.Weapon = CreateWeaponClub()
	bomzh.IsNPC = true

	var plaza Location

	plaza.Name = Plaza
	plaza.Host = &bomzh
	plaza.Description = "площадь деревни скрытого листа, в луже по центру лежит мощный бомж Керилл"

	var bar Location
	bar.Name = Bar
	bar.Host = &barman
	bar.Description = "бар с одним видом пива - нефильтрованная пшеничка, на заккуску только чеснок, за стойкой очень мощный бармен Васян"

	var game Game
	game.Party = make(map[int64](Char))
	game.ActiveParty = make(map[int64](*Char))

	game.Locations = []Location{plaza, bar}
	game.CurrentLocation = &plaza
	game.CombatFC = make(chan bool)

	return &game
}

func (g *Game) SetCurrentLocation(i int) {
	g.CurrentLocation = &g.Locations[i]
}

func (g *Game) Lookaround() string {
	result := "вы осматриваете " + g.CurrentLocation.Name + "\n"
	result += g.CurrentLocation.Description + "\n"
	result += "похоже тут главный " + g.CurrentLocation.Host.Title + " " + g.CurrentLocation.Host.Name

	return result
}
