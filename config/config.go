package config

type AppConfig struct {
	TelegramBotToken string
	WakaTimeAPIToken string
	OpenWeatherToken string
	MotherShip       int64
	JokePath         string // TODO make it filepath
	TwitterCookie    string // TODO make it filepath
}
