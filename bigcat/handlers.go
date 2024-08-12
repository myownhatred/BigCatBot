package bigcat

import (
	"Guenhwyvar/config"
	"Guenhwyvar/servitor"
	"fmt"
	"log/slog"

	tele "gopkg.in/telebot.v3"
)

var (
	menu     = &tele.ReplyMarkup{ResizeKeyboard: true, RemoveKeyboard: true}
	selector = &tele.ReplyMarkup{}

	btnHelp     = menu.Text("Помось")
	btnAnimeMaw = menu.Text("Анимемав")
	btnSettings = menu.Text("всякое")
	btnClear    = menu.Text("хз что это")

	mawButs  = &tele.ReplyMarkup{}
	btnAni   = mawButs.Data("AnimeMaw", "animemaw", "anime")
	btnGrob  = mawButs.Data("GrobMaw", "grobmaw", "grob")
	btnDummy = mawButs.Data("Dummy", "dummy")
)

type BotHandler struct {
	tgbot  *tele.Bot
	serv   *servitor.Servitor
	flags  *silly
	brain  *BigBrain
	comfig *config.AppConfig
	logger *slog.Logger
}

func (bh *BotHandler) AddHandler() {
	// main chat to send warning to
	mothership := tele.Chat{
		ID: bh.comfig.MotherShip,
	}

	bh.tgbot.Handle(tele.OnText, func(c tele.Context) error {
		return CommandHandler(c, bh.serv, bh.flags, bh.brain, bh.comfig, bh.logger)
	})
	bh.tgbot.Handle("/start", func(c tele.Context) error {
		return c.Send("пливет!", menu)
	})
	bh.tgbot.Handle(&btnAnimeMaw, func(c tele.Context) error {
		return c.Send("вот тебе немного помощи")
	})
	bh.tgbot.Handle(&btnSettings, func(c tele.Context) error {
		return c.Send("ещё не придумал")
	})
	bh.tgbot.Handle(&btnClear, func(c tele.Context) error {
		menu.RemoveKeyboard = true
		menu = &tele.ReplyMarkup{}
		return c.Send("случилось страшное")
	})

	bh.tgbot.Handle(tele.OnCallback, func(c tele.Context) error {
		return CallbackHandler(c, bh.serv)
	})

	bh.tgbot.Handle(tele.OnAddedToGroup, func(c tele.Context) error {
		_, err := bh.tgbot.Send(&mothership, fmt.Sprintf("Трявога! 🤬 \nБлядина @%s добавил меня в чат %s, примите мэры!", c.Message().Sender.Username, c.Message().Chat.Title))
		return err
	})

	bh.tgbot.Handle(&btnGrob, func(c tele.Context) error {
		_ = c.Respond(&tele.CallbackResponse{})
		opening, err := bh.serv.GetRandomOpening("grob")
		if err != nil {
			return c.Send(fmt.Sprintf("Неполучилось с опенингом: %s", err.Error()))
		}
		return c.Send(fmt.Sprintf("Тебе выпало послушать %s - %s", opening.Description, opening.Link))
	})
	bh.tgbot.Handle(&btnDummy, func(c tele.Context) error {
		_ = c.Respond(&tele.CallbackResponse{})
		return c.Send("копочка-наёбочка")
	})
	bh.tgbot.Handle(tele.OnPhoto, func(c tele.Context) error {
		// Metatron checks and actions
		// private chat only
		if c.Message().Private() {
			// check if user is on the bot/metatron list and set metatron flag on
			if _, ok := bh.brain.UsersFlags[c.Sender().ID]; ok {
				bh.logger.Info("user found:", c.Sender().ID)
			} else {
				bh.logger.Info("user not found:", c.Sender().ID)
				return nil
			}
			val := bh.brain.UsersFlags[c.Sender().ID]
			if val.MetatronFordwardFlag {
				return c.ForwardTo(&tele.Chat{ID: val.MetatronChat})
			}
		}
		return nil
	})
	bh.tgbot.Handle(tele.OnVideo, func(c tele.Context) error {
		// Metatron checks and actions
		// private chat only
		if c.Message().Private() {
			// check if user is on the bot/metatron list and set metatron flag on
			if _, ok := bh.brain.UsersFlags[c.Sender().ID]; ok {
				bh.logger.Info("user found:", c.Sender().ID)
			} else {
				bh.logger.Info("user not found:", c.Sender().ID)
				return nil
			}
			val := bh.brain.UsersFlags[c.Sender().ID]
			if val.MetatronFordwardFlag {
				return c.ForwardTo(&tele.Chat{ID: val.MetatronChat})
			}
		}
		return nil
	})
	bh.tgbot.Handle(tele.OnDocument, func(c tele.Context) error {
		//TODO: rewrite for good
		if !bh.flags.AnimeOpeningsUploadFlag {
			return nil
		}
		bh.flags.AnimeOpeningsUploadFlag = false
		size := c.Message().Document.FileSize
		user := c.Message().Sender.Username
		fileTele := &c.Message().Document.File
		msg := fmt.Sprintf("о ти (%s) прислал файлек и он весет %d байд!\nflag to opening is :%v\n", user, size, bh.flags.AnimeOpeningsUploadFlag)
		//_ c.Send(msg)
		_ = bh.tgbot.Download(fileTele, "openingsfile.csv")
		report, _ := bh.serv.UploadOpenings("openingsfile.csv")
		msg += report
		return c.Send(msg)
	})

}
