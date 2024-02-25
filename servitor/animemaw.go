package servitor

import (
	"Guenhwyvar/bringer"
	"Guenhwyvar/entities"
	"encoding/csv"
	"fmt"
	"os"
)

type AnimeMawServ struct {
	bringer bringer.AnimeMaw
}

func NewAnimeMawServ(bringer bringer.AnimeMaw) *AnimeMawServ {
	return &AnimeMawServ{
		bringer: bringer,
	}
}

func (m *AnimeMawServ) GetOpeningsFilePath() (filePath string) {
	openings, _ := m.bringer.GetOpeningsFromDB()

	file, err := os.Create("openings.csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, op := range openings {
		writer.Write(append([]string{}, op.Description, op.Link))
	}

	return "openings.csv"
}

func (m *AnimeMawServ) UploadOpenings(filePath string) (report string, err error) {
	file, err := os.Open("openingsfile.csv")
	if err != nil {
		return "", err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 2
	lines, _ := reader.ReadAll()
	report += fmt.Sprintf("Всего линий в файле на зяливку: %d\n", lines)

	uploadStuff := make([]entities.AnimeOpening, 0)

	for _, line := range lines {
		var open entities.AnimeOpening
		open.Description = line[0]
		open.Link = line[1]
		uploadStuff = append(uploadStuff, open)
	}

	report += fmt.Sprintf("Из файла насобиралось %d опенингов\n", len(uploadStuff))
	rows, err := m.bringer.PutOpeningsToDB(uploadStuff)
	if err != nil {
		return report, err
	}
	report += fmt.Sprintf("в бязу попяло %d опенингов\n", rows)

	return report, nil

}

func (m *AnimeMawServ) UploadOpeningsByURL(url string) (report string, err error) {
	openings, _ := m.bringer.GetOpeningsFromURL(url)
	rowsCount, err := m.bringer.PutOpeningsToDB(openings)
	report = fmt.Sprintf("Всего залили стаффчика %d штук", rowsCount)
	return report, err
}

func (m *AnimeMawServ) GetRandomOpening(typ string) (opening entities.AnimeOpening, err error) {
	opening, err = m.bringer.GetRandomOpening(typ)
	if err != nil {
		return opening, err
	}
	return opening, nil
}
