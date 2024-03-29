package bringer

import (
	"Guenhwyvar/config"
	"Guenhwyvar/entities"
	"database/sql"

	"github.com/go-resty/resty/v2"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

type WakaStuff interface {
	GetDailyWaka() (reportMessage string, err error)
}

type Memser interface {
	CreateGuiltyCatMeme(text string) (filePath string, err error)
	CreateHoldMeme(text string) (filePath string, err error)
}

type AnimeMaw interface {
	GetOpeningsFromDB() (openings []entities.AnimeOpening, err error)
	GetOpeningsFromURL(url string) (openings []entities.AnimeOpening, err error)
	PutOpeningsToDB(openings []entities.AnimeOpening) (affRows int64, err error)
	GetRandomOpening(typ string) (opening entities.AnimeOpening, err error)
}

type TimeWithOut interface {
	GetTimeWithOutList(chatID int64) (list []entities.TimeWithOut, err error)
	GetTimeWithOutTimerByID(id int) (event entities.TimeWithOut, err error)
	ResetTimer(id int) (err error)
	AddNewTimer(name string, chatID int64) (err error)
}

type FreeMaw interface {
	GetRandomMawFromDB(typ string) (maw entities.FreeMaw, err error)
	PutFreeMawToDB(maw entities.FreeMaw) (err error)
	FreeMawDBReport() (report string, err error)
}

type Comfiger interface {
	GetAppComfigFromViper() (config *config.AppConfig, err error)
}

type Twitter interface {
	TwitterGetVideo(link string) (filePath string, err error)
}

type Bringer struct {
	gormPost *gorm.DB
	db       *sql.DB
	WakaStuff
	Memser
	AnimeMaw
	TimeWithOut
	FreeMaw
	Comfiger
	Twitter
}

func NewBringer(r *resty.Client, scrap *twitterscraper.Scraper, v *viper.Viper, db *sql.DB) *Bringer {
	gormP, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		// TODO: wrap it better - could work without DB
		panic(err)
	}
	return &Bringer{
		gormPost:    gormP,
		db:          db,
		Comfiger:    NewComfigerViper(v),
		WakaStuff:   NewWakaStuff(r, v),
		Memser:      NewMemserGG(),
		AnimeMaw:    NewAnimeMawPostgres(gormP, r),
		TimeWithOut: NewTimeWithOutPostgres(db),
		FreeMaw:     NewFreeMawPostgres(db),
		Twitter:     NewTwitterScrapper(scrap),
	}
}
