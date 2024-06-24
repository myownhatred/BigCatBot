package bringer

import (
	"Guenhwyvar/config"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type ComfigViper struct {
	v *viper.Viper
}

func NewComfigerViper(v *viper.Viper) *ComfigViper {
	return &ComfigViper{
		v: v,
	}
}

func (c *ComfigViper) GetAppComfigFromViper() (*config.AppConfig, error) {
	var comfig *config.AppConfig = new(config.AppConfig)
	if c.v.IsSet("twittercookie") {
		comfig.TwitterCookie = c.v.GetString("twittercookie")
	}
	if c.v.IsSet("jokepath") {
		comfig.JokePath = c.v.GetString("jokepath")
	} else {
		// TODO make it just log message
		log.Fatalf("can't find dad's joke snatchel")
	}
	if c.v.IsSet("misterx") {
		comfig.MisterX = c.v.GetInt64("misterx")
	}
	if c.v.IsSet("mothership") {
		comfig.MotherShip = c.v.GetInt64("mothership")
		return comfig, nil
	} else {
		err := fmt.Errorf("no mothership config value set")
		return comfig, err
	}

}
