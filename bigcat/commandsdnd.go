package bigcat

import (
	"Guenhwyvar/servitor"
	"fmt"
	"log/slog"

	tele "gopkg.in/telebot.v3"
)

func DnDJoinActive(c tele.Context, serv *servitor.Servitor, brain *BigBrain) (err error) {
	username := "АНОНИМ_ЛЕГИВОН"
	if c.Sender().Username != "" {
		username = c.Sender().Username
	}
	// check if user have char in Party map
	val, ok := brain.Game.Party[c.Sender().ID]
	if ok {
		charlink := &val
		brain.Game.ActiveParty[c.Sender().ID] = charlink
		return c.Send(fmt.Sprintf("%s завербовался для активных действий", username))
	} else {
		return c.Send("вы, господин хороший, чара то сначала нарольте, а потом уже адвеньтюрьте так сказать")
	}
}

func DnDCharStats(c tele.Context, brain *BigBrain) (err error) {
	username := "АНОНИМ_ЛЕГИВОН"
	if c.Sender().Username != "" {
		username = c.Sender().Username
	}
	// check if user have char in Party map
	_, ok := brain.Game.Party[c.Sender().ID]
	if ok {
		charlink := brain.Game.Party[c.Sender().ID]
		message := fmt.Sprintf("%s ваш перец - %s\n❤️ = %d\n", username, string(charlink.Race)+"-"+string(charlink.Class), charlink.Hitpoints)
		_, ok := brain.Game.ActiveParty[c.Sender().ID]
		if ok {
			message += "вы также челен активной партии вот вам за это +1 хп\n"
			char := brain.Game.Party[c.Sender().ID]
			char.Hitpoints++
			brain.Game.Party[c.Sender().ID] = char
			message += "посмотрим что стало с хп в базе:\n"
			charlink := brain.Game.Party[c.Sender().ID]
			message += fmt.Sprintf("%s ваш перец - %s\n❤️ = %d\n", username, string(charlink.Race)+"-"+string(charlink.Class), charlink.Hitpoints)
		}
		return c.Send(message)
	} else {
		return c.Send("вы, господин хороший, чара то сначала нарольте, а потом уже адвеньтюрьте так сказать")
	}
}

func DnDCombat2(c tele.Context, serv *servitor.Servitor, brain *BigBrain) (err error) {
	message := brain.Game.Combat()
	message += "бой начался"
	c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, message)
	message = ""
	for brain.Game.CombatFlag {
		for _, char := range brain.Game.CombatOrder {
			if char.IsNPC {
				serv.Logger.Info("npc cycle start")
				// stupid way to pick target for npc
				for i, v := range brain.Game.ActiveParty {
					endFlag := true
					for validIndex, validTarget := range brain.Game.CombatOrder {
						serv.Logger.Info("combat order range over",
							slog.String("v.Name", v.Name),
							slog.String("validTarget.Name", validTarget.Name))
						if v.Name == validTarget.Name && validTarget.Hitpoints > 0 {
							endFlag = false
							message += fmt.Sprintf("%s бьёт🤜 по %s\n", char.Name, validTarget.Name)
							dmg, messagedmg := char.GetAttackDamage(validTarget.AC)
							message += messagedmg
							message += fmt.Sprintf("\nхп цели: %d - %d = %d\n", brain.Game.CombatOrder[validIndex].Hitpoints, dmg, brain.Game.CombatOrder[validIndex].Hitpoints-dmg)
							brain.Game.CombatOrder[validIndex].Hitpoints -= dmg
							brain.Game.Party[i] = *brain.Game.CombatOrder[validIndex]
							if brain.Game.CombatOrder[validIndex].Hitpoints <= 0 {
								message += "цель perished💀\n"
							}
						}
					}
					if endFlag {
						message += "пахот нпц винс"
						brain.Game.CombatFlag = false
					}
					c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, message)
					serv.Logger.Info("npc cycle comes to end, must be brake")
					break
				}
				c.Send(&tele.Chat{ID: c.Chat().ID}, message)
				message = ""
			} else {
				serv.Logger.Info("playuer cycle start")
				if char.Hitpoints <= 0 {
					continue
				}
				// stupid way to get player id
				var id int64
				for k, v := range brain.Game.Party {
					if v.Name == char.Name {
						id = k
					}
				}
				mes, buttons, _ := DnDTargetsButtonsPriv(c, serv, brain, c.Chat().ID)
				serv.Logger.Info("sending buttons to player")
				c.Bot().Send(&tele.Chat{ID: id}, mes, buttons)
				<-brain.Game.CombatFC
				serv.Logger.Info("returning from buttons part")
			}
		}
	}
	message += "комбату конец"
	return c.Send(message)
}

