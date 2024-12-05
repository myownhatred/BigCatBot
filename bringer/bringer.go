package bringer

import (
	"Guenhwyvar/config"
	"Guenhwyvar/entities"
	"database/sql"
	"log/slog"
	"time"

	"github.com/go-resty/resty/v2"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	twitterscraper "github.com/imperatrona/twitter-scraper"
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
	FreeMawVectorTypeAdd(qtype entities.VectorType) (err error)
	FreeMawVectorTypeByID(ID int) (qtype entities.VectorType, err error)
	FreeMawVectorTypes() (report []entities.VectorType, err error)
	FreeMawVectorAdd(vec entities.FreeVector) (err error)
	FreeMawVectorGetRandomByType(typ int) (vec entities.FreeVector, err error)
}

type Comfiger interface {
	GetAppComfigFromViper() (config *config.AppConfig, err error)
}

type Twitter interface {
	TwitterGetVideo(link string) (filePath string, err error)
	TwitterGetHourlyPicture(acc string) (filePath string, err error)
}

type GetRekt interface {
	GetWeatherDayForecast(place string) (report string, err error)
	GetCurrentWeather(place string) (report string, err error)
	GetRandomMTG() (url string, err error)
	GetFreeSteamGames() (string, error)
	SendGenerationReq(modelID int, prompt string) (err error)
	GetGenerationStatus() (status string, err error)
	GetGeneratorStatus() (singa entities.Signa, err error)
}

type Police interface {
	UserDefaultCheck(UserID int64, username, firstname, lastname, command string) (ID int64, err error)
	MetatronChatAdd(ChatID int64, ChatName string) (err error)
	MetatronChatList() (IDs []int64, ChatIDs []int64, Names []string, err error)
	Achieves(GRID int) (IDs []int, GRIDs []int, Names []string, Ranks []int, Descrs []string, err error)
	AchGroups() (IDs []int, GroupNames []string, err error)
	UserAchs(UserID int64) (IDs []int, UserIDs []int64, AchIDs []int, Dates []time.Time, Chats []string, ChatIDs []int64, err error)
	UserAchAdd(UserID int64, AID int, chat string, chatID int64) (UAID int, err error)
	UserByUsername(username string) (UID int64, err error)
}

type Bringer struct {
	gormPost *gorm.DB
	db       *sql.DB
	logger   *slog.Logger
	WakaStuff
	Memser
	AnimeMaw
	TimeWithOut
	FreeMaw
	Comfiger
	Twitter
	GetRekt
	Police
}

func NewBringer(r *resty.Client, scrap *twitterscraper.Scraper, v *viper.Viper, db *sql.DB, logger *slog.Logger) *Bringer {
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
		logger:      logger,
		Comfiger:    NewComfigerViper(v),
		WakaStuff:   NewWakaStuff(r, v),
		Memser:      NewMemserGG(),
		AnimeMaw:    NewAnimeMawPostgres(gormP, r),
		TimeWithOut: NewTimeWithOutPostgres(db),
		FreeMaw:     NewFreeMawPostgres(db, logger),
		Twitter:     NewTwitterScrapper(scrap),
		GetRekt:     NewGetRect(r, v),
		Police:      NewPolicePostgres(db, logger),
	}
}
