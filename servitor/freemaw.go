package servitor

import (
	"Guenhwyvar/bringer"
	"Guenhwyvar/entities"
)

type FreeMawServ struct {
	bringer bringer.FreeMaw
}

func NewFreeMawServ(bringer bringer.FreeMaw) *FreeMawServ {
	return &FreeMawServ{
		bringer: bringer,
	}
}

func (m *FreeMawServ) GetFreeMaw(typ string) (maw entities.FreeMaw, err error) {
	return m.bringer.GetRandomMawFromDB(typ)
}

func (m *FreeMawServ) PutFreeMaw(maw entities.FreeMaw) (err error) {
	return m.bringer.PutFreeMawToDB(maw)
}

func (m *FreeMawServ) GetFreeMawReport() (report string, err error) {
	return "скоро тут будет красивый ачот", nil
}
