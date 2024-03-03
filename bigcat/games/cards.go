package bigcat

import (
	"fmt"

	"github.com/chen3feng/stl4go"
)

// gaming cards structs
const (
	Ace = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Valet
	Dama
	Korol
	// mast
	Hearts   = 0
	Diamonds = 1
	Spades   = 2
	Clubs    = 3
)

var Masty []string = []string{"Черви", "Бубны", "Пики", "Трефы"}
var Values []string = []string{"Туз", "Двойка", "Тройка", "Четвёрка", "Пятёрка", "Шестёрка", "Семёрка", "Восьмёрка", "Девятка", "Десятка", "Валет", "Дама", "Король"}
var AHearts []string = []string{"🂱", "🂲", "🂳", "🂴", "🂵", "🂶", "🂷", "🂸", "🂹", "🂺", "🂻", "🂽", "🂾"}
var ADiamonds []string = []string{"🃁", "🃂", "🃃", "🃄", "🃅", "🃆", "🃇", "🃈", "🃉", "🃊", "🃋", "🃍", "🃎"}
var ASpades []string = []string{"🂡", "🂢", "🂣", "🂤", "🂥", "🂦", "🂧", "🂨", "🂩", "🂪", "🂫", "🂭", "🂮"}
var AClubs []string = []string{"🃑", "🃒", "🃓", "🃔", "🃕", "🃖", "🃗", "🃘", "🃙", "🃚", "🃛", "🃝", "🃞"}
var Full [][]string = [][]string{AHearts, ADiamonds, ASpades, AClubs}

type Card struct {
	Value int
	Mast  int
}

type Deck struct {
	DeckVec stl4go.Vector[Card]
}

func (d *Deck) NewStack() {
	// making full card stack
	for i := 0; i < 4; i++ {
		for j := 0; j < 13; j++ {
			d.DeckVec.Append(Card{Value: j, Mast: i})
		}
	}
}

func (d *Deck) TopDeck() Card {
	card := d.DeckVec.At(0)
	d.DeckVec.Remove(1)
	return card
}

func (c *Card) OneShot() string {
	return Full[c.Mast][c.Value]
}

func (c *Card) String() string {
	return fmt.Sprintf("%s %s", Values[c.Value], Masty[c.Mast])
}
