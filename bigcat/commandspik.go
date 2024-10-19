package bigcat

import (
	"Guenhwyvar/servitor"
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func CmdPikMenuMain(c tele.Context) (err error) {
	message := "что сделать с картинкой?\n"
	incButtons := &tele.ReplyMarkup{ResizeKeyboard: true}

	var rows []tele.Row
	//rows = append(rows, incButtons.Row(incButtons.Data("АИ", "pikAI")))
	rows = append(rows, incButtons.Row(incButtons.Data("Добавление пыхчи в зокрома", "pikWeekMenu")))
	rows = append(rows, incButtons.Row(incButtons.Data("Сылка", "pikTGLink")))

	rows = append(rows, incButtons.Row(incButtons.Data("скрыть", "sweep")))
	incButtons.Inline(rows...)
	return c.Send(message, incButtons)
}

// add pik to asset dir
func CmdPikMenuWeek(c tele.Context, serv *servitor.Servitor, brain *BigBrain) (err error) {
	message := "какого типа ваша ежедневная картинка?\n"
	incButtons := &tele.ReplyMarkup{ResizeKeyboard: true}

	var rows []tele.Row
	rows = append(rows, incButtons.Row(incButtons.Data("Понедельник", fmt.Sprintf("pikWeek_Monday_%d", c.Sender().ID))))
	rows = append(rows, incButtons.Row(incButtons.Data("Вторник", fmt.Sprintf("pikWeek_Tuesday_%d", c.Sender().ID))))
	rows = append(rows, incButtons.Row(incButtons.Data("Среда", fmt.Sprintf("pikWeek_Wednesday_%d", c.Sender().ID))))
	rows = append(rows, incButtons.Row(incButtons.Data("Четверг", fmt.Sprintf("pikWeek_Thursday_%d", c.Sender().ID))))
	rows = append(rows, incButtons.Row(incButtons.Data("Пятница", fmt.Sprintf("pikWeek_Friday_%d", c.Sender().ID))))
	rows = append(rows, incButtons.Row(incButtons.Data("Суббота", fmt.Sprintf("pikWeek_Saturday_%d", c.Sender().ID))))
	rows = append(rows, incButtons.Row(incButtons.Data("Воскресенье", fmt.Sprintf("pikWeek_Sunday_%d", c.Sender().ID))))
	rows = append(rows, incButtons.Row(incButtons.Data("Любой", fmt.Sprintf("pikWeek_Any_%d", c.Sender().ID))))
	rows = append(rows, incButtons.Row(incButtons.Data("Доброе утро", fmt.Sprintf("pikWeek_Gm_%d", c.Sender().ID))))
	rows = append(rows, incButtons.Row(incButtons.Data("Недоброе утро", fmt.Sprintf("pikWeek_Angry_%d", c.Sender().ID))))

	rows = append(rows, incButtons.Row(incButtons.Data("скрыть", "sweep")))
	incButtons.Inline(rows...)
	return c.Send(message, incButtons)
}

func CmdPikWeekTest(c tele.Context, serv *servitor.Servitor, brain *BigBrain) (err error) {
	fileS, err := serv.RandomFileFromDir("./assets/week/Test")
	if err != nil {
		return err
	}
	serv.Logger.Info("bigcat commandspik", "got next filename from dir ", fileS)
	file := tele.FromDisk("./assets/week/Test/" + fileS)
	if strings.HasSuffix(fileS, ".mp4") {
		gif := &tele.Animation{
			File:     file,
			FileName: file.FileLocal,
			Caption:  "test gifky",
		}
		return c.Send(gif)
	}
	pho := &tele.Photo{
		File:    file,
		Caption: "тест воскресенья",
	}
	return c.Send(pho)
}
