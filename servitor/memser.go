package servitor

import "Guenhwyvar/bringer"

type MemserServ struct {
	bringer bringer.Memser
}

func NewMemserServ(bringer bringer.Memser) *MemserServ {
	return &MemserServ{bringer: bringer}
}

func (m *MemserServ) CreateGuiltyCatMeme(text string) (pathImg string, err error) {
	return m.bringer.CreateGuiltyCatMeme(text)
}
