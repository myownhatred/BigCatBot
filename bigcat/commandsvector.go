package bigcat

import (
	"Guenhwyvar/entities"
	"Guenhwyvar/servitor"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func CmdVectorAddNewType(c tele.Context, serv *servitor.Servitor) (err error) {
	if c.Message().Payload == "" {
		return c.Send("пазязя дайте название викторинькиx!")
	}
	var typ entities.VectorType
	typ.Name = c.Message().Payload
	err = serv.FreeMawVectorTypeAdd(typ)
	if err != nil {
		return c.Send("произосло " + err.Error())
	}
	return c.Send("успесно добавили новый викториновский тип (тайп) 💕")
}

func CmdVectorGetTypes(c tele.Context, serv *servitor.Servitor) (err error) {
	report, err := serv.FreeMawVectorTypes()
	if err != nil {
		return c.Send("произосло " + err.Error())
	}
	message := ""
	for _, typ := range report {
		message += strconv.Itoa(typ.ID) + " - " + typ.Name
		if typ.Protected {
			message += "🛑"
		}
		message += "\n"
	}
	return c.Send(message)
}

func CmdVectorAddNew(c tele.Context, serv *servitor.Servitor) (err error) {
	if c.Message().Payload == "" {
		return c.Send("пазязя предоставьте 👐 данные для добавления нового воплосика!")
	}
	data := strings.Split(c.Message().Payload, ",")
	if len(data) < 4 {
		return c.Send("пазязя введите как минимум 4 параметра, разделённих запятими:\n- номир викториньки\n- ссылоцку на калтиноцку\n- воплосик\n- ответики чирис запитую, можно несколько")
	}
	for i, s := range data {
		data[i] = strings.TrimSpace(s)
	}

	var question entities.FreeVector
	question.TypeID, err = strconv.Atoi(data[0])
	if err != nil {
		return c.Send("осибка с номером викториньки - " + err.Error() + "\nвы точно ввели целое цисло?")
	}
	report, err := serv.FreeMawVectorTypes()
	if err != nil {
		return c.Send("осибка пли запроси номеров викторин - " + err.Error())
	}
	if !checkVectorTypesHaveNumber(report, question.TypeID) {
		return c.Send("осибка с номером викториньки - вы ввели несусествуюсий номер викторинки")
	}
	question.PicLink = data[1]
	if !strings.Contains(data[1], "http") {
		return c.Send("казеться вы ввели neco-ректную ссылотьку 😈")
	}
	question.Question = data[2]
	for i := 3; i < len(data); i++ {
		question.Answers = append(question.Answers, entities.VectorAnswer{ID: i, QuestionID: 6666, Answer: data[i]})
	}
	question.UserID = c.Sender().ID
	err = serv.FreeMawVectorAdd(question)
	if err != nil {
		return c.Send("осибка при внесении в базьку 😿 " + err.Error())
	}
	return c.Send("воплосик успесно внесёнь в базьку 💕")
}

func checkVectorTypesHaveNumber(report []entities.VectorType, number int) bool {
	for _, typ := range report {
		if typ.ID == number {
			return true
		}
	}
	return false
}
