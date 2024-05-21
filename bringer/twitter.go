package bringer

import (
	"context"
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

func (s *TwitterScrapper) TwitterGetHourlyPicture(acc string) (filePath string, err error) {
	count := 0
	for tweet := range s.s.GetTweets(context.Background(), acc, 2) {
		if tweet.Error != nil {
			return "", tweet.Error
		}
		if count > 0 {
			return tweet.Photos[0].URL, nil
		}
		// skipping first tweet on account just in case it is pinned
		// I don't know why I just don't pick second tweet
		// TODO check how it works and make it better
		count++
	}
	return "found nothing", nil
}
