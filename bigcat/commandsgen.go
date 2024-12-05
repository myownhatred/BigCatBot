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
		report := fmt.Sprintf("Генератор активен\nПозывной:%s\nМодели:\n", s.Callsign)
		for _, m := range s.Models {
			report += fmt.Sprintf("%d - %s\n", m.ID, m.Name)
		}
		return c.Send(report)
	} else {
		return c.Send("Нет активных генераторов")
	}

}
