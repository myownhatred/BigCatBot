package servitor

import (
	"Guenhwyvar/bringer"
	"strings"
)

type TwitterServ struct {
	bringer bringer.Twitter
}

func NewTwitterServ(bringer bringer.Twitter) *TwitterServ {
	return &TwitterServ{bringer: bringer}
}

func (t *TwitterServ) TwitterGetVideo(link string) (pathImg string, err error) {
	twitString := strings.Split(link, "/status/")
	id := twitString[1][:19]
	return t.bringer.TwitterGetVideo(id)
}
