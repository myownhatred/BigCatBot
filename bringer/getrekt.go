package bringer

import (
	"Guenhwyvar/entities"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly"
	"github.com/spf13/viper"
)

type GetRektResty struct {
	r *resty.Client
	v *viper.Viper
}

type FileState struct {
	State string `json:"state"`
}

func NewGetRect(r *resty.Client, v *viper.Viper) *GetRektResty {
	return &GetRektResty{
		r: r,
		v: v,
	}
}

func (rekt *GetRektResty) GetWeatherDayForecast(place string) (report string, err error) {
	// TODO: make base URL great again
	OpenWeatherBaseURL := "https://api.openweathermap.org"
	requestLine := fmt.Sprintf("%s/data/2.5/forecast?q=%s&exclude=minutely&appid=%s&units=metric&lang=ru", OpenWeatherBaseURL, place, rekt.v.GetString("openweathertoken"))
	response, err := rekt.r.R().Get(requestLine)
	if err != nil {
		return "", fmt.Errorf("error getting weather forecast: %s", err)
	}
	// TODO: add to debug logging level

	var forecast Forecast
	if err = json.Unmarshal(response.Body(), &forecast); err != nil {
		return "", fmt.Errorf("error unmarshalling forecast to json: %s", err)
	}
	// TODO: add timezone/timeshift to function signature
	// timezone shift to GMT+7
	currentTime := time.Now()
	newTime := currentTime.Add(1 * time.Hour)
	day := newTime.Day()
	var windf, hurrf, coldf, hotf, rainf bool
	report += "Прогноз:\n"
	for _, w := range forecast.List {
		// getting datetime of forecast record
		tm := time.Unix(w.Dt, 0)
		dayf := tm.Day()
		if dayf == day {
			h, _, _ := tm.Clock()
			if h > 6 && h < 22 {
				if w.WsM.Temp < 10 {
					coldf = true
				}
				if w.WsM.Temp > 28 {
					hotf = true
				}
				if w.Wind.Speed > 7 {
					windf = true
				}
				if w.Wind.Speed > 25 {
					hurrf = true
				}
				if w.Pop > 0 {
					rainf = true
				}
			}
			switch h {
			case 0, 1:
				report += "Утро(ебать)☕\n"
				report += hourReport(w.WsM.Temp, w.Wind.Speed, w.Ws[0].Description)
			case 5, 6:
				report += "Обед🍴🍲\n"
				report += hourReport(w.WsM.Temp, w.Wind.Speed, w.Ws[0].Description)
			case 11, 12:
				report += "Вечер 𓀐𓂸🤱🏻🤰\n"
				report += hourReport(w.WsM.Temp, w.Wind.Speed, w.Ws[0].Description)
			}
		}
	}
	if coldf {
		report += "прохладно нужен кофтан\n"
	}
	if hotf {
		report += "жара, тёлки могут скинуть шмотки, не забудьте презики\n"
	}
	if windf && !hurrf {
		report += "ветренно, если вы карлик - может унести в казахстан\n"
	}
	if hurrf {
		report += "ураган, может унести в талнах даже если не карлик\n"
	}
	if rainf {
		report += "возможен додж, возьмите зонтикс\n"
	}

	return report, nil
}

func (rekt *GetRektResty) GetCurrentWeather(place string) (report string, err error) {
	// TODO: make base URL greath again
	OpenWeatherBaseURL := "https://api.openweathermap.org"
	requestLine := fmt.Sprintf("%s/data/2.5/weather?q=%s&appid=%s&units=metric&lang=ru", OpenWeatherBaseURL, place, rekt.v.GetString("openweathertoken"))
	response, err := rekt.r.R().Get(requestLine)
	if err != nil {
		return "", err
	}
	var w WeatherResponse
	err = json.Unmarshal(response.Body(), &w)
	if err != nil {
		return "", err
	}

	report = fmt.Sprintf("Погода для %s\n", w.Name)
	report += fmt.Sprintf("Дипазон паходы: %2.f - %2.f\n", w.WsM.TempMin, w.WsM.TempMax)
	report += fmt.Sprintf("Сёдня жарит %2.f градусиф\nАщущения как %2.f\nМокрость %d\nДовление аж %d\n", w.WsM.Temp, w.WsM.FeelsLike, w.WsM.Humidity, w.WsM.Pressure)
	report += fmt.Sprintf("Ветрище %3.fм/с дует с порывами до %3.fм/с по азимуту %d\n", w.Wind.Speed, w.Wind.Gust, w.Wind.Degree)
	report += fmt.Sprintf("За последний час нападал дощь %3.2f мм\n", w.Rain.Hour)
	report += fmt.Sprintf("А снега навалило %3.2f мм\n", w.Snow.Hour)
	report += fmt.Sprintf("Видимость %d Облаковость %d\nВ целом:", w.Vis, w.Cl.All)
	for _, wes := range w.Ws {
		report += fmt.Sprintf("%s ", wes.Description)
	}
	report += "\n"

	return report, nil
}

