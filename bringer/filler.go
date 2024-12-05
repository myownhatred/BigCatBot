package bringer

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Filler struct {
	logger *slog.Logger
}

func NewFiller(logger *slog.Logger) *Filler {
	return &Filler{logger: logger}
}

func (fil *Filler) MediaDayOfWeekFile() (file string, err error) {
	fil.logger.Info("morning picture for the day")
	fileD := filepath.FromSlash("./assets/good_morning_miserable.jpg")
	currentDay := time.Now().Weekday()
	switch currentDay {
	case time.Sunday:
		fil.logger.Info("case of Sunday")
		fileS, err := getRandomFileFromDir("./assets/week/Sunday")
		if err != nil {
			return fileD, err
		}
		fil.logger.Info("Sunday pic",
			slog.String("filename ", fileS),
		)
		file = filepath.FromSlash("./assets/week/Sunday/" + fileS)
		return file, nil
	case time.Monday:
		fil.logger.Info("case of Monday")
		fileS, err := getRandomFileFromDir("./assets/week/Monday")
		if err != nil {
			return fileD, err
		}
		fil.logger.Info("Monday pic",
			slog.String("filename ", fileS),
		)
		file = filepath.FromSlash("./assets/week/Monday/" + fileS)
		return file, nil
	case time.Tuesday:
		fil.logger.Info("case of Tuesday")
		fileS, err := getRandomFileFromDir("./assets/week/Tuesday")
		if err != nil {
			return fileD, err
		}
		fil.logger.Info("Tuesday pic",
			slog.String("filename ", fileS),
		)
		file = filepath.FromSlash("./assets/week/Tuesday/" + fileS)
		return file, nil
	case time.Wednesday:
		fil.logger.Info("case of Wednesday")
		fileS, err := getRandomFileFromDir("./assets/week/Wednesday")
		if err != nil {
			return fileD, err
		}
		fil.logger.Info("Wednesday pic",
			slog.String("filename ", fileS),
		)
		file = filepath.FromSlash("./assets/week/Wednesday/" + fileS)
		return file, nil
	case time.Thursday:
		fil.logger.Info("case of Thursday")
		fileS, err := getRandomFileFromDir("./assets/week/Thursday")
		if err != nil {
			return fileD, err
		}
		fil.logger.Info("Thursday pic",
			slog.String("filename ", fileS),
		)
		file = filepath.FromSlash("./assets/week/Thursday/" + fileS)
		return file, nil
	case time.Friday:
		fil.logger.Info("case of Friday")
		fileS, err := getRandomFileFromDir("./assets/week/Friday")
		if err != nil {
			return fileD, err
		}
		fil.logger.Info("Friday pic",
			slog.String("filename ", fileS),
		)
		file = filepath.FromSlash("./assets/week/Friday/" + fileS)
		return file, nil
	case time.Saturday:
		fil.logger.Info("case of Saturday")
		fileS, err := getRandomFileFromDir("./assets/week/Saturday")
		if err != nil {
			return fileD, err
		}
		fil.logger.Info("Saturday pic",
			slog.String("filename ", fileS),
		)
		file = filepath.FromSlash("./assets/week/Saturday/" + fileS)
		return file, nil
	default:
		return file, nil
	}
}

func (fil *Filler) GeneratorPickup() (file string, err error) {
	file = filepath.FromSlash("./" + "image.png")
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
