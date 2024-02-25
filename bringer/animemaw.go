package bringer

import (
	"Guenhwyvar/entities"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type AnimeMawPostgres struct {
	db *gorm.DB
	r  *resty.Client
}

func NewAnimeMawPostgres(db *gorm.DB, r *resty.Client) *AnimeMawPostgres {
	return &AnimeMawPostgres{
		db: db,
		r:  r,
	}
}

func (m *AnimeMawPostgres) GetOpeningsFromDB() (openings []entities.AnimeOpening, err error) {
	// TODO: add error check
	_ = m.db.Find(&openings)

	return openings, nil
}

func (m *AnimeMawPostgres) GetOpeningsFromURL(url string) (openings []entities.AnimeOpening, err error) {
	response, err := m.r.R().Get(url)
	if err != nil {
		return openings, err
	}
	file := string(response.Body())
	lines := strings.Split(file, "\n")
	for _, line := range lines {
		var op entities.AnimeOpening
		pieces := strings.Split(line, ",")
		if len(pieces) != 2 {
			fmt.Println(line)
			continue
		}
		op.Description = pieces[0]
		op.Link = pieces[1]
		openings = append(openings, op)
	}
	return openings, nil
}

func (m *AnimeMawPostgres) PutOpeningsToDB(openings []entities.AnimeOpening) (affRows int64, err error) {
	_ = m.db.Exec("TRUNCATE TABLE anime_openings")
	result := m.db.Create(&openings)
	return result.RowsAffected, result.Error
}

func (m *AnimeMawPostgres) GetRandomOpening(typ string) (opening entities.AnimeOpening, err error) {
	switch typ {
	case "anime":
		opening, err = GetRandomOpeningByTyp("anime_openings", m.db)
	case "grob":
		opening, err = GetRandomOpeningByTyp("grobs", m.db)
	}
	if err != nil {
		return opening, err
	}
	return opening, nil
}

func GetRandomOpeningByTyp(typ string, db *gorm.DB) (opening entities.AnimeOpening, err error) {
	err = db.Raw(fmt.Sprintf("SELECT id, description, link FROM %s ORDER BY random() LIMIT 1", typ)).Scan(&opening).Error
	if err != nil {
		return opening, err
	}
	return opening, nil
}
