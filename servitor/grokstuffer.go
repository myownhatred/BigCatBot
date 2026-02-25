package servitor

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"Guenhwyvar/bringer"
)

type GrokStufferServ struct {
	bringer bringer.Grokker
}

func NewGrokStufferServ(bringer bringer.Grokker) *GrokStufferServ {
	return &GrokStufferServ{
		bringer: bringer,
	}
}

func (g *GrokStufferServ) SimpleAnswer(prompt string) (string, string, error) {
	return g.bringer.SimpleAnswer(prompt)
}

func (g *GrokStufferServ) DNDBiogen(prompt string) (string, string, error) {
	return g.bringer.DNDBiogen(prompt)
}

func (g *GrokStufferServ) GenGrok(prompt, role string, temp float64) (string, string, error) {
	return g.bringer.GenGrok(prompt, role, temp)
}

func (g *GrokStufferServ) SendMessageInConversation(chatID int64, content string) (string, error) {
	return g.bringer.SendMessageInConversation(chatID, content)
}

func (g *GrokStufferServ) DeleteConversation(chatID int64) {
	g.bringer.DeleteConversation(chatID)
}

func (g *GrokStufferServ) CreateChatDayReport() (string, error) {
	// funky way to find last chatsummary file
	// TODO redo it to some nice way like date stuff or something
	entries, err := os.ReadDir("./chatlogs/")
	if err != nil {
		return "", err
	}

	var latestFile string
	var latestModTime time.Time
	var found bool

	// Iterate through files to find the most recently modified .json file
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue // Skip files with inaccessible metadata
		}

		modTime := info.ModTime()
		if !found || modTime.After(latestModTime) {
			latestFile = filepath.Join("./chatlogs/", entry.Name())
			latestModTime = modTime
			found = true
		}
	}
	// Check if a valid file was found
	if !found {
		return "", fmt.Errorf("no chatlogs found")
	}

	return g.bringer.AnalChatDay(latestFile)
}
