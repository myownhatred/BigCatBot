package freevector

import (
	"Guenhwyvar/entities"
	"strings"
)

type VectorChanS struct {
	Uid  int64
	Text string
}

type VectorCore struct {
	CurrentQuestion entities.FreeVector
	VectorChan      chan VectorChanS
}

func NewVectorCore() (vc *VectorCore) {
	return &VectorCore{
		CurrentQuestion: entities.FreeVector{},
		VectorChan:      make(chan VectorChanS),
	}
}

func (vc *VectorCore) CheckAnswer(t string) (result bool) {
	for _, a := range vc.CurrentQuestion.Answers {
		if strings.ToLower(a.Answer) == t {
			return true
		}
	}
	return false
}
