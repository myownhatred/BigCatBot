package bigcat

import (
	"Guenhwyvar/servitor"
	"log"

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
}

type silly struct {
	AnimeOpeningsUploadFlag bool
}

func New(tgBot *tele.Bot, serv *servitor.Servitor, str string) *BigCat {
	comfig, err := serv.GetAppComfig()
	if err != nil {
		log.Fatalf("cant load comfig to kitty:%s", err.Error())
	}
	flag := &silly{}
	handler := &BotHandler{
		tgbot:  tgBot,
		serv:   serv,
		flags:  flag,
		comfig: comfig,
	}
	handler.AddHandler()
	return &BigCat{
		tgBot:    tgBot,
		serv:     serv,
		brain:    flag,
		bigBrain: &BigBrain{},
		clock:    cron.New(cron.WithSeconds()),
		storage:  str,
	}
}

func (c *BigCat) Start() {
	c.loadComfig()
	//c.clock.AddFunc("* * * * ", func(){
	//})
	//c.clock.AddFunc("*/15 * * * * *", func() {
	//	c.tgBot.Send(&tele.Chat{ID: c.bigBrain.Comfig.MotherShip}, "начинаем набалтывать")
	//})
	//c.clock.AddFunc("*/15 * * * * *", func() { fmt.Println("miun") })
	c.tgBot.Start()
}

func (c *BigCat) loadComfig() {
	comfig, err := c.serv.GetAppComfig()
	if err != nil {
		log.Fatalf("cant load comfig to kitty:%s", err.Error())
	} else {
		c.bigBrain.Comfig = *comfig
	}
}
