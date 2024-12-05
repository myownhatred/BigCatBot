package bigcat

import (
	"Guenhwyvar/servitor"
	"log/slog"
	"strings"

	tele "gopkg.in/telebot.v4"
)

func CallbackPikRouter(c tele.Context, serv *servitor.Servitor, brain *BigBrain) error {
	cbUniq := c.Callback().Data
	payloads := strings.Split(cbUniq, "\fpik")
	serv.Logger.Info("bigcat callbackspik", "router got next stuff as payload ", payloads[1])
	args := strings.Split(payloads[1], "_")
	serv.Logger.Info("bigcat callbackspik", "router got next stuff as args[0] ", args[0])
	if len(args) > 1 {
		serv.Logger.Info("bigcat callbackspik", "router got next stuff as args[1] ", args[1])
	}
	switch args[0] {
	case "WeekMenu":
		if brain.ChatContent[c.Chat().ID].LastPicture != "" {
			c.Delete()
			return CmdPikMenuWeek(c, serv, brain)
		}
		c.Delete()
		return c.Send("дайте картиночку в чятик")
	case "Week":
		if brain.ChatContent[c.Chat().ID].LastPicture == "" {
			c.Delete()
			return c.Send("дайте картиночку в чятик")
		}
		switch args[1] {
		case "Monday":
			return downloadToWeekDir(c, serv, brain, args[1])
		case "Tuesday":
			return downloadToWeekDir(c, serv, brain, args[1])
		case "Wednesday":
			return downloadToWeekDir(c, serv, brain, args[1])
		case "Thursday":
			return downloadToWeekDir(c, serv, brain, args[1])
		case "Friday":
			return downloadToWeekDir(c, serv, brain, args[1])
		case "Saturday":
			return downloadToWeekDir(c, serv, brain, args[1])
		case "Sunday":
			return downloadToWeekDir(c, serv, brain, args[1])
		case "Any":
			return downloadToWeekDir(c, serv, brain, args[1])
		case "Gm":
			return downloadToWeekDir(c, serv, brain, args[1])
		case "Angry":
			return downloadToWeekDir(c, serv, brain, args[1])
		}
	case "TGLink":
		if brain.ChatContent[c.Chat().ID].LastPicture != "" {
			c.Delete()
			return c.Send(brain.ChatContent[c.Chat().ID].LastPicture)
		}
		c.Delete()
		return c.Send("дайте картиночку в чятик")
	}
	return nil
}

func downloadToWeekDir(c tele.Context, serv *servitor.Servitor, brain *BigBrain, dir string) error {
	file, err := c.Bot().FileByID(brain.ChatContent[c.Chat().ID].LastPicture)
	if err != nil {
		c.Delete()
		serv.Logger.Error("bigcat", "CallbackPikRouter error getting file by ID from ChatContent", err.Error())
		return c.Send("ошибочка вышла:" + err.Error())
	}
	err = c.Bot().Download(&file, "./assets/week/"+dir+"/"+file.FileID+".jpg")
	if err != nil {
		serv.Logger.Error("bigcat", "CallbackPikRouter error downloading file to dir ",
			slog.String("dir name:", dir),
			slog.String("error:", err.Error()))
		c.Delete()
		return c.Send("ошибочка вышла:" + err.Error())
	}
	c.Delete()
	return c.Send("изображение добавлено в " + dir)
}
