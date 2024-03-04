package bringer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

type TwitterScrapper struct {
	s *twitterscraper.Scraper
}

func NewTwitterScrapper(s *twitterscraper.Scraper) *TwitterScrapper {
	return &TwitterScrapper{
		s: s,
	}
}

func (s *TwitterScrapper) TwitterGetVideo(link string) (filePath string, err error) {

	tweet, err := s.s.GetTweet(link)
	if err != nil {
		return "", err
	}

	log.Println(tweet.Text)
	vidlinl := ""
	if len(tweet.Videos) > 0 {
		vidlinl = tweet.Videos[0].URL
	} else {
		return "", fmt.Errorf("у твиттика нету видосиков")
	}
	out, err := os.Create(link + ".mp4")
	if err != nil {
		return "", fmt.Errorf("не могу саздать файлик с твиттиком: %s", err.Error())
	}
	defer out.Close()
	resp, err := http.Get(vidlinl)
	if err != nil {
		return "", fmt.Errorf("не могу скачать файлик видосика: %s", err.Error())
	}
	defer resp.Body.Close()
	_, _ = io.Copy(out, resp.Body)
	return link + ".mp4", nil
}
