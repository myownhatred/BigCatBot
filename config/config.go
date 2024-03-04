package config

type AppConfig struct {
	TelegramBotToken string
	WakaTimeAPIToken string
	MotherShip       int64
	JokePath         string // TODO make it filepath
	TwitterCookie    string // TODO make it filepath
}
