package bigcat

import (
	"Guenhwyvar/servitor"
	"fmt"

	tele "gopkg.in/telebot.v4"
)

func CmdCitizensAllSimpe(c tele.Context, brain *BigBrain) error {
	report := ""

	for _, v := range brain.Users {
		report += fmt.Sprintf("%d - %s\n", v.UserID, v.Username)
	}
	return c.Send(report)
}

func CmdCitizensAllBase(c tele.Context, serv *servitor.Servitor) (err error) {
	report := ""

	u, err := serv.Police.GetAllUsers()
	if err != nil {
		return c.Send("случилась бида при доставании юзеров из базы " + err.Error())
	}
	for _, v := range u {
		report += fmt.Sprintf("%d - %s\n", v.UserID, v.Username)
	}
	return c.Send(report)
}
