package dnd

import (
	"sort"
	"strconv"
)

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
	Locations       []Location
	CurrentLocation *Location
	CombatOrder     []Char
	CombatFlag      bool
}

func NewGame() *Game {
	var barman Char
	barman, _ = CharFromData(10, 10, 10, 10, 10, 10, 10, 8, 0, 2, 2)
	barman.Name = "Васян"
	barman.Title = "бармен"
	barman.Class = ""
	barman.Weapon = CreateWeaponClub()
	barman.IsNPC = true

	var bomzh Char
	bomzh, _ = CharFromData(10, 10, 10, 10, 10, 10, 10, 4, 0, 2, 2)
	bomzh.Name = "Керилл"
	bomzh.Title = "бомж"
	bomzh.Class = ""
	bomzh.Weapon = CreateWeaponClub()
	bomzh.IsNPC = true

	var plaza Location

	plaza.Name = Plaza
	plaza.Host = &bomzh
	plaza.Description = "площадь деревни скрытого листа, в луже по центру лежит бомж Керилл"

	var bar Location
	bar.Name = Bar
	bar.Host = &barman
	bar.Description = "бар с одним видом пива - нефильтрованная пшеничка, на заккуску только чеснок"

	var game Game
	game.Party = make(map[int64](Char))

	game.Locations = []Location{plaza, bar}
	game.CurrentLocation = &plaza

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

func (g *Game) Combat() string {
	if g.CombatFlag {
		message := "наши байцы будут выступать в таком порядке:\n"
		for i, c := range g.CombatOrder {
			message += strconv.Itoa(i+1) + " - " + c.Name + " с инициативой " + strconv.Itoa(c.Initiative) + "\n"
		}

		return message
	}
	var order []Char

	//g.Locations[0].Host.Initiative = g.Locations[0].Host.GetInitiative()
	g.CurrentLocation.Host.Initiative = g.CurrentLocation.Host.GetInitiative()
	order = append(order, *g.CurrentLocation.Host)
	for _, c := range g.Party {
		c.Initiative = c.GetInitiative()
		order = append(order, c)
	}
	sort.Sort(ByInitiative(order))
	message := "наши байцы будут выступать в таком порядке:\n"
	for i, c := range order {
		message += strconv.Itoa(i+1) + " - " + c.Name + " с инициативой " + strconv.Itoa(c.Initiative) + "\n"
	}
	g.CombatOrder = order
	g.CombatFlag = true
	return message
}

// sorting by initiatives
type ByInitiative []Char

func (a ByInitiative) Len() int { return len(a) }
func (a ByInitiative) Less(i, j int) bool {
	// Sort by initiative in ascending order
	return a[i].Initiative > a[j].Initiative
}
func (a ByInitiative) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
