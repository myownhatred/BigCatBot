package servitor

import (
	"Guenhwyvar/bringer"
	"io/ioutil"
	"log/slog"
	"math/rand"
	"time"

	tele "gopkg.in/telebot.v3"
)

type MediaCreatorServ struct {
	logger *slog.Logger
	twit   bringer.Twitter
	rekt   bringer.GetRekt
}

func NewMediaCreatorServ(twitter bringer.Twitter, rekter bringer.GetRekt, logger *slog.Logger) *MediaCreatorServ {
	return &MediaCreatorServ{
		logger: logger,
		twit:   twitter,
		rekt:   rekter,
	}
}

func (mc *MediaCreatorServ) MediaManulFile() (file tele.File, err error) {
	//TODO make it configurable
	sourses := [8]string{"redpandaeveryhr", "OtterAnHour", "FennecEveryHr",
		"PossumEveryHour", "ServalEveryHR", "raccoonhourly", "https://scryfall.com/random", "file/manyls"}
	rand.Seed(time.Now().Unix())
	toss := rand.Intn(len(sourses))
	mc.logger.Info("coin toss result",
		slog.Int("coin ", toss),
		slog.String("source ", sourses[toss]),
	)
	switch sourses[toss] {
	case "file/manyls":
		mc.logger.Info("case of manul")
		files, err := ioutil.ReadDir("./manyls")
		if err != nil {
			return file, err
		}
		rand.Seed(time.Now().UnixNano())
		randIndex := rand.Intn(len(files))
		mc.logger.Info("manul file pic",
			slog.String("filename ", files[randIndex].Name()),
		)
		file = tele.FromDisk("./manyls/" + files[randIndex].Name())
		return file, nil
	case "https://scryfall.com/random":
		mc.logger.Info("case of MTG")
		filePath, err := mc.rekt.GetRandomMTG()
		if err != nil {
			mc.logger.Warn("error getting MTG",
				slog.String("error message ", err.Error()),
			)
			return file, err
		}
		mc.logger.Info("MTG url acquired",
			slog.String("url ", filePath),
		)
		file = tele.FromURL(filePath)
		return file, nil
	default:
		mc.logger.Info("default case twitter",
			slog.String("source ", sourses[toss]),
		)
		filePath, err := mc.twit.TwitterGetHourlyPicture(sourses[toss])
		if err != nil {
			return file, err
		}
		file = tele.FromURL(filePath)
		return file, nil
	}

}
