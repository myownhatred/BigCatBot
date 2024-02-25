package servitor

import (
	"Guenhwyvar/bringer"
	"Guenhwyvar/config"
	"Guenhwyvar/entities"
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

type Servitor struct {
	WakaStuff
	Memser
	AnimeMaw
	FreeMaw
	TimeWithOut
	Comfiger
}

func NewServitor(bringer *bringer.Bringer) *Servitor {
	return &Servitor{
		WakaStuff:   NewWakaStuffServ(bringer.WakaStuff),
		Memser:      NewMemserServ(bringer.Memser),
		AnimeMaw:    NewAnimeMawServ(bringer.AnimeMaw),
		FreeMaw:     NewFreeMawServ(bringer.FreeMaw),
		TimeWithOut: NewTimeWithOutServ(bringer.TimeWithOut),
		Comfiger:    NewComfigerServ(bringer.Comfiger),
	}
}
