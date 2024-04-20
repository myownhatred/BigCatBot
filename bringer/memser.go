package bringer

import (
	"Guenhwyvar/lib/memser"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

type MemserGG struct{}

func NewMemserGG() *MemserGG {
	return &MemserGG{}
}

func (m *MemserGG) GetDayPicture(day string) (filePath string, err error) {
	rand.Seed(time.Now().UnixNano())
	files, err := ioutil.ReadDir(day)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", nil
	}
	randIndex := rand.Intn(len(files))
	randomFile := files[randIndex].Name()
	return randomFile, nil
}

func (m *MemserGG) CreateGuiltyCatMeme(text string) (filePath string, err error) {
	// TODO: make filepath part of configuration or some sort of
	pathImg, err := memser.TextOnImg(text)
	if err != nil {
		log.Printf("failed to generate meme %v\n", err)
		pathImg = "storage/notponyal.png"
		return "", err
	}

	return pathImg, nil
}

func (m *MemserGG) CreateHoldMeme(text string) (filePath string, err error) {
	// TODO: make filepath part of configuration
	pathImg, err := memser.HoldMeme(text)
	if err != nil {
		log.Printf("failed to generate hold meme %v\n", err)
		pathImg = "storage/hold.png"
		return "", err
	}

	return pathImg, nil
}
