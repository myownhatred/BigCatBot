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
}

const (
	Bar    = "Бар"
	Temple = "Храм"
	Tavern = "Таверна"
)

type Game struct {
	Party           map[int64]Char
	Locations       []Location
	CurrentLocation Location
}

func NewGame() *Game {
	var barman Char
	barman, _ = CharFromData(14, 11, 14, 10, 12, 8, 0, 2, 2)
	barman.Name = "Васян"
	barman.Title = "бармен"
	barman.Class = ""

	var bar Location
	bar.Name = Bar
	bar.Host = &barman

	var game Game
	game.Party = make(map[int64](Char))

	game.Locations = []Location{bar}

	return &game
}

func (g *Game) SetCurrentLocation() {
	g.CurrentLocation = g.Locations[0]
}

func (g *Game) Lookaround() string {
	result := "вы осматриваете " + g.CurrentLocation.Name + "\n"
	result += g.CurrentLocation.Description + "\n"
	result += "похоже тут главный " + g.CurrentLocation.Host.Title + " " + g.CurrentLocation.Host.Name

	return result
}

func (g *Game) Combat() string {
	var order []Char

	g.Locations[0].Host.Initiative = g.Locations[0].Host.GetInitiative()
	order = append(order, *g.Locations[0].Host)
	for _, c := range g.Party {
		c.Initiative = c.GetInitiative()
		order = append(order, c)
	}
	sort.Sort(ByInitiative(order))
	message := "наши байцы будут выступать в таком порядке:\n"
	for i, c := range order {
		message += strconv.Itoa(i) + " - " + c.Name + " с инициативой " + strconv.Itoa(c.Initiative) + "\n"
	}
	return message
}

// sorting by initiatives
type ByInitiative []Char

func (a ByInitiative) Len() int { return len(a) }
func (a ByInitiative) Less(i, j int) bool {
	// Sort by initiative in ascending order
	return a[i].Initiative < a[j].Initiative
}
func (a ByInitiative) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
