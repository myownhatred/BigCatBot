package bringer

import (
	"Guenhwyvar/lib/memser"
	"log"
)

type MemserGG struct{}

func NewMemserGG() *MemserGG {
	return &MemserGG{}
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
