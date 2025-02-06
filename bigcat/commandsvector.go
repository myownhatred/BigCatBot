package bigcat

import (
	"Guenhwyvar/entities"
	freevector "Guenhwyvar/lib/vector"
	"Guenhwyvar/servitor"
	"fmt"
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
	data := strings.Split(c.Message().Payload, ";")
	if len(data) < 4 {
		return c.Send("пазязя введите как минимум 4 параметра, разделённих тоцкой с запятой (запятоцкой):\n- номир викториньки\n- ссылоцку на калтиноцку\n- воплосик\n- ответики чирис запитую, можно несколько")
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
	answersArray := strings.Split(data[3], ",")
	for i := 0; i < len(answersArray); i++ {
		answersArray[i] = strings.TrimSpace(answersArray[i])
		// QuestionID:6666 is placeholder and be changed at the moment of DB insert
		question.Answers = append(question.Answers, entities.VectorAnswer{ID: i, QuestionID: 6666, Answer: answersArray[i]})
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
	//c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, strconv.Itoa(num))
	s.Logger.Info("vector game",
		slog.Int("mumber", num))
	gF := brain.ChatFlags[c.Chat().ID]
	gF.VectorGame = true
	brain.ChatFlags[c.Chat().ID] = gF
	var rs = make(map[int64]int)
	vc := freevector.NewVectorCore()
	for i := 0; i < 10; i++ {
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
			Caption: fmt.Sprintf("Вопрос %d - %s", i+1, question.Question),
		}
		c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, m)
		firstHelp := false
		secondHelp := false
		thirdHelp := false

		vecscore := 3
		gameStart := time.Now()

		for {
			eslaped := time.Since(gameStart)
			select {
			case uidText := <-vc.VectorChan:
				// if we have 0 in chan - means /vectorstop has arrived
				if uidText.Uid == 0 {
					c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, "заканчиваем шпилу!")
					goto endshpil
				}

				if vc.CheckAnswer(strings.ToLower(uidText.Text)) {
					uname := brain.Users[uidText.Uid].Username
					s.Logger.Info("vector game loop question case of right answer",
						slog.Int("mumber ", i))
					c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, "правильный атвет, "+uname)
					rs[uidText.Uid] += vecscore
					s.FreeMaw.FreeMawVectorUpsertScore(uidText.Uid, num, vecscore)
					c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, answerString)
					time.Sleep(4 * time.Second)
					goto nextQuestion
				}

			// checking stuff every second
			case <-time.After(1 * time.Second):
				if eslaped >= 15*time.Second && !firstHelp {
					s.Logger.Info("vector game loop question case of 15 secs",
						slog.Int("mumber ", i))
					c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, "подсказька - "+one)
					firstHelp = true

				}
				if eslaped >= 30*time.Second && !secondHelp {
					s.Logger.Info("vector game loop question case of 30 secs",
						slog.Int("mumber ", i))
					c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, "подсказька - "+two)
					secondHelp = true
					vecscore--
				}
				if eslaped >= 45*time.Second && !thirdHelp {
					s.Logger.Info("vector game loop question case of 30 secs",
						slog.Int("mumber ", i))
					c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, "подсказька - "+three)
					thirdHelp = true
					vecscore--
				}
				if eslaped >= 60*time.Second {
					s.Logger.Info("vector game loop question case of 60 secs",
						slog.Int("mumber ", i))
					c.Bot().Send(&tele.Chat{ID: c.Chat().ID}, "какие то вайбы кожанного позора, ответ - "+answerString)
					time.Sleep(4 * time.Second)
					goto nextQuestion
				}

			}
		}
	nextQuestion:
		continue
	}
endshpil:
	gF.VectorGame = false
	brain.ChatFlags[c.Chat().ID] = gF
	res := ""
	for k, v := range rs {
		res += fmt.Sprintf("плеер %s набрал %d очкобеней\n", brain.Users[k].Username, v)
	}
	return c.Send(res + "конец раунда, биатч")
}

func CmdVectorGetScores(c tele.Context, serv *servitor.Servitor, brain *BigBrain) error {
	report := "Очкобени\n"
	scores, err := serv.FreeMawVectorGetTopScores(10)
	if err != nil {
		return c.Send("проблема с получением очечей: %s", err.Error())
	}
	for _, s := range scores {
		report += brain.Users[s.UID].Username + " " + strconv.Itoa(s.Score) + "\n"
	}
	return c.Send(report)
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
	if len(longest) <= 2 {
		return "**", "**", "**"
	}
	if len(longest) <= 3 {
		return "***", "***", "***"
	}
	revealFirst := int(math.Round(float64(len(longest) * 20 / 100)))
	logger.Info("vector game estimates ",
		slog.Int("10% ", revealFirst))
	revealSecond := int(math.Round(float64(len(longest) * 30 / 100)))
	logger.Info("vector game estimates ",
		slog.Int("20% ", revealSecond))
	revealThird := int(math.Round(float64(len(longest) * 40 / 100)))
	logger.Info("vector game estimates ",
		slog.Int("40% ", revealThird))
	// Convert to rune slice to handle multi-byte characters properly
	if len(longest) == 4 {
		revealFirst = 1
		revealSecond = 1
		revealThird = 2
	}
	if len(longest) == 5 {
		revealFirst = 1
		revealSecond = 1
		revealThird = 2
	}
	if len(longest) == 6 {
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
