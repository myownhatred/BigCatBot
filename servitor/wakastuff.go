package servitor

import (
	"Guenhwyvar/bringer"
	"log"
)

type WakaStuffServ struct {
	bringer bringer.WakaStuff
}

func NewWakaStuffServ(bringer bringer.WakaStuff) *WakaStuffServ {
	return &WakaStuffServ{bringer: bringer}
}

func (w *WakaStuffServ) GetWakaStuff() string {
	text, err := w.bringer.GetDailyWaka()
	if err != nil {
		log.Printf("Error getting daily waka stuff: %v", err)
		text = "Проблемки с узнаванием прогресса однако"
		return text
	}
	return text
}
