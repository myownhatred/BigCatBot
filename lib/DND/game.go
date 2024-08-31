package dnd

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
