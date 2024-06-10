package servitor

import (
	"Guenhwyvar/bringer"
)

type GetRektServ struct {
	bringer bringer.GetRekt
}

func NewGetRectServ(bringer bringer.GetRekt) *GetRektServ {
	return &GetRektServ{bringer: bringer}
}

func (rekt *GetRektServ) GetWeatherDayForecast(place string) (report string, err error) {
	report, err = rekt.bringer.GetWeatherDayForecast(place)
	if err != nil {
		return "", err
	}
	return report, nil
}

func (rekt *GetRektServ) GetCurrentWeather(place string) (report string, err error) {
	report, err = rekt.bringer.GetCurrentWeather(place)
	if err != nil {
		return "", err
	}
	return report, nil
}

func (rekt *GetRektServ) GetFreeSteamGames() (report string, err error) {
	report, err = rekt.bringer.GetFreeSteamGames()
	if err != nil {
		return "", err
	}
	return report, nil
}
