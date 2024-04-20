package bigcat

import (
	"Guenhwyvar/lib/memser"
	"Guenhwyvar/servitor"
	"fmt"
	"log/slog"
	"os"

	cron "github.com/robfig/cron/v3"
	tele "gopkg.in/telebot.v3"
)

type BigCat struct {
	tgBot    *tele.Bot
	serv     *servitor.Servitor
	bigBrain *BigBrain
	brain    *silly
	clock    *cron.Cron
	storage  string
	logger   *slog.Logger
}

type silly struct {
	AnimeOpeningsUploadFlag bool
}

func New(tgBot *tele.Bot, serv *servitor.Servitor, str string, logger *slog.Logger) *BigCat {
	comfig, err := serv.GetAppComfig()
	if err != nil {
		logger.Error("cant load comfig to kitty:%s", err.Error())
		os.Exit(1)
	}
	flag := &silly{}
	handler := &BotHandler{
		tgbot:  tgBot,
		serv:   serv,
		flags:  flag,
		comfig: comfig,
		logger: logger,
	}
	handler.AddHandler()
	return &BigCat{
		tgBot:    tgBot,
		serv:     serv,
		brain:    flag,
		bigBrain: &BigBrain{},
		clock:    cron.New(cron.WithSeconds()),
		storage:  str,
		logger:   logger,
	}
}

func (c *BigCat) Start() {
	c.loadComfig()
	// CRON JOBS
	// mobilizatsya
	c.clock.AddFunc("15 0 2 * * *", func() {
		pik, err := memser.DaysMob()
		if err != nil {
			c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, fmt.Sprintf("при созании картинки для сбросика таймер случчилась бида:%s", err.Error()))
		}
		m := &tele.Photo{
			File: tele.FromDisk(pik),
		}
		c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, m)
	})
	// weather report
	c.clock.AddFunc("20 59 23 * * *", func() {
		report, err := c.serv.GetWeatherDayForecast("красноярск")
		if err != nil {
			c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, err.Error())
		}
		c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, report)
	})
	c.clock.Start()
	c.tgBot.Start()
}

func (c *BigCat) loadComfig() {
	comfig, err := c.serv.GetAppComfig()
	if err != nil {
		c.logger.Error("cant load comfig to kitty:%s", err.Error())
		os.Exit(1)
	} else {
		c.bigBrain.Comfig = *comfig
	}
}
