package config

import (
	"errors"
	"fmt"
)

type AppConfig struct {
	Meta         MetaConfig         `mapstructure:"meta"`
	ChatsAndPeps ChatsAndPepsConfig `mapstructure:"chatsandpeps"`
	Twitter      TwitterConfig      `mapstructure:"twitter"`
	Telegram     TelegramConfig     `mapstructure:"telegram"`
	Misc         MiscConifg         `mapstructure:"misc"`
	API          APIConfig          `mapstructure:"api"`
	Postgres     PostgresConfig     `mapstructure:"postgres"`
	Grok         GrokConfig         `mapstructure:"grok"`
}

type MetaConfig struct {
	Version  string `mapstructure:"version"`
	LogLevel string `mapstructure:"loglevel"`
}

type ChatsAndPepsConfig struct {
	MotherShip int64 `mapstructure:"mothership"`
	JonChat    int64 `mapstructure:"jonchat"`
	MisterX    int64 `mapstructure:"misterx"` // TODO remove after user policing implemented
	RapGodX    int64 `mapstructure:"rapgodx"`
}

type TwitterConfig struct {
	CookieFile string `mapstructure:"twittercookie"` // TODO make it filepatuh
	Delay      int    `mapstructure:"delay"`
}

type TelegramConfig struct {
	TelegramBotToken string `mapstructure:"telegramtoken"`
}

type MiscConifg struct {
	JokePath string `mapstructure:"joke"` // TODO make it filepath
}

type APIConfig struct {
	WakaTimeAPIToken string `mapstructure:"wakakey"`
	OpenWeatherToken string `mapstructure:"openweathertoken"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"postgreshost"`
	Port     string `mapstructure:"postgresport"`
	User     string `mapstructure:"postgresuser"`
	Password string `mapstructure:"postgrespass"`
	DBName   string `mapstructure:"postgresdbname"`
	SSLMode  string `mapstructure:"postgressslmode"`
}

type GrokConfig struct {
	Grokeys []string `mapstructure:"grokeys"`
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode)
}

func (m MetaConfig) Validate() error {
	vaidLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !vaidLevels[m.LogLevel] {
		return fmt.Errorf("invalid log_level '%s', must be one of: debug, info, warn, error", m.LogLevel)
	}

	return nil
}

func (c ChatsAndPepsConfig) Validate() error {
	var errs []error
	if c.MotherShip == 0 {
		errs = append(errs, fmt.Errorf("Mothership chat ID is missing"))
	}
	if c.MisterX == 0 {
		errs = append(errs, fmt.Errorf("MisterX ID is missing"))
	}

	return errors.Join(errs...)
}

func (t TelegramConfig) Validate() error {
	if t.TelegramBotToken == "" {
		return fmt.Errorf("No telegram bot token provided")
	}

	return nil
}

func (a APIConfig) Validate() error {
	var errs []error
	if a.OpenWeatherToken == "" {
		errs = append(errs, fmt.Errorf("No openweather token provided"))
	}
	if a.WakaTimeAPIToken == "" {
		errs = append(errs, fmt.Errorf("No wakatime token provided"))
	}

	return errors.Join(errs...)
}

func (p PostgresConfig) Validate() error {
	var errs []error
	if p.Host == "" {
		errs = append(errs, fmt.Errorf("No postgress host provided"))
	}
	if p.Port == "" {
		errs = append(errs, fmt.Errorf("No postgress port provided"))
	}
	if p.DBName == "" {
		errs = append(errs, fmt.Errorf("No postgress dbname provided"))
	}
	if p.User == "" {
		errs = append(errs, fmt.Errorf("No postgress user provided"))
	}
	if p.Password == "" {
		errs = append(errs, fmt.Errorf("No postgress password provided"))
	}
	if p.SSLMode == "" {
		errs = append(errs, fmt.Errorf("No postgress sslmode provided"))
	}

	return errors.Join(errs...)
}

func (g GrokConfig) Validate() error {
	if len(g.Grokeys) == 0 {
		return fmt.Errorf("No grok API-keys provided")
	}

	return nil
}

func (c *AppConfig) Validate() error {
	var errs []error
	if err := c.Meta.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("Meta validation error %w", err))
	}
	if err := c.ChatsAndPeps.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("Chat & Peps validation error %w", err))
	}
	if err := c.Telegram.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("Telegram validation error %w", err))
	}
	if err := c.API.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("API validation error %w", err))
	}
	if err := c.Postgres.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("Postgres validation error %w", err))
	}
	if err := c.Grok.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("Grok validation error %w", err))
	}

	return errors.Join(errs...)
}
