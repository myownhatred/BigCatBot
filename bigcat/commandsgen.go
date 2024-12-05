package bigcat

import (
	"Guenhwyvar/servitor"
	"fmt"

	tele "gopkg.in/telebot.v4"
)

func CmdGetGeneratedImage(c tele.Context, serv *servitor.Servitor) (err error) {
	m, _ := serv.MediaCreator.GeneratorPickup()
	pho := &tele.Photo{
		File:    m,
		Caption: "#aigenerash",
	}
	return c.Send(pho)
}

func CmdGetGeneratorStatus(c tele.Context, serv *servitor.Servitor) (err error) {
	s, err := serv.GetGeneratorStatus()
	if err != nil {
		return c.Send("беда:" + err.Error())
	}
	if s.Armed {
		report := fmt.Sprintf("Позывной генератора:%s", s.Callsign)
		incButtons := &tele.ReplyMarkup{ResizeKeyboard: true}

		var rows []tele.Row

		//return c.Send(message, incButtons)
		for _, m := range s.Models {
			rows = append(rows, incButtons.Row(incButtons.Data(m.Name, fmt.Sprintf("gen%d", m.ID))))
		}
		rows = append(rows, incButtons.Row(incButtons.Data("скрыть", "sweep")))

		incButtons.Inline(rows...)

		return c.Send(report, incButtons)
	} else {
		return c.Send("Нет активных генераторов")
	}
}
