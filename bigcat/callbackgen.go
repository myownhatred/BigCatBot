package bigcat

import (
	"Guenhwyvar/servitor"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

func CallbackGenerateImage(c tele.Context, serv *servitor.Servitor, brain *BigBrain) error {
	cbUniq := c.Callback().Data
	payloads := strings.Split(cbUniq, "\fgen")
	// payloads = gen%d - %d is gen model number
	serv.Logger.Info("bigcat callbacksgen", "got next stuff as payload ", payloads[1])
	c.Delete()
	num, err := strconv.Atoi(payloads[1])
	if err != nil {
		c.Delete()
		return c.Send("проблема с конвертацией номера модели " + payloads[1])
	}
	brain.GenTrapMap[c.Sender().ID] = GeneratorTrap{
		UID:     c.Sender().ID,
		ChatID:  c.Chat().ID,
		ModelID: num,
	}

	return c.Send("введите промпт для генерации")
}