func (rekt *GetRektResty) SendGenerationReq(modelID int, prompt string) (err error) {
	rekt.r.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	generatorURL := "https://localhost:6009/generate"
	_, err = rekt.r.R().
		SetBody(GenerationRequest{ID: modelID, Prompt: prompt}).
		Post(generatorURL)
	if err != nil {
		return err
	}
	return nil
}

func (rekt *GetRektResty) GetGenerationStatus() (status string, err error) {
	rekt.r.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	statusURL := "https://localhost:6009/status"
	resp, err := rekt.r.R().Get(statusURL)
	if err != nil {
		return "", err
	}
	var s FileState
	err = json.Unmarshal(resp.Body(), &s)
	if err != nil {
		return "", err
	}
	return s.State, nil
}

func (rekt *GetRektResty) GetGeneratorStatus() (singa entities.Signa, err error) {
	rekt.r.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	statusURL := "https://localhost:6009/cnc"
	resp, err := rekt.r.R().Get(statusURL)
	if err != nil {
		return entities.Signa{}, err
	}
	var s entities.Signa
	err = json.Unmarshal(resp.Body(), &s)
	if err != nil {
		return entities.Signa{}, err
	}
	return s, nil
}

func hourReport(temp, wind float32, desc string) string {
	return fmt.Sprintf("🌡️:%d°С 💨: %d м/с в целом:%s\n", int(temp), int(wind), desc)
}

func (rekt *GetRektResty) GetRandomMTG() (url string, err error) {
	site := "https://scryfall.com/random"
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)
	c.OnHTML("[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		if strings.Contains(link, "https://cards.scryfall.io") {
			url = link
		}
	})
	c.Visit(site)
	return url, nil
}

func (rekt *GetRektResty) GetFreeSteamGames() (string, error) {
	res, err := rekt.r.R().
		Get("https://raw.githubusercontent.com/InJeCTrL/NeedFree/master/free_goods_detail.json")
	if err != nil {
		return "", err
	}

	var ret FreeGames
	if err = json.Unmarshal(res.Body(), &ret); err != nil {
		return "", err
	}
	otvet := ""
	otvet += fmt.Sprintf("Сегодня навалили халявы %d титулов\n", ret.TotalCount)
	for _, gre := range ret.FreeList {
		otvet += fmt.Sprintf("%s - %s\n", gre[0], gre[1])
	}
	otvet += fmt.Sprintf("последний раз проверялось - %s\n", ret.UpdateTime)
	return otvet, nil
}

// structs
type Weather struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type WeatherMain struct {
	Temp      float32 `json:"temp"`
	FeelsLike float32 `json:"feels_like"`
	TempMin   float32 `json:"temp_min"`
	TempMax   float32 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

type Wind struct {
	Speed  float32 `json:"speed"`
	Degree int     `json:"deg"`
	Gust   float32 `json:"gust"`
}

type Clouds struct {
	All int `json:"all"`
}

type Rain struct {
	Hour       float32 `json:"1h"`
	ThreeHours float32 `json:"3h"`
}

type Snow struct {
	Hour       float32 `json:"1h"`
	ThreeHours float32 `json:"3h"`
}

type MiscInfo struct {
	Sunrise int `json:"sunrise"`
	Sunset  int `json:"sunset"`
}

type WeatherResponse struct {
	Dt   int64       `json:"dt"`
	Ws   []Weather   `json:"weather"`
	WsM  WeatherMain `json:"main"`
	Vis  int         `json:"visibility"`
	Cl   Clouds      `json:"clouds"`
	Pop  float32     `json:"pop"`
	Wind Wind        `json:"wind"`
	Rain Rain        `json:"rain"`
	Snow Snow        `json:"snow"`
	MI   MiscInfo    `json:"sys"`
	Name string      `json:"name"`
}

type Forecast struct {
	Cod     string            `json:"cod"`
	Message int               `json:"message"`
	Cnt     int               `json:"cnt"`
	List    []WeatherResponse `json:"list"`
}

type FreeGames struct {
	TotalCount int        `json:"total_count"`
	FreeList   [][]string `json:"free_list"`
	UpdateTime string     `json:"update_time"`
}

type GenerationRequest struct {
	ID     int    `json:"id"`
	Prompt string `json:"prompt"`
}
