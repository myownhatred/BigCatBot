package bigcat

import (
	"fmt"
	"time"

	tele "gopkg.in/telebot.v4"
)

func (bh *BotHandler) TimersHandler(c tele.Context) error {
	list, err := bh.serv.GetTimeWithOutList(c.Chat().ID)
	if err != nil {
		return c.Send("не получилось со списком")
	}
	if len(list) == 0 {
		return c.Send("таймеров нет")
	}
	message := "Текущие таймеры:\n"
	for id, item := range list {
		duration := time.Since(item.Time)
		days := int(duration.Hours()) / 24
		hours := int(duration.Hours()) % 24
		message += fmt.Sprintf("%s: %s — %02d дейз %02d хаурс\n",
			butifulMumbers(id+1), item.Name, days, hours)
	}
	return c.Send(message)
}
