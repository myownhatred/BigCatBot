package bigcat

import (
	"Guenhwyvar/config"
	dnd "Guenhwyvar/lib/DND"
	"Guenhwyvar/lib/memser"
	"Guenhwyvar/servitor"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	cron "github.com/robfig/cron/v3"
	tele "gopkg.in/telebot.v4"
)

type BigCat struct {
	tgBot    *tele.Bot
	serv     *servitor.Servitor
	comfig   *config.AppConfig
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

func New(tgBot *tele.Bot, comfig *config.AppConfig, serv *servitor.Servitor, str string, logger *slog.Logger) *BigCat {
	biggy := NewBigBrain()
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
	// user cache upload
	users, _ := c.serv.Police.GetAllUsers()
	for _, u := range users {
		c.bigBrain.Users[u.UserID] = u
	}
	// CRON JOBS
	// save summary
	c.clock.AddFunc("0 0 17 * * *", func() {
		err := c.serv.MemoryManager.SaveTheDay()
		if err != nil {
			c.logger.Error(fmt.Sprintf("error saving daily log on timer: %w", err))
		} else {
			c.logger.Info("got saved stuff nice")
		}
		// update grok limits
		for _, u := range c.bigBrain.Users {
			u.GrokToks = 10
			c.bigBrain.Users[u.UserID] = u
		}
	})

	// mobilizatsya
	c.clock.AddFunc("15 0 2 * * *", func() {
		//pik, err := memser.DaysMob()
		pik, err := memser.DaysToSomething()
		if err != nil {
			c.tgBot.Send(&tele.Chat{ID: c.comfig.ChatsAndPeps.MotherShip}, fmt.Sprintf("при созании картинки для будущего события случилась беда:%s", err.Error()))
		}
		m := &tele.Photo{
			File: tele.FromDisk(pik),
		}
		c.logger.Info(fmt.Sprintf("created pik %s", m.File))
		//c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, m)
	})
	// weather report
	c.clock.AddFunc("20 59 23 * * *", func() {
		c.bigBrain.Game = dnd.NewGame()
		report, err := c.serv.GetWeatherDayForecast("красноярск")
		if err != nil {
			c.tgBot.Send(&tele.Chat{ID: c.comfig.ChatsAndPeps.MotherShip}, err.Error())
		}
		c.tgBot.Send(&tele.Chat{ID: c.comfig.ChatsAndPeps.MotherShip}, report)
	})
	// c.clock.AddFunc("40 59 23 * * *", func() {
	// 	c.tgBot.Send(&tele.Chat{ID: c.comfig.ChatsAndPeps.MotherShip}, "порядок играния в днд - ролим чаров -> /dndjoin - вступить в комбат ЧВК редан (тока они будут драца)\n-> /dndmf - начнёца бой до победы кароч, выбор цели через приват")
	// })
	// manul spam
	c.clock.AddFunc("0 0 * * * *", func() {
		c.logger.Info("manul spam executed")
		m, cap, err := c.serv.MediaCreator.MediaManulFile()
		pho := &tele.Photo{
			File:    m,
			Caption: "#" + cap + "каждыйчас",
		}
		if err != nil {
			c.tgBot.Send(&tele.Chat{ID: c.comfig.ChatsAndPeps.MotherShip}, err.Error())
		}
		c.tgBot.Send(&tele.Chat{ID: c.comfig.ChatsAndPeps.MotherShip}, pho)
	})
	// gm spam
	c.clock.AddFunc("15 0 1 * * *", func() {
		//c.clock.AddFunc("0 * * * * *", func() {
		c.logger.Info("GM spam executed")
		m, err := c.serv.MediaCreator.MediaDayOfWeekFile()
		if err != nil {
			c.tgBot.Send(&tele.Chat{ID: c.comfig.ChatsAndPeps.MotherShip}, err.Error())
		}
		if strings.HasSuffix(m.FileLocal, ".mp4") {
			gif := &tele.Animation{
				File:     m,
				FileName: m.FileLocal,
				Caption:  "ДОБРОЕ УТРО",
			}
			c.tgBot.Send(&tele.Chat{ID: c.comfig.ChatsAndPeps.MotherShip}, gif)
			c.tgBot.Send(&tele.Chat{ID: c.comfig.ChatsAndPeps.JonChat}, gif)
		} else {
			pho := &tele.Photo{
				File:    m,
				Caption: "ДОБРОЕ УТРО",
			}
			c.tgBot.Send(&tele.Chat{ID: c.comfig.ChatsAndPeps.MotherShip}, pho)
			c.tgBot.Send(&tele.Chat{ID: c.comfig.ChatsAndPeps.JonChat}, pho)
		}
		c.logger.Info("anal of yesterday executed")
		report, err := c.serv.CreateChatDayReport()
		if err != nil {
			c.tgBot.Send(&tele.Chat{ID: c.comfig.ChatsAndPeps.MotherShip}, fmt.Sprintf("Проблемес с аналом дня: %v", err))
		}
		urlPart := strings.TrimPrefix(strconv.FormatInt(-c.comfig.ChatsAndPeps.MotherShip, 10), "100")
		report = strings.ReplaceAll(report, "msg://", "https://t.me/c/"+urlPart+"/")
		c.tgBot.Send(&tele.Chat{ID: c.comfig.ChatsAndPeps.MotherShip}, report)
	})

	// steam check
	// c.clock.AddFunc("0 0 3 * * *", func() {
	// 	c.logger.Info("steam spam executed")
	// 	report, err := c.serv.GetFreeSteamGames()
	// 	if err != nil {
	// 		c.tgBot.Send(&tele.Chat{ID: c.comfig.ChatsAndPeps.MotherShip}, err.Error())
	// 	}
	// 	c.tgBot.Send(&tele.Chat{ID: c.comfig.ChatsAndPeps.MotherShip}, report)
	// })
	c.clock.Start()
	c.tgBot.Start()
}
