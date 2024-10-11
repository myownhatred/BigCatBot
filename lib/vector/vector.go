package freevector

import (
	"Guenhwyvar/entities"
	"strings"
)

type VectorCore struct {
	CurrentQuestion entities.FreeVector
}

func NewVectorCore() (vc *VectorCore) {
	return &VectorCore{}
}

func (vc *VectorCore) CheckAnswer(t string) (result bool) {
	for _, a := range vc.CurrentQuestion.Answers {
		if strings.ToLower(a.Answer) == t {
			return true
		}
	}
	return false
}
