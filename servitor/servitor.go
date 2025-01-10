package servitor

import (
	"Guenhwyvar/bringer"
	"Guenhwyvar/config"
	"Guenhwyvar/entities"
	"Guenhwyvar/lib/citizen"
	"log/slog"
	"time"

	tele "gopkg.in/telebot.v4"
)

type WakaStuff interface {
	GetWakaStuff() string
}

type Memser interface {
	CreateGuiltyCatMeme(text string) (filePath string, err error)
}

type AnimeMaw interface {
	GetOpeningsFilePath() (filePath string)
	UploadOpenings(filePath string) (report string, err error)
	UploadOpeningsByURL(url string) (report string, err error)
	GetRandomOpening(typ string) (opening entities.AnimeOpening, err error)
}

type FreeMaw interface {
	GetFreeMaw(typ string) (maw entities.FreeMaw, err error)
	PutFreeMaw(maw entities.FreeMaw) (err error)
	GetFreeMawReport() (report string, err error)
	FreeMawVectorTypeAdd(qtype entities.VectorType) (err error)
	FreeMawVectorTypeByID(ID int) (qtype entities.VectorType, err error)
	FreeMawVectorTypes() (report []entities.VectorType, err error)
	FreeMawVectorAdd(vec entities.FreeVector) (err error)
	FreeMawVectorGetRandomByType(typ int) (vec entities.FreeVector, err error)
	FreeMawVectorUpsertScore(uid int64, vectorType int, score int) error
	FreeMawVectorGetTopScores(limit int) (scores []entities.VectorScore, err error)
}

type TimeWithOut interface {
	GetTimeWithOutList(chatID int64) (list []entities.TimeWithOut, err error)
	GetTimeWithOutTimerByID(id int) (event entities.TimeWithOut, err error)
	ResetTimer(id int) (err error)
	AddNewTimer(name string, chatID int64) (err error)
}

type Comfiger interface {
	GetAppComfig() (comfig *config.AppConfig, err error)
}

type Twitter interface {
	TwitterGetVideo(link string) (filePath string, err error)
}

type GetRekt interface {
	GetWeatherDayForecast(place string) (report string, err error)
	GetCurrentWeather(place string) (report string, err error)
	GetFreeSteamGames() (report string, err error)
	SendGenerationReq(modelID int, prompt string) (err error)
	GetGenerationStatus() (status string, err error)
	GetGeneratorStatus() (singa entities.Signa, err error)
}

type MediaCreator interface {
	MediaManulFile() (tele.File, error)
	MediaDayOfWeekFile() (tele.File, error)
	RandomFileFromDir(dirPath string) (string, error)
	GeneratorPickup() (file tele.File, err error)
}

type Police interface {
	UserDefaultCheck(UserID int64, username, firstname, lastname, command string) (ID int64, err error)
	MetatronChatAdd(ChatID int64, ChatName string) (err error)
	MetatronChatList() (IDs []int64, ChatIDs []int64, Names []string, err error)
	Achieves(GRID int) (IDs []int, GRIDs []int, Names []string, Ranks []int, Descrs []string, err error)
	UserAchs(UserID int64) (IDs []int, UserIDs []int64, AchIDs []int, Dates []time.Time, Chats []string, ChatIDs []int64, err error)
	UserAchAdd(UserID int64, AID int, chat string, chatID int64) (UAID int, err error)
	UserByUsername(username string) (UID int64, err error)
	GetAllUsers() (allUsers []citizen.Citizen, err error)
}

type Servitor struct {
	Logger *slog.Logger
	WakaStuff
	Memser
	AnimeMaw
	FreeMaw
	TimeWithOut
	Comfiger
	Twitter
	GetRekt
	MediaCreator
	Police
}

func NewServitor(bringer *bringer.Bringer, logger *slog.Logger) *Servitor {
	return &Servitor{
		Logger:       logger,
		WakaStuff:    NewWakaStuffServ(bringer.WakaStuff),
		Memser:       NewMemserServ(bringer.Memser),
		AnimeMaw:     NewAnimeMawServ(bringer.AnimeMaw),
		FreeMaw:      NewFreeMawServ(bringer.FreeMaw),
		TimeWithOut:  NewTimeWithOutServ(bringer.TimeWithOut),
		Comfiger:     NewComfigerServ(bringer.Comfiger),
		Twitter:      NewTwitterServ(bringer.Twitter),
		GetRekt:      NewGetRectServ(bringer.GetRekt),
		MediaCreator: NewMediaCreatorServ(bringer.Twitter, bringer.GetRekt, logger),
		Police:       NewPoliceServ(bringer.Police, logger),
	}
}
