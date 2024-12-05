package bigcat

import (
	"Guenhwyvar/entities"
	freevector "Guenhwyvar/lib/vector"
	"Guenhwyvar/servitor"
	"log/slog"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
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

func CmdVectorGame(c tele.Context, s *servitor.Servitor, brain *BigBrain) (err error) {
	// one game at a time
	if brain.ChatFlags[c.Chat().ID].VectorGame {
		return c.Send("викторина уже запущена❗️")
	}
	if c.Message().Payload == "" {
		return c.Send("пазязя предоставьте 👐 номер викторины!")
	}
	num, err := strconv.Atoi(c.Message().Payload)
	if err != nil {
		return c.Send("осибка с номером викториньки\nвы точно ввели целое цисло?")
	}
	c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, strconv.Itoa(num))
	s.Logger.Info("vector game",
		slog.Int("mumber", num))
	gF := brain.ChatFlags[c.Chat().ID]
	gF.VectorGame = true
	brain.ChatFlags[c.Chat().ID] = gF
	vc := freevector.NewVectorCore()
	for i := 0; i < 4; i++ {
		question, err := s.FreeMawVectorGetRandomByType(num)
		answerString := ""
		for ix, a := range question.Answers {
			if ix == len(question.Answers)-1 {
				answerString += a.Answer
			} else {
				answerString += a.Answer + " / "
			}
		}

		if err != nil {
			return c.Send("осибка с номером при получении вопроса викторины:\n" + err.Error())
		}
		vc.CurrentQuestion = question
		brain.VectorGame[c.Chat().ID] = *vc
		s.Logger.Info("vector game calling question",
			slog.Int("mumber ", i))
		one, two, three := createHelpLines(question, s.Logger)
		m := &tele.Photo{
			File:    tele.FromURL(question.PicLink),
			Caption: question.Question,
		}
		c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, m)
		firstHelp := 14 * time.Second
		secondHelp := 15 * time.Second
		thirdHelp := 16 * time.Second
		endGame := 17 * time.Second
		sex := ""
	timers:
		select {
		case <-time.After(firstHelp):
			s.Logger.Info("vector game loop question case of 15 secs",
				slog.Int("mumber ", i))
			c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, "подсказька - "+one)
			firstHelp += time.Hour
			goto timers
		case <-time.After(secondHelp):
			s.Logger.Info("vector game loop question case of 30 secs",
				slog.Int("mumber ", i))
			c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, "подсказька - "+two)
			secondHelp += time.Hour
			goto timers
		case <-time.After(thirdHelp):
			s.Logger.Info("vector game loop question case of 45 secs",
				slog.Int("mumber ", i))
			c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, "подсказька - "+three)
			thirdHelp += time.Hour
			goto timers
		case <-time.After(endGame):
			s.Logger.Info("vector game loop question case of 60 secs",
				slog.Int("mumber ", i))
			c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, "какие то вайбы кожанного позора, ответ - "+answerString)
			time.Sleep(4 * time.Second)
			continue
		case sex = <-brain.VectorChan:
			s.Logger.Info("vector game loop question case of right answer",
				slog.Int("mumber ", i))
			c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, "правильный атвет,"+sex)
			c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, answerString)
			time.Sleep(4 * time.Second)
			continue
		}
		//c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, one+" - "+two+" - "+three)
	}
	gF.VectorGame = false
	brain.ChatFlags[c.Chat().ID] = gF
	return c.Send("конец раунда, биатч")
}

func checkVectorTypesHaveNumber(report []entities.VectorType, number int) bool {
	for _, typ := range report {
		if typ.ID == number {
			return true
		}
	}
	return false
}

func createHelpLines(q entities.FreeVector, logger *slog.Logger) (first, second, third string) {
	longest := ""
	for _, a := range q.Answers {
		if len(a.Answer) > len(longest) {
			longest = a.Answer
		}
	}
	longest = q.Answers[rand.Intn(len(q.Answers))].Answer
	logger.Info("vector game helping answer is",
		slog.String("answer ", longest))
	// short lines exceptions
	if len(longest) <= 3 {
		return "***", "***", "***"
	}
	revealFirst := int(math.Round(float64(len(longest) * 10 / 100)))
	logger.Info("vector game estimates ",
		slog.Int("10% ", revealFirst))
	revealSecond := int(math.Round(float64(len(longest) * 20 / 100)))
	logger.Info("vector game estimates ",
		slog.Int("20% ", revealSecond))
	revealThird := int(math.Round(float64(len(longest) * 30 / 100)))
	logger.Info("vector game estimates ",
		slog.Int("40% ", revealThird))
	// Convert to rune slice to handle multi-byte characters properly
	if len(longest) <= 4 {
		revealFirst = 1
		revealSecond = 1
		revealThird = 2
	}
	if len(longest) <= 6 {
		revealFirst = 1
		revealSecond = 2
		revealThird = 3
	}
	runes := []rune(longest)
	revealed := make([]rune, len(runes))
	logger.Info("vector game masking words with * and spaces")
	for i := range revealed {
		if runes[i] == ' ' {
			revealed[i] = ' '
		} else {
			revealed[i] = '*'
		}
	}
	logger.Info("vector game starting feeling loop")

	for i := 0; i < revealThird; {
		j := rand.Intn(len(revealed))
		if revealed[j] == '*' {
			revealed[j] = runes[j]
			i++
			logger.Info("vector game help creation",
				slog.String("help line ", string(revealed)),
				slog.Int("number of reveals ", i))
			if i == revealFirst {
				first = string(revealed)
				second = string(revealed)
			}
			if i == revealSecond {
				second = string(revealed)
			}
		}
	}
	third = string(revealed)
	return first, second, third
}

func checkAnswer(q entities.FreeVector, variant string) bool {
	for _, a := range q.Answers {
		if a.Answer == variant {
			return true
		}
	}
	return false
}
