package servitor

import (
	"Guenhwyvar/bringer"
	"Guenhwyvar/entities"
)

type TimeWithOutServ struct {
	bringer bringer.TimeWithOut
}

func NewTimeWithOutServ(bringer bringer.TimeWithOut) *TimeWithOutServ {
	return &TimeWithOutServ{
		bringer: bringer,
	}
}

func (t *TimeWithOutServ) GetTimeWithOutList(chatID int64) (list []entities.TimeWithOut, err error) {
	list, err = t.bringer.GetTimeWithOutList(chatID)
	if err != nil {
		return list, err
	}
	return list, nil
}

func (t *TimeWithOutServ) GetTimeWithOutTimerByID(id int) (event entities.TimeWithOut, err error) {
	return t.bringer.GetTimeWithOutTimerByID(id)
}

func (t *TimeWithOutServ) ResetTimer(id int) (err error) {
	return t.bringer.ResetTimer(id)
}

func (t *TimeWithOutServ) AddNewTimer(name string, chatID int64) (err error) {
	return t.bringer.AddNewTimer(name, chatID)
}
