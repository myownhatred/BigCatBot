package main

import (
	"Guenhwyvar/bigcat"
	"Guenhwyvar/bringer"
	"Guenhwyvar/servitor"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"log/slog"

	"github.com/go-resty/resty/v2"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	tele "gopkg.in/telebot.v3"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {

	opts := PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := NewPrettyHandler(os.Stdout, opts)
	logger := slog.New(handler)
	//logger.Debug(
	//	"executing database query",
	//	slog.String("query", "SELECT * FROM users"),
	//)
	logger.Info("aux init started", slog.String("version", "0.1.3"))

	scrap := twitterscraper.New()
	scrap.WithDelay(5)

	f, _ := os.Open("cookies.json")
	var cookies []*http.Cookie
	json.NewDecoder(f).Decode(&cookies)
	scrap.SetCookies(cookies)
	if scrap.IsLoggedIn() {
		fmt.Println("Twitter is Ok")
	}

	configFile := pflag.String("config", "", "Path to config file")
	pflag.Parse()
	dodgeViper := viper.New()
	dodgeViper.BindPFlags(pflag.CommandLine)

	if *configFile != "" {
		dodgeViper.SetConfigFile(*configFile)
		dodgeViper.SetConfigType("yaml")
		if err := dodgeViper.ReadInConfig(); err != nil {
			fmt.Printf("Failed to read config file: %v\n", err)
			os.Exit(1)
		}
	}

	// db init
	psqlString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dodgeViper.GetString("postgreshost"), dodgeViper.GetString("postgresport"),
		dodgeViper.GetString("postgresuser"), dodgeViper.GetString("postgrespass"),
		dodgeViper.GetString("postgresdbname"))

	dbPostgres, err := sql.Open("postgres", psqlString)
	if err != nil {
		// probably could be just log message - we can live
		// without DB
		fmt.Printf("Failed to connect to DB: %v\n", err)
		os.Exit(1)
	}
	// close it, just in case
	defer dbPostgres.Close()

	bringa := bringer.NewBringer(resty.New(), scrap, dodgeViper, dbPostgres, logger)
	serva := servitor.NewServitor(bringa, logger)

	value2 := dodgeViper.GetString("telegramtoken")
	pref := tele.Settings{
		Token: value2,
		Poller: &tele.LongPoller{
			Timeout: 10 * time.Second,
			AllowedUpdates: []string{
				"message",
				"edited_message",
				"callback_query",
			},
		},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		fmt.Println("error creating new bot")
		log.Fatal(err)
		return
	}

	bigcat := bigcat.New(bot, serva, "start", logger)
	bigcat.Start()

	log.Print("service started")

}
