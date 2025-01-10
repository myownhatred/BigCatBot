package bigcat

import (
	dnd "Guenhwyvar/lib/DND"
	"Guenhwyvar/lib/memser"
	"Guenhwyvar/servitor"
	"fmt"
	"log/slog"
	"os"
	"strings"

	cron "github.com/robfig/cron/v3"
	tele "gopkg.in/telebot.v4"
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
	ManulSpam               bool
}

func New(tgBot *tele.Bot, serv *servitor.Servitor, str string, logger *slog.Logger) *BigCat {
	comfig, err := serv.GetAppComfig()
	biggy := NewBigBrain()
	if err != nil {
		logger.Error("cant load comfig to kitty:" + err.Error())
		os.Exit(1)
	}
	flag := &silly{
		ManulSpam: true,
	}
	handler := &BotHandler{
		tgbot:  tgBot,
		serv:   serv,
		flags:  flag,
		brain:  biggy,
		comfig: comfig,
		logger: logger,
	}
	handler.AddHandler()
	return &BigCat{
		tgBot:    tgBot,
		serv:     serv,
		brain:    flag,
		bigBrain: biggy,
		clock:    cron.New(cron.WithSeconds()),
		storage:  str,
		logger:   logger,
	}
}

func (c *BigCat) Start() {
	c.loadComfig()
	// user cache upload
	users, _ := c.serv.Police.GetAllUsers()
	for _, u := range users {
		c.bigBrain.Users[u.UserID] = u
	}
	// CRON JOBS
	// mobilizatsya
	c.clock.AddFunc("15 0 2 * * *", func() {
		//pik, err := memser.DaysMob()
		pik, err := memser.DaysToSomething()
		if err != nil {
			c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, fmt.Sprintf("при созании картинки для будущего события случилась беда:%s", err.Error()))
		}
		m := &tele.Photo{
			File: tele.FromDisk(pik),
		}
		c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, m)
	})
	// weather report
	c.clock.AddFunc("20 59 23 * * *", func() {
		c.bigBrain.Game = dnd.NewGame()
		report, err := c.serv.GetWeatherDayForecast("красноярск")
		if err != nil {
			c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, err.Error())
		}
		c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, report)
	})
	// c.clock.AddFunc("40 59 23 * * *", func() {
	// 	c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, "порядок играния в днд - ролим чаров -> /dndjoin - вступить в комбат ЧВК редан (тока они будут драца)\n-> /dndmf - начнёца бой до победы кароч, выбор цели через приват")
	// })
	// manul spam
	c.clock.AddFunc("0 0 * * * *", func() {
		c.logger.Info("manul spam executed")
		m, err := c.serv.MediaCreator.MediaManulFile()
		pho := &tele.Photo{
			File:    m,
			Caption: "#чтотокаждыйчас",
		}
		if err != nil {
			c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, err.Error())
		}
		c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, pho)
	})

	// gm spam
	c.clock.AddFunc("15 0 1 * * *", func() {
		//c.clock.AddFunc("0 * * * * *", func() {
		c.logger.Info("GM spam executed")
		m, err := c.serv.MediaCreator.MediaDayOfWeekFile()
		if err != nil {
			c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, err.Error())
		}
		if strings.HasSuffix(m.FileLocal, ".mp4") {
			gif := &tele.Animation{
				File:     m,
				FileName: m.FileLocal,
				Caption:  "ДОБРОЕ УТРО",
			}
			c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, gif)
			c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.JonChat}, gif)
		} else {
			pho := &tele.Photo{
				File:    m,
				Caption: "ДОБРОЕ УТРО",
			}
			c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, pho)
			c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.JonChat}, pho)
		}
	})

	// steam check
	// c.clock.AddFunc("0 0 3 * * *", func() {
	// 	c.logger.Info("steam spam executed")
	// 	report, err := c.serv.GetFreeSteamGames()
	// 	if err != nil {
	// 		c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, err.Error())
	// 	}
	// 	c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, report)
	// })
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
