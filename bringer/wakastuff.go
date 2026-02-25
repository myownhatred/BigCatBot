package bringer

import (
	"Guenhwyvar/config"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type WakaStuffResty struct {
	r *resty.Client
	c *config.AppConfig
}

func NewWakaStuff(r *resty.Client, c *config.AppConfig) *WakaStuffResty {
	return &WakaStuffResty{
		r: r,
		c: c,
	}
}

func (w *WakaStuffResty) GetDailyWaka() (string, error) {
	dateFormat := "2006-01-02"
	now := time.Now()
	tomorrow := now.Add(time.Hour * 24)
	nowString := now.Format(dateFormat)
	tomorrowString := tomorrow.Format(dateFormat)
	requestLine := fmt.Sprintf("https://wakatime.com/api/v1/users/current/summaries?start=%s&end=%s&api_key=%s", nowString, tomorrowString, w.c.API.WakaTimeAPIToken)

	response, err := w.r.R().Get(requestLine)
	if err != nil {
		log.Printf("wakatime api request failed: %s", err)
		return "", err
	}

	var data WakaStruct
	var message string

	json.Unmarshal(response.Body(), &data)

	timeString := data.CumulativeTotal.Text
	strings.Replace(timeString, "hr", "ч", -1)
	strings.Replace(timeString, "mins", "мин", -1)

	message = fmt.Sprintf("Даницка севодня поработаль над монолитами примерно %s", timeString)

	return message, nil

}

type WakaStruct struct {
	Data []struct {
		Languages []struct {
			Name         string  `json:"name"`
			TotalSeconds float64 `json:"total_seconds"`
			Digital      string  `json:"digital"`
			Decimal      string  `json:"decimal"`
			Text         string  `json:"text"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Seconds      int     `json:"seconds"`
			Percent      float64 `json:"percent"`
		} `json:"languages"`
		GrandTotal struct {
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			TotalSeconds float64 `json:"total_seconds"`
			Digital      string  `json:"digital"`
			Decimal      string  `json:"decimal"`
			Text         string  `json:"text"`
		} `json:"grand_total"`
		Editors []struct {
			Name         string  `json:"name"`
			TotalSeconds float64 `json:"total_seconds"`
			Digital      string  `json:"digital"`
			Decimal      string  `json:"decimal"`
			Text         string  `json:"text"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Seconds      int     `json:"seconds"`
			Percent      float64 `json:"percent"`
		} `json:"editors"`
		OperatingSystems []struct {
			Name         string  `json:"name"`
			TotalSeconds float64 `json:"total_seconds"`
			Digital      string  `json:"digital"`
			Decimal      string  `json:"decimal"`
			Text         string  `json:"text"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Seconds      int     `json:"seconds"`
			Percent      float64 `json:"percent"`
		} `json:"operating_systems"`
		Categories []struct {
			Name         string  `json:"name"`
			TotalSeconds float64 `json:"total_seconds"`
			Digital      string  `json:"digital"`
			Decimal      string  `json:"decimal"`
			Text         string  `json:"text"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Seconds      int     `json:"seconds"`
			Percent      float64 `json:"percent"`
		} `json:"categories"`
		Dependencies []struct {
			Name         string  `json:"name"`
			TotalSeconds float64 `json:"total_seconds"`
			Digital      string  `json:"digital"`
			Decimal      string  `json:"decimal"`
			Text         string  `json:"text"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Seconds      int     `json:"seconds"`
			Percent      float64 `json:"percent"`
		} `json:"dependencies"`
		Machines []struct {
			Name          string  `json:"name"`
			TotalSeconds  float64 `json:"total_seconds"`
			MachineNameID string  `json:"machine_name_id"`
			Digital       string  `json:"digital"`
			Decimal       string  `json:"decimal"`
			Text          string  `json:"text"`
			Hours         int     `json:"hours"`
			Minutes       int     `json:"minutes"`
			Seconds       int     `json:"seconds"`
			Percent       float64 `json:"percent"`
		} `json:"machines"`
		Projects []struct {
			Name         string  `json:"name"`
			TotalSeconds float64 `json:"total_seconds"`
			Color        any     `json:"color"`
			Digital      string  `json:"digital"`
			Decimal      string  `json:"decimal"`
			Text         string  `json:"text"`
			Hours        int     `json:"hours"`
			Minutes      int     `json:"minutes"`
			Seconds      int     `json:"seconds"`
			Percent      float64 `json:"percent"`
		} `json:"projects"`
		Range struct {
			Start    time.Time `json:"start"`
			End      time.Time `json:"end"`
			Date     string    `json:"date"`
			Text     string    `json:"text"`
			Timezone string    `json:"timezone"`
		} `json:"range"`
	} `json:"data"`
	Start           time.Time `json:"start"`
	End             time.Time `json:"end"`
	CumulativeTotal struct {
		Seconds float64 `json:"seconds"`
		Text    string  `json:"text"`
		Digital string  `json:"digital"`
		Decimal string  `json:"decimal"`
	} `json:"cumulative_total"`
	DailyAverage struct {
		Holidays                      int    `json:"holidays"`
		DaysMinusHolidays             int    `json:"days_minus_holidays"`
		DaysIncludingHolidays         int    `json:"days_including_holidays"`
		Seconds                       int    `json:"seconds"`
		SecondsIncludingOtherLanguage int    `json:"seconds_including_other_language"`
		Text                          string `json:"text"`
		TextIncludingOtherLanguage    string `json:"text_including_other_language"`
	} `json:"daily_average"`
}
