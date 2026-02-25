package main

import (
	"Guenhwyvar/bigcat"
	"Guenhwyvar/bringer"
	"Guenhwyvar/config"
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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	twitterscraper "github.com/tbdsux/twitter-scraper"
	tele "gopkg.in/telebot.v4"
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
	logger.Info("aux init started")

	scrap := twitterscraper.New()
	scrap.WithDelay(5)

	configFile := pflag.String("config", "", "Path to config file")
	pflag.Parse()
	dodgeViper := viper.New()
	dodgeViper.BindPFlags(pflag.CommandLine)
	comfig := &config.AppConfig{}

	if *configFile == "" {
		fmt.Println("config file is required: --config <path>")
		os.Exit(1)
	}

	if *configFile != "" {
		dodgeViper.SetConfigFile(*configFile)
		dodgeViper.SetConfigType("yaml")
		if err := dodgeViper.ReadInConfig(); err != nil {
			fmt.Printf("Failed to read config file: %v\n", err)
			os.Exit(1)
		}
		dodgeViper.Unmarshal(comfig)
		valerr := comfig.Validate()
		if valerr != nil {
			panic(valerr)
		}
		fmt.Printf("config: %+v\n", comfig)
	}

	// db init
	psqlString := comfig.Postgres.DSN()

	dbPostgres, err := sql.Open("postgres", psqlString)
	if err != nil {
		// probably could be just log message - we can live
		// without DB
		fmt.Printf("Failed to connect to DB: %v\n", err)
		os.Exit(1)
	}
	// close it, just in case
	defer dbPostgres.Close()

	f, _ := os.Open(comfig.Twitter.CookieFile)
	var cookies []*http.Cookie
	json.NewDecoder(f).Decode(&cookies)
	scrap.SetCookies(cookies)
	if scrap.IsLoggedIn() {
		fmt.Println("Twitter is Ok")
	}

	bringa := bringer.NewBringer(resty.New(), scrap, comfig, dbPostgres, logger)
	serva := servitor.NewServitor(bringa, logger, comfig)

	pref := tele.Settings{
		Token: comfig.Telegram.TelegramBotToken,
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

	bigcat := bigcat.New(bot, comfig, serva, "start", logger)
	bigcat.Start()

	logger.Info("service started")

}