func DnDTargetsButtonsPriv(c tele.Context, serv *servitor.Servitor, brain *BigBrain, hostchat int64) (message string, buttons *tele.ReplyMarkup, err error) {
	message += "Выберите цель\n"
	incButtons := &tele.ReplyMarkup{ResizeKeyboard: true}

	var rows []tele.Row
	for id, item := range brain.Game.CombatOrder {
		rows = append(rows, incButtons.Row(incButtons.Data(fmt.Sprintf("%d: %s", id+1, item.Name), fmt.Sprintf("dndAttackTarget%d_%d", id, hostchat))))

	}
	rows = append(rows, incButtons.Row(incButtons.Data("скрыть", "sweep")))
	incButtons.Inline(rows...)
	return message, incButtons, nil
}

func DnDAttackByCallback(c tele.Context, serv *servitor.Servitor, brain *BigBrain, num int, chatID int64) (err error) {
	serv.Logger.Info("callback attack starts")
	if !brain.Game.CombatFlag {
		return c.Send("битва не может начатться, сделайте /dndcombat")
	}
	if num < 0 || num > len(brain.Game.CombatOrder) {
		return c.Send("номер залупатара вне грониц /dndcombat")
	}
	me := brain.Game.Party[c.Sender().ID]

	// assign target in order array (cringe)
	meIndex := 0
	for ind, char := range brain.Game.CombatOrder {
		if me.Name == char.Name {
			meIndex = ind
			break
		}
	}

	if brain.Game.CombatOrder[meIndex].Hitpoints <= 0 {
		return c.Send("вы мертвы и не можете совершать действия")
	}
	target := brain.Game.CombatOrder[num]
	var targetID int64 = 0
	for k, val := range brain.Game.Party {
		if target.Name == val.Name {
			targetID = k
		}
	}
	if target.Hitpoints <= 0 {
		return c.Send("ваша цель метрва")
	}
	brain.Game.CombatOrder[meIndex].Target = brain.Game.CombatOrder[num]
	message := ""
	if me.Name == target.Name {
		message += fmt.Sprintf("%s решил хуярит сам по себе (чистый термояд-дегенерат)\n", me.Name)
	} else {
		message += fmt.Sprintf("%s выбрал целью %s\n", me.Name, target.Name)
	}
	message += fmt.Sprintf("%s бьёт👊🏾 по %s\n", me.Name, target.Name)
	dmg, messagedmg := me.GetAttackDamage(target.AC)
	message += messagedmg
	message += fmt.Sprintf("\nхп цели: %d - %d = %d\n", target.Hitpoints, dmg, target.Hitpoints-dmg)
	target.Hitpoints -= dmg
	if !target.IsNPC {
		brain.Game.Party[targetID] = *target
	}
	if target.Hitpoints <= 0 {
		message += "цель perished💀\n"
		if target.Name == "Керилл" || target.Name == "Васян" {
			message += "Со смертью этого персонажа комбат в этом месте заканчивается, идите в другое или живите дальше в проклятом мире, который сами и создали\n"
			brain.Game.CombatFlag = false
		}
	}
	serv.Logger.Info("sending message after callback call to chat ")
	serv.Logger.Info(string(chatID))
	c.Bot().Send(&tele.Chat{ID: chatID}, message)
	brain.Game.CombatFC <- true
	return nil
}
