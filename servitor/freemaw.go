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

func (m *FreeMawServ) FreeMawVectorTypeAdd(qtype entities.VectorType) (err error) {
	return m.bringer.FreeMawVectorTypeAdd(qtype)
}

func (m *FreeMawServ) FreeMawVectorTypeByID(ID int) (qtype entities.VectorType, err error) {
	return m.bringer.FreeMawVectorTypeByID(ID)
}
func (m *FreeMawServ) FreeMawVectorTypes() (report []entities.VectorType, err error) {
	return m.bringer.FreeMawVectorTypes()
}
func (m *FreeMawServ) FreeMawVectorAdd(vec entities.FreeVector) (err error) {
	return m.bringer.FreeMawVectorAdd(vec)
}
func (m *FreeMawServ) FreeMawVectorGetRandomByType(typ int) (vec entities.FreeVector, err error) {
	return m.bringer.FreeMawVectorGetRandomByType(typ)
}
