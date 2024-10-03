package dnd

import (
	"fmt"
	"sort"
	"strconv"
)

// index and telegram user ID for currently playing char
func (g *Game) CombatIndexData() (int, int64, error) {
	return 0, 0, nil
}

func (g *Game) CombatIfAllPlayersDead() bool {
	for _, char := range g.ActiveParty {
		if char.Hitpoints > 0 && !char.IsNPC {
			return false
		}
	}
	return true
}

// TODO remake to check all NPC
func (g *Game) CombatIfAllNPCDead() bool {
	return g.CurrentLocation.Host.Hitpoints <= 0
}

func (g *Game) CombatStart() string {
	if g.CombatFlag {
		message := "наши байцы будут выступать в таком порядке:\n"
		for i, c := range g.CombatOrder {
			message += strconv.Itoa(i+1) + " - " + c.Name + " с инициативой " + strconv.Itoa(c.Initiative) + "\n"
		}
		return message
	}
	var order []*Char

	//g.Locations[0].Host.Initiative = g.Locations[0].Host.GetInitiative()
	g.CurrentLocation.Host.Initiative = g.CurrentLocation.Host.GetInitiative()
	order = append(order, g.CurrentLocation.Host)
	for k := range g.ActiveParty {
		char := g.ActiveParty[k]
		char.Initiative = char.GetInitiative()
		order = append(order, char)
		g.ActiveParty[k] = char
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

func (g *Game) CombatRouter() (message string, userID int64) {
	charLink := g.CombatOrder[g.CombatIndex]
	for userID, char := range g.ActiveParty {
		if char.Name == charLink.Name {
			g.CombatToNextChar()
			return g.CombatCBMessage, userID
		}
	}
	g.CombatToNextChar()
	return g.CombatTurn(0), 0
}

func (g *Game) CombatTurn(userID int64) string {
	message := ""
	// case of NPC
	if userID == 0 {
		// TODO redo for possible several NPC in one location
		char := g.CurrentLocation.Host
		targetID := g.CombatPickRandomPlayer()
		targetChar := g.ActiveParty[targetID]
		message += fmt.Sprintf("%s бьёт🤜 по %s\n", char.Name, targetChar.Name)
		dmg, messagedmg := char.GetAttackDamage(targetChar.AC)
		message += messagedmg
		message += fmt.Sprintf("\nхп цели: %d - %d = %d\n", targetChar.Hitpoints, dmg, targetChar.Hitpoints-dmg)
		targetChar.Hitpoints -= dmg
		if targetChar.Hitpoints <= 0 {
			message += "цель perished💀\n"
		}
		// "save" target after attack
		g.ActiveParty[targetID] = targetChar
		return message
	}
	// case of real motherfucker
	return ""
}

func (g *Game) CombatToNextChar() {
	if g.CombatIndex == len(g.CombatOrder)-1 {
		g.CombatIndex = 0
	} else {
		g.CombatIndex++
	}
}

func (g *Game) CombatPickRandomPlayer() (userID int64) {
	for i, v := range g.ActiveParty {
		if v.Hitpoints > 0 {
			return i
		}
	}
	return 0
}

func (g *Game) CombatReadyCheck() (message string, flag bool) {
	if g.CombatFlag {
		return "бой уже идёт, это лишнее", false
	}
	if g.CombatIfAllPlayersDead() {
		return "все игроки мертвы, боя не будет", false
	}
	if g.CombatIfAllNPCDead() {
		return "все нпц в этой локации мертвы, боя не будет", false
	}
	return "", true
}

// sorting by initiatives
type ByInitiative []*Char

func (a ByInitiative) Len() int { return len(a) }
func (a ByInitiative) Less(i, j int) bool {
	// Sort by initiative in ascending order
	return a[i].Initiative > a[j].Initiative
}
func (a ByInitiative) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
