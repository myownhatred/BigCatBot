package servitor

import (
	"Guenhwyvar/bringer"
	"Guenhwyvar/config"
)

type ComfigerServ struct {
	bringer bringer.Comfiger
}

func NewComfigerServ(bringer bringer.Comfiger) *ComfigerServ {
	return &ComfigerServ{
		bringer: bringer,
	}
}

func (c *ComfigerServ) GetAppComfig() (comfig *config.AppConfig, err error) {
	return c.bringer.GetAppComfigFromViper()
}
