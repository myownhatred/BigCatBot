package servitor

import (
	"Guenhwyvar/bringer"
	"fmt"
	"io/ioutil"
	"log/slog"
	"math/rand"
	"os"
	"time"

	tele "gopkg.in/telebot.v4"
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

// TODO move file operations to bringer
func (mc *MediaCreatorServ) MediaDayOfWeekFile() (file tele.File, err error) {
	mc.logger.Info("morning picture for the day")
	fileD := tele.FromDisk("./assets/good_morning_miserable.jpg")
	currentDay := time.Now().Weekday()
	switch currentDay {
	case time.Sunday:
		mc.logger.Info("case of Sunday")
		fileS, err := getRandomFileFromDir("./assets/week/Sunday")
		if err != nil {
			return fileD, err
		}
		mc.logger.Info("Sunday pic",
			slog.String("filename ", fileS),
		)
		file = tele.FromDisk("./assets/week/Sunday/" + fileS)
		return file, nil
	case time.Monday:
		mc.logger.Info("case of Monday")
		fileS, err := getRandomFileFromDir("./assets/week/Monday")
		if err != nil {
			return fileD, err
		}
		mc.logger.Info("Monday pic",
			slog.String("filename ", fileS),
		)
		file = tele.FromDisk("./assets/week/Monday/" + fileS)
		return file, nil
	case time.Tuesday:
		mc.logger.Info("case of Tuesday")
		fileS, err := getRandomFileFromDir("./assets/week/Tuesday")
		if err != nil {
			return fileD, err
		}
		mc.logger.Info("Tuesday pic",
			slog.String("filename ", fileS),
		)
		file = tele.FromDisk("./assets/week/Tuesday/" + fileS)
		return file, nil
	case time.Wednesday:
		mc.logger.Info("case of Wednesday")
		fileS, err := getRandomFileFromDir("./assets/week/Wednesday")
		if err != nil {
			return fileD, err
		}
		mc.logger.Info("Wednesday pic",
			slog.String("filename ", fileS),
		)
		file = tele.FromDisk("./assets/week/Wednesday/" + fileS)
		return file, nil
	case time.Thursday:
		mc.logger.Info("case of Thursday")
		fileS, err := getRandomFileFromDir("./assets/week/Thursday")
		if err != nil {
			return fileD, err
		}
		mc.logger.Info("Thursday pic",
			slog.String("filename ", fileS),
		)
		file = tele.FromDisk("./assets/week/Thursday/" + fileS)
		return file, nil
	case time.Friday:
		mc.logger.Info("case of Friday")
		fileS, err := getRandomFileFromDir("./assets/week/Friday")
		if err != nil {
			return fileD, err
		}
		mc.logger.Info("Friday pic",
			slog.String("filename ", fileS),
		)
		file = tele.FromDisk("./assets/week/Friday/" + fileS)
		return file, nil
	case time.Saturday:
		mc.logger.Info("case of Saturday")
		fileS, err := getRandomFileFromDir("./assets/week/Saturday")
		if err != nil {
			return fileD, err
		}
		mc.logger.Info("Saturday pic",
			slog.String("filename ", fileS),
		)
		file = tele.FromDisk("./assets/week/Saturday/" + fileS)
		return file, nil
	default:
		return file, nil
	}
}

func (mc *MediaCreatorServ) MediaManulFile() (file tele.File, caption string, err error) {
	//TODO make it configurable
	// twi - "OtterAnHour" "raccoonhourly"
	// sourses := [8]string{"redpandaeveryhr", "FennecEveryHr",
	// 	"PossumEveryHour", "ServalEveryHR", "https://scryfall.com/random",
	// 	"file/manyls", "file/nintendo", "file/japan"}
	sourses := [5]string{"https://scryfall.com/random",
		"file/manyls", "file/bigcats", "file/umumusume", "file/endfield"}
	toss := rand.Intn(len(sourses))
	mc.logger.Info("coin toss result",
		slog.Int("coin ", toss),
		slog.String("source ", sourses[toss]),
	)
	switch sourses[toss] {
	case "file/manyls":
		capt := "манул"
		mc.logger.Info("case of manul")
		files, err := ioutil.ReadDir("./manyls")
		if err != nil {
			return file, capt, err
		}
		rand.Seed(time.Now().UnixNano())
		randIndex := rand.Intn(len(files))
		mc.logger.Info("manul file pic",
			slog.String("filename ", files[randIndex].Name()),
		)
		file = tele.FromDisk("./manyls/" + files[randIndex].Name())
		return file, capt, nil
	case "file/bigcats":
		capt := "бигкэт"
		mc.logger.Info("case of bigcats")
		files, err := ioutil.ReadDir("./bigcats")
		if err != nil {
			return file, capt, err
		}
		randIndex := rand.Intn(len(files))
		mc.logger.Info("bigcats file pic",
			slog.String("filename ", files[randIndex].Name()),
		)
		file = tele.FromDisk("./bigcats/" + files[randIndex].Name())
		return file, capt, nil
	case "file/umumusume":
		capt := "умамусума"
		mc.logger.Info("case of umumusume")
		files, err := ioutil.ReadDir("./umumusume")
		if err != nil {
			return file, capt, err
		}
		randIndex := rand.Intn(len(files))
		mc.logger.Info("umumusume file pic",
			slog.String("filename ", files[randIndex].Name()),
		)
		file = tele.FromDisk("./umumusume/" + files[randIndex].Name())
		return file, capt, nil
	case "file/endfield":
		capt := "эндфилд"
		mc.logger.Info("case of endfield")
		files, err := ioutil.ReadDir("./endfield")
		if err != nil {
			return file, capt, err
		}
		randIndex := rand.Intn(len(files))
		mc.logger.Info("endfield file pic",
			slog.String("filename ", files[randIndex].Name()),
		)
		file = tele.FromDisk("./endfield/" + files[randIndex].Name())
		return file, capt, nil
	case "https://scryfall.com/random":
		mc.logger.Info("case of MTG")
		filePath, err := mc.rekt.GetRandomMTG()
		if err != nil {
			mc.logger.Warn("error getting MTG",
				slog.String("error message ", err.Error()),
			)
			return file, "мтг", err
		}
		mc.logger.Info("MTG url acquired",
			slog.String("url ", filePath),
		)
		file = tele.FromURL(filePath)
		return file, "мтг", nil
	default:
		mc.logger.Info("default case twitter",
			slog.String("source ", sourses[toss]),
		)
		filePath, err := mc.twit.TwitterGetHourlyPicture(sourses[toss])
		if err != nil {
			return file, "твиттер", err
		}
		file = tele.FromURL(filePath)
		return file, "твиттер", nil
	}

}

func (mc *MediaCreatorServ) GeneratorPickup() (file tele.File, err error) {
	file = tele.FromDisk("./image.png")
	return file, nil
}

func getRandomFileFromDir(dirPath string) (string, error) {
	// Read the directory and get a list of files
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	// Filter out non-files (directories, etc.)
	var validFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() {
			validFiles = append(validFiles, file)
		}
	}

	// Check if we have any valid files
	if len(validFiles) == 0 {
		return "", fmt.Errorf("no files found in directory: %s", dirPath)
	}

	// Seed the random number generator

	// Pick a random file
	randomIndex := rand.Intn(len(validFiles))
	randomFile := validFiles[randomIndex]

	return randomFile.Name(), nil
}

func (mc *MediaCreatorServ) RandomFileFromDir(dirPath string) (string, error) {
	// Read the directory and get a list of files
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	// Filter out non-files (directories, etc.)
	var validFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() {
			validFiles = append(validFiles, file)
		}
	}

	// Check if we have any valid files
	if len(validFiles) == 0 {
		return "", fmt.Errorf("no files found in directory: %s", dirPath)
	}

	// Seed the random number generator

	// Pick a random file
	randomIndex := rand.Intn(len(validFiles))
	randomFile := validFiles[randomIndex]

	return randomFile.Name(), nil
}
