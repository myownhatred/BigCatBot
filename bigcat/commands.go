package bigcat

import (
	bigcat "Guenhwyvar/bigcat/games"
	"Guenhwyvar/config"
	"Guenhwyvar/entities"
	dnd "Guenhwyvar/lib/DND"
	"Guenhwyvar/lib/memser"
	"Guenhwyvar/servitor"
	"fmt"
	"io/ioutil"
	"log"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

// command handlers functions

const (
	// base
	HelpCmd      = "/help"
	HelpCmdFull  = "/help@GuenhwyvarBot"
	StartCmd     = "/start"
	ChangelogCmd = "/changelog"
	// progress
	Progress = "/progress"
	SetTask  = "/settask"
	GetTask  = "/gettask"
	Task     = "/task"
	AllTasks = "/all"
	// polls
	Dilemma = "/dilemma"
	// memes
	Meme   = "/meme"
	Hold   = "/hold"
	Choice = "/choice"
	// misc
	JokeHelp = "/joke"
	Clear    = "/clear"
	Ping     = "/ping"
	// maws
	GetAnimeOp = "/getopening"
	SetAnimeOp = "/setopening"
	// freemaws
	MawGet        = "/maw"
	MawAdd        = "/mawadd"
	MawRep        = "/mawreport"
	MawListInline = "/mawlist"
	AnimeMaw      = "/animemaw"
	GrobMaw       = "/grobmaw"
	// timers
	TimeWithOut = "/timewo"
	NewTimer    = "/newtimewo"
	// dnd stuff
	RollChar = "/rollcharhard"
	Party    = "/dndparty"
	Combat   = "/dndcombat"
	Attack   = "/dndattack222"
	Turn     = "/dndturn222"
	DnDJoin  = "/dndjoin"
	DnDStats = "/dndstats"
	DnDMF    = "/dndmf"
	// vector stuff
	VectorAddNewType = "/vectoraddtype"
	VectorGetTypes   = "/vectortypes"
	VectorAddNew     = "/vectoraddq"
	VectorGame       = "/vectorgame"
	// card stuff
	Card = "/card"
	// weather
	WeatherForecastDay = "/wday"
	WeatherCurrent     = "/weather"
	// steam
	GetFreeSteamGames = "/steam"
	// police mematron
	MetatronChatAdd     = "/chatadd"
	MetatronChatList    = "/chatlist"
	MetatronChatSend    = "/chatsend"
	MetatronChatForward = "/chatforward"
	// police achive
	UserAchList = "/achlist"
	UserAchAdd  = "/achtest"
	UserAchAdd2 = "/achtesttwo"
)

func CommandHandler(c tele.Context, serv *servitor.Servitor, flags *silly, brain *BigBrain, comfig *config.AppConfig, logger *slog.Logger) error {
	username := "АНОНИМ_ЛЕГИВОН"
	lastname := "Doe"
	if c.Sender().Username != "" {
		username = c.Sender().Username
	}
	if c.Sender().LastName != "" {
		lastname = c.Sender().LastName
	}
	logger.Info("incomingtext message",
		slog.Int64("chatID:", c.Chat().ID),
		slog.Int64("userID:", c.Sender().ID),
		slog.String("username:", username),
		slog.String("firstname:", c.Sender().FirstName),
		slog.String("lastname:", lastname),
		slog.String("message:", c.Message().Text))
	msgText := strings.Split(c.Message().Text, " ")

	// metatron check

	if c.Message().Private() {
		// check if user is on the bot/metatron list and set metatron flag on
		if _, ok := brain.UsersFlags[c.Sender().ID]; ok {
			logger.Info("user found:" + strconv.FormatInt(c.Sender().ID, 10))
			val := brain.UsersFlags[c.Sender().ID]
			if val.MetatronFordwardFlag {
				if msgText[0] == "/stop" {
					val.MetatronFordwardFlag = false
					brain.UsersFlags[c.Sender().ID] = val
					return c.Send("остановили форварденг")
				}
				//photo := c.Message().Photo
				return c.ForwardTo(&tele.Chat{ID: val.MetatronChat})
			}
		} else {
			logger.Info("user not found:" + strconv.FormatInt(c.Sender().ID, 10))
		}
	}
	// vector check
	if brain.ChatFlags[c.Chat().ID].VectorGame {
		logger.Info("vector game",
			slog.String("checking if message of this chact is really an answer, sender ", username))
		vc := brain.VectorGame[c.Chat().ID]
		if vc.CheckAnswer(strings.ToLower(c.Message().Text)) {
			brain.VectorChan <- username
		}
	}
	command := msgText[0]
	// check if link (twitter)
	if strings.HasPrefix(command, "https://twitter.com") || strings.HasPrefix(command, "https://x.com") {
		pathVid, err := serv.TwitterGetVideo(command)
		if err != nil {
			return c.Send(err.Error())
		}
		vido := &tele.Video{File: tele.FromDisk(pathVid)}
		return c.Send(vido)
	}
	// full command (command@botname) handling
	if len(strings.Split(command, "@")) > 1 {
		//fmt.Println("long command incoming")
		coms := strings.Split(command, "@")
		//fmt.Printf("command part is %s\n", coms[0])
		command = coms[0]
	}
	switch command {
	case StartCmd:
		return c.Send(msgHello)
	case HelpCmd:
		return c.Send(msgHelp)
	case HelpCmdFull:
		return c.Send(msgHelp)
	case ChangelogCmd:
		return c.Send(changelog)
	case Clear:
		menu = &tele.ReplyMarkup{RemoveKeyboard: true}
		return c.Send("чисти-чисти", menu)
	case Ping:
		_, err := c.Bot().Send(&tele.Chat{ID: comfig.MotherShip}, "pong")
		return err
	case Meme:
		return CreateGuiltyCatMeme(c, serv)
	case Hold:
		return holdmeme(c)
	case Choice:
		return CreateChoiceMeme(c)
	case Dilemma:
		return DilemmaPoll(c)
	case JokeHelp:
		return jokeHelp(c, comfig)
	case Progress:
		return c.Send(serv.GetWakaStuff())
	case GetAnimeOp:
		return SendAnimeOp(c, serv)
	case SetAnimeOp:
		return AnimeOpUpload(c, flags, serv)
	case AnimeMaw:
		return RandomAnimeMaw(c, serv)
	case MawListInline:
		return ShowInlineKeys(c)
	case TimeWithOut:
		return ShowTimeWithOut(c, serv)
	case NewTimer:
		return AddNewTimer(c, serv)
	case MawGet:
		return FreeMawGet(c, serv)
	case MawAdd:
		return FreeMawAdd(c, serv)
	case MawRep:
		return FreeMawRep(c, serv)
	case RollChar:
		return DnDRollChar(c, serv, brain)
	case Party:
		return DnDParty(c, serv, brain)
	case Combat:
		return DnDCombat(c, serv, brain)
	case Attack:
		return DnDAttack(c, serv, brain)
	case Turn:
		return DnDCombatTurn(c, serv, brain)
	case DnDJoin:
		return DnDJoinActive(c, serv, brain)
	case DnDStats:
		return DnDCharStats(c, brain)
	case DnDMF:
		return DnDCombat2(c, serv, brain)
	case Card:
		return GetRandomCard(c)
	case WeatherCurrent:
		return CmdWeatherCurrent(c, serv)
	case WeatherForecastDay:
		return CmdWeatherForecastDay(c, serv)
	case GetFreeSteamGames:
		return CmdGetFreeSteamGames(c, serv)
	case MetatronChatAdd:
		return CmdMetatronChatAdd(c, serv)
	case MetatronChatList:
		return CmdMetatronChatList(c, serv)
	case MetatronChatSend:
		return CmdMetatronChatSend(c, serv, username)
	case MetatronChatForward:
		return CmdMetatronChatForward(c, serv, brain, username)
	case UserAchList:
		return CmdUserAchList(c, serv)
	case UserAchAdd:
		return CmdUserAchAdd(c, serv)
	case VectorAddNewType:
		return CmdVectorAddNewType(c, serv)
	case VectorGetTypes:
		return CmdVectorGetTypes(c, serv)
	case VectorAddNew:
		return CmdVectorAddNew(c, serv)
	case VectorGame:
		return CmdVectorGame(c, serv, brain)
	default:
		return nil
	}

}

func CreateGuiltyCatMeme(c tele.Context, serv *servitor.Servitor) error {
	text := c.Message().Payload
	pathImg, err := serv.CreateGuiltyCatMeme(text)
	if err != nil {
		log.Printf("failed to generate meme %v\n", err)
		pathImg = "storage/notponyal.png"
	}

	poto := &tele.Photo{File: tele.FromDisk(pathImg)}
	return c.Send(poto)
}

func CreateChoiceMeme(c tele.Context) error {
	text := c.Message().Payload
	question := ""
	elements := strings.Split(text, ",")
	if len(elements) == 3 {
		question = elements[2]
	}
	if len(elements) <= 1 {
		return c.Send("для генерации мемчика дай строку вида левая кнопка,правая кнопка,вопрос внизу\nпоследний параметор апцинальный так то!")
	}
	pathImg, err := memser.Choice(elements[0], elements[1], question)
	if err != nil {
		c.Send("осибка генерации мемчика с выбором:%s", err.Error())
	}
	poto := &tele.Photo{File: tele.FromDisk(pathImg)}
	return c.Send(poto)
}

func DilemmaPoll(c tele.Context) error {
	text := c.Message().Payload
	question := ""
	elements := strings.Split(text, ",")
	if len(elements) == 3 {
		question = elements[2]
	}
	if len(elements) < 3 {
		return c.Send("для генерации опроса дай строку вида левая кнопка,правая кнопка,вопрос внизу")
	}
	pathImg, err := memser.Choice(elements[0], elements[1], question)
	if err != nil {
		c.Send("осибка генерации мемчика с выбором:%s", err.Error())
	}
	poto := &tele.Photo{File: tele.FromDisk(pathImg)}
	_ = c.Send(poto)

	left := tele.PollOption{
		Text: elements[0],
	}
	right := tele.PollOption{
		Text: elements[1],
	}

	poll := &tele.Poll{
		Type:     tele.PollRegular,
		Question: question,
		Options:  []tele.PollOption{left, right},
	}

	return c.Send(poll)

}

func SendAnimeOp(c tele.Context, serv *servitor.Servitor) error {
	file := tele.FromDisk(serv.GetOpeningsFilePath())
	dod := &tele.Document{
		File:     file,
		FileName: "openings.csv",
		Caption:  "опенинги",
	}
	return c.Send(dod)
}

func AnimeOpUpload(c tele.Context, flags *silly, serv *servitor.Servitor) error {
	// TODO: make argument of command pass by to servitor

	link := "https://raw.githubusercontent.com/myownhatred/Monolith/main/openings.csv"
	report, err := serv.UploadOpeningsByURL(link)
	if err != nil {
		report = err.Error()
	}
	return c.Send(report)
}

func holdmeme(c tele.Context) (err error) {
	text := c.Message().Payload
	pathImg, err := memser.HoldMeme(text)
	if err != nil {
		log.Printf("failed to generate hold meme %v\n", err)
		pathImg = "storage/hold.png"
	}

	poto := &tele.Photo{File: tele.FromDisk(pathImg)}
	return c.Send(poto)
}

func jokeHelp(c tele.Context, comfig *config.AppConfig) (err error) {
	msg := "sting"
	// TODO make this use filePath
	data, err := ioutil.ReadFile(comfig.JokePath)
	if err != nil {
		return err
	}
	msg = string(data)

	return c.Send(msg)
}

func RandomAnimeMaw(c tele.Context, serv *servitor.Servitor) (err error) {
	opening, err := serv.GetRandomOpening("anime")
	if err != nil {
		return c.Send(fmt.Sprintf("Неполучилось с опенингом: %s", err.Error()))
	}
	return c.Send(fmt.Sprintf("Тебе выпало послушать %s - %s", opening.Description, opening.Link))
}

func ShowInlineKeys(c tele.Context) (err error) {
	mawButs := &tele.ReplyMarkup{}
	btnAni := mawButs.Data("AnimeMaw", "animemaw")
	btnGrob := mawButs.Data("GrobMaw", "grobmaw")
	btnFree := mawButs.Data("FreeMaw", "freemaw")
	mawButs.Inline(mawButs.Row(btnAni, btnGrob, btnFree))
	return c.Send("Mavs", mawButs)
}

func ShowTimeWithOut(c tele.Context, serv *servitor.Servitor) (err error) {
	list, err := serv.GetTimeWithOutList(c.Chat().ID)
	if err != nil {
		return c.Send("не получилось со списком чот ")
	}
	incButtons := &tele.ReplyMarkup{ResizeKeyboard: true}

	message := "Список инцедентив:\n"
	var rows []tele.Row
	for id, item := range list {
		currentTime := time.Now()
		duration := currentTime.Sub(item.Time)
		days := int(duration.Hours()) / 24
		hours := int(duration.Hours()) % 24
		formattedDuration := fmt.Sprintf(" %02d дейз %02d хаурс", days, hours)
		rows = append(rows, incButtons.Row(incButtons.Data(fmt.Sprintf("%s: %s", butifulMumbers(id+1), formattedDuration), fmt.Sprintf("two%d", item.ID))))
		message += fmt.Sprintf("%s: %s %s\n", butifulMumbers(id+1), item.Name, formattedDuration)
	}
	message += "⬇️ нажми на номер чтобы сбросить счётчик ⬇️"
	rows = append(rows, incButtons.Row(incButtons.Data("скрыть", "sweep")))
	incButtons.Inline(rows...)
	return c.Send(message, incButtons)
}

func AddNewTimer(c tele.Context, serv *servitor.Servitor) (err error) {
	if c.Message().Payload == "" {
		return c.Send("пазязя дайте название событию!")
	}
	err = serv.AddNewTimer(c.Message().Payload, c.Chat().ID)
	message := fmt.Sprintf("%s, новий беспроблемный таймер создан!", c.Message().Sender.Username)
	if err != nil {
		message = fmt.Sprintf("вышла ошибочка: %s", err.Error())
	}
	return c.Send(message)
}

func FreeMawAdd(c tele.Context, serv *servitor.Servitor) (err error) {
	// empty payload means open typ
	// TODO: add typ parsing to payload
	// OR make it inline option dunno
	typ := "open"
	// finding the link
	// TODO add link validity check
	if c.Message().Payload == "" {
		return c.Send("пазязя дайте название сонга и ссылоцу!")
	}
	// TODO add flag option to get description and link
	// in next message
	mawargs := strings.Split(c.Message().Payload, "http")
	if len(mawargs) != 2 {
		return c.Send("в вашем сообсении должна быть аписание сонга и потом ссылоцка!")
	}
	maw := entities.FreeMaw{
		Type:        typ,
		Description: mawargs[0],
		Link:        "http" + mawargs[1],
	}
	err = serv.PutFreeMaw(maw)
	if err != nil {
		return c.Send("произосло: " + err.Error())
	}
	return c.Send("добавели васу песенку!")
}

func FreeMawGet(c tele.Context, serv *servitor.Servitor) (err error) {
	typ := "open"
	maw, err := serv.GetFreeMaw(typ)
	if err != nil {
		return c.Send("произосло:" + err.Error())
	}
	return c.Send(fmt.Sprintf("слусай %s %s", maw.Description, maw.Link))
}

func CmdWeatherCurrent(c tele.Context, serv *servitor.Servitor) (err error) {
	if c.Message().Payload == "" {
		return c.Send("пазязя дайте название дерёвни для узнания погоды!")
	}
	report, err := serv.GetCurrentWeather(c.Message().Payload)
	if err != nil {
		return c.Send("произосло: " + err.Error())
	}
	return c.Send(report)
}

func CmdWeatherForecastDay(c tele.Context, serv *servitor.Servitor) (err error) {
	if c.Message().Payload == "" {
		return c.Send("пазязя дайте название дерёвни для узнания погоды!")
	}
	report, err := serv.GetWeatherDayForecast(c.Message().Payload)
	if err != nil {
		return c.Send("произосло: " + err.Error())
	}
	return c.Send(report)
}

func CmdGetFreeSteamGames(c tele.Context, serv *servitor.Servitor) (err error) {
	report, err := serv.GetFreeSteamGames()
	if err != nil {
		return c.Send("пли вызовые халяви стима произосло: " + err.Error())
	}
	return c.Send(report)
}

func CmdMetatronChatAdd(c tele.Context, serv *servitor.Servitor) (err error) {
	err = serv.MetatronChatAdd(c.Chat().ID, c.Chat().Title)
	if err != nil {
		return c.Send("сомесинг вронг аддинг зыс чат ту зе метатрон лист: " + err.Error())
	} else {
		return c.Send("chat added to hyperconnection communnnication ultranetwork")
	}
}

func CmdMetatronChatList(c tele.Context, serv *servitor.Servitor) (err error) {
	IDs, _, Names, err := serv.MetatronChatList()
	if err != nil {
		return c.Send("сомесинг вент вронг листинг чатс: " + err.Error())
	} else {
		report := ""
		for i, id := range IDs {
			report += fmt.Sprintf("%d - %s\n", id, Names[i])
		}
		return c.Send(report)
	}
}

func CmdMetatronChatSend(c tele.Context, serv *servitor.Servitor, username string) (err error) {
	if c.Message().Payload == "" {
		return c.Send("введите номер чатика из списька а потом сообщеньку")
	}
	cmdargs := strings.Split(c.Message().Payload, " ")
	num, err := strconv.Atoi(cmdargs[0])
	if err != nil {
		return c.Send("введите номер чатика из списька а потом сообщеньку")
	}
	if num == 0 {
		return c.Send("номер чатика должен быть большее нулика")
	}
	// prepare message and chat coordinates
	result := ""
	for i, str := range cmdargs {
		if i != 0 {
			result = strings.Join([]string{result, str}, " ")
		}
	}
	IDs, ChatIDs, Names, err := serv.MetatronChatList()
	if len(IDs) < num {
		return c.Send("номер чатика должен быить из списка чатиков /chatlist")
	}
	if c.Chat().Title == "" {
		result += fmt.Sprintf("\n\n^^^^^^^^^\nот %s", username)
	} else {
		result += fmt.Sprintf("\n\n^^^^^^^^^\nот %s\nИз чата %s", username, c.Chat().Title)
	}
	c.Bot().Send(&tele.Chat{ID: ChatIDs[num-1]}, result)
	report := fmt.Sprintf("Отправили ваше письмо в чат %s", Names[num-1])
	return c.Send(report)
}

func CmdMetatronChatForward(c tele.Context, serv *servitor.Servitor, brain *BigBrain, username string) (err error) {
	if c.Message().Payload == "" {
		return c.Send("введите номер чатика для форварьдинга из списька")
	}
	cmdargs := strings.Split(c.Message().Payload, " ")
	num, err := strconv.Atoi(cmdargs[0])
	if err != nil {
		return c.Send("введите номер чатика из списька")
	}
	if num == 0 {
		return c.Send("номер чатика должен быть большее нулика")
	}
	IDs, ChatIDs, _, err := serv.MetatronChatList()
	if len(IDs) < num {
		return c.Send("номер чатика должен быить из списка чатиков /chatlist")
	}
	// add forward flag and chat ID to bot brain
	rule := UserRules{
		MetatronChat:         ChatIDs[num-1],
		MetatronFordwardFlag: true,
	}
	brain.UsersFlags[c.Sender().ID] = rule
	return c.Send("форварденг включин")
}

func CmdUserAchList(c tele.Context, serv *servitor.Servitor) (err error) {
	SenderID := c.Sender().ID
	report := ""
	//	IDs, UserIDs, AchIDs, Dates, Chats, ChatIDs := serv.UserAchs(SenderID)
	IDs, _, AchIDs, _, _, _, err := serv.UserAchs(SenderID)
	if err != nil {
		return c.Send(err.Error())
	}
	// TODO remake achive list to map?
	OAchIDs, _, ONames, _, ODescriptions, err := serv.Achieves(0)
	// range over IDs of user's achives
	for _, ID := range AchIDs {
		// range over Achieves table to get info for each ID
		for j, OID := range OAchIDs {
			if ID == OID {
				report += "⭐️" + ONames[j] + "⭐️ - " + ODescriptions[j] + "\n"
			}
		}
	}
	report += fmt.Sprintf("у юзера всего ачив %d", len(IDs))
	return c.Send(report)
}

func DnDParty(c tele.Context, serv *servitor.Servitor, brain *BigBrain) (err error) {
	message := "наша ⚔️мощная⚔️ ватага:\n"
	for _, pers := range brain.Party {
		message += pers.Name + " " + pers.Title + " : " + string(pers.Race) + "-" + string(pers.Class) + "\n"
	}
	message += "вы находитесь возле деревне скрытого листа и можете идти в такие места: "
	incButtons := &tele.ReplyMarkup{ResizeKeyboard: true}

	var rows []tele.Row
	rows = append(rows, incButtons.Row(incButtons.Data("Бар \"Пьяный Шакал\"", "DNDtoBar")))
	rows = append(rows, incButtons.Row(incButtons.Data("Площадь деревни скрытого листа", "DNDtoPlaza")))

	rows = append(rows, incButtons.Row(incButtons.Data("скрыть", "sweep")))
	incButtons.Inline(rows...)
	return c.Send(message, incButtons)
}

func DnDCombatTurn(c tele.Context, serv *servitor.Servitor, brain *BigBrain) (err error) {
	if !brain.Game.CombatFlag {
		return c.Send("зачем ходить комбат если комбата и нет?")
	}
	message := "летс мортар комбат бегинс\n"
	for _, char := range brain.Game.CombatOrder {
		if char.IsNPC {
			// stupid way to pick target for npc
			for _, v := range brain.Game.Party {
				endFlag := true
				for validIndex, validTarget := range brain.Game.CombatOrder {
					if v.Name == validTarget.Name && validTarget.Hitpoints > 0 {
						endFlag = false
						message += fmt.Sprintf("%s бьёт по %s\n", char.Name, validTarget.Name)
						dmg, messagedmg := char.GetAttackDamage(validTarget.AC)
						message += messagedmg
						message += fmt.Sprintf("\nхп цели: %d - %d = %d\n", brain.Game.CombatOrder[validIndex].Hitpoints, dmg, brain.Game.CombatOrder[validIndex].Hitpoints-dmg)
						brain.Game.CombatOrder[validIndex].Hitpoints -= dmg
						//brain.Game.Party[c.Sender().ID] = *brain.Game.CombatOrder[validIndex]
						if brain.Game.CombatOrder[validIndex].Hitpoints <= 0 {
							message += "цель perished\n"
						}
					}
				}
				if endFlag {
					message += "пахот нпц винс"
					brain.Game.CombatFlag = false
				}
				break
			}
		}
		if char.Target == nil {
			continue
		}
		if char.Hitpoints <= 0 {
			continue
		}
		target := char.Target
		message += fmt.Sprintf("%s бьёт по %s\n", char.Name, target.Name)
		dmg, messagedmg := char.GetAttackDamage(target.AC)
		message += messagedmg
		message += fmt.Sprintf("\nхп цели: %d - %d = %d\n", target.Hitpoints, dmg, target.Hitpoints-dmg)
		target.Hitpoints -= dmg
		if target.Hitpoints <= 0 {
			message += "цель perished\n"
			if target.Name == "Керилл" || target.Name == "Васян" {
				message += "Со смертью этого персонажа комбат в этом месте заканчивается, идите в другое или живите дальше в проклятом мире, который сами и создали\n"
				brain.Game.CombatFlag = false
				break
			}
		}

	}
	return c.Send(message)
}

func DnDAttack(c tele.Context, serv *servitor.Servitor, brain *BigBrain) (err error) {
	if !brain.Game.CombatFlag {
		return c.Send("битва не может начатться, сделайте /dndcombat")
	}
	if c.Message().Payload == "" {
		return c.Send("введите номер залупатора из списка /dndcombat")
	}
	num, err := strconv.Atoi(c.Message().Payload)
	if err != nil {
		return c.Send("ваш интеллект явно ниже 10, введите чесло!")
	}
	if num < 0 || num > len(brain.Game.CombatOrder) {
		return c.Send("номер залупатара вне грониц /dndcombat")
	}
	me := brain.Game.Party[c.Sender().ID]

	// assign target in order array (cringe)
	meIndex := 0
	for ind, char := range brain.Game.CombatOrder {
		if me.Name == char.Name {
			meIndex = ind
			break
		}
	}
	if brain.Game.CombatOrder[meIndex].Hitpoints <= 0 {
		return c.Send("вы мертвы и не можете совершать действия")
	}
	target := brain.Game.CombatOrder[num-1]
	if target.Hitpoints <= 0 {
		return c.Send("ваша цель метрва")
	}
	brain.Game.CombatOrder[meIndex].Target = brain.Game.CombatOrder[num-1]
	message := ""
	if me.Name == target.Name {
		message += fmt.Sprintf("%s решил хуярит сам по себе (чистый термояд-дегенерат)\n", me.Name)
	} else {
		message += fmt.Sprintf("%s выбрал целью %s\n", me.Name, target.Name)
	}
	return c.Send(message)
}

func DnDAttackFC(c tele.Context, brain *BigBrain, num int) (err error) {
	if !brain.Game.CombatFlag {
		return c.Send("битва не может начатться, сделайте /dndcombat")
	}
	if num < 0 || num > len(brain.Game.CombatOrder) {
		return c.Send("номер залупатара вне грониц /dndcombat")
	}
	me := brain.Game.Party[c.Sender().ID]

	// assign target in order array (cringe)
	meIndex := 0
	for ind, char := range brain.Game.CombatOrder {
		if me.Name == char.Name {
			meIndex = ind
			break
		}
	}
	if brain.Game.CombatOrder[meIndex].Hitpoints <= 0 {
		return c.Send("вы мертвы и не можете совершать действия")
	}
	target := brain.Game.CombatOrder[num-1]
	if target.Hitpoints <= 0 {
		return c.Send("ваша цель метрва")
	}
	brain.Game.CombatOrder[meIndex].Target = brain.Game.CombatOrder[num-1]
	message := ""
	if me.Name == target.Name {
		message += fmt.Sprintf("%s решил хуярит сам по себе (чистый термояд-дегенерат)\n", me.Name)
	} else {
		message += fmt.Sprintf("%s выбрал целью %s\n", me.Name, target.Name)
	}
	message += fmt.Sprintf("%s бьёт по %s\n", me.Name, target.Name)
	dmg, messagedmg := me.GetAttackDamage(target.AC)
	message += messagedmg
	message += fmt.Sprintf("\nхп цели: %d - %d = %d\n", target.Hitpoints, dmg, target.Hitpoints-dmg)
	target.Hitpoints -= dmg
	if target.Hitpoints <= 0 {
		message += "цель perished\n"
		if target.Name == "Керилл" || target.Name == "Васян" {
			message += "Со смертью этого персонажа комбат в этом месте заканчивается, идите в другое или живите дальше в проклятом мире, который сами и создали\n"
			brain.Game.CombatFlag = false
		}
	}
	return c.Send(message)
}

func DnDListActionButtons(c tele.Context, serv *servitor.Servitor, brain *BigBrain) (err error) {
	message := "Выберите действие\n"
	incButtons := &tele.ReplyMarkup{ResizeKeyboard: true}
	var rows []tele.Row
	rows = append(rows, incButtons.Row(incButtons.Data("Атака", "dndAttacks")))
	rows = append(rows, incButtons.Row(incButtons.Data("Трюки", "dndCantrips")))
	rows = append(rows, incButtons.Row(incButtons.Data("Заклинания", "dndSpells")))
	rows = append(rows, incButtons.Row(incButtons.Data("скрыть", "sweep")))
	incButtons.Inline(rows...)
	return c.Send(message, incButtons)
}

func DnDAttackButtons(c tele.Context, serv *servitor.Servitor, brain *BigBrain) (err error) {
	message := "Выберите тип атаки\n"
	incButtons := &tele.ReplyMarkup{ResizeKeyboard: true}
	var rows []tele.Row
	rows = append(rows, incButtons.Row(incButtons.Data("Мили", "dndMeleeAttack")))
	rows = append(rows, incButtons.Row(incButtons.Data("Дальняя", "dndRangeAttack")))
	rows = append(rows, incButtons.Row(incButtons.Data("Заклинания", "dndSpells")))
	rows = append(rows, incButtons.Row(incButtons.Data("скрыть", "sweep")))
	incButtons.Inline(rows...)
	return c.Send(message, incButtons)
}

func DnDTargetsButtons(c tele.Context, serv *servitor.Servitor, brain *BigBrain) (err error) {
	message := "Выберите цель\n"
	incButtons := &tele.ReplyMarkup{ResizeKeyboard: true}

	var rows []tele.Row
	for id, item := range brain.Game.CombatOrder {
		rows = append(rows, incButtons.Row(incButtons.Data(fmt.Sprintf("%d: %s", id+1, item.Name), fmt.Sprintf("dndAttackTarget%d", id))))

	}
	rows = append(rows, incButtons.Row(incButtons.Data("скрыть", "sweep")))
	incButtons.Inline(rows...)
	return c.Send(message, incButtons)
}

func CmdUserAchAdd(c tele.Context, serv *servitor.Servitor) (err error) {
	SenderID := c.Sender().ID
	username := "АНОНИМ_ЛЕГИВОН"
	lastname := "Doe"
	if c.Sender().Username != "" {
		username = c.Sender().Username
	}
	if c.Sender().LastName != "" {
		lastname = c.Sender().LastName
	}
	firstname := c.Sender().FirstName
	_, err = serv.UserDefaultCheck(SenderID, username, firstname, lastname, "/achtest")
	if err != nil {
		return c.Send(err.Error())
	}
	AchID, err := serv.UserAchAdd(SenderID, 1, c.Chat().Title, c.Chat().ID)
	if err != nil {
		return c.Send(err.Error())
	}
	if AchID == 0 {
		return c.Send("успесно добавили тестовую очивочку")
	} else {
		return c.Send("у вас узе видимо есть эта ацивоцка")
	}
}

func FreeMawRep(c tele.Context, serv *servitor.Servitor) (err error) {
	return c.Send("report is not awalablash")
}

func DnDRollChar(c tele.Context, serv *servitor.Servitor, brain *BigBrain) error {
	SenderID := c.Sender().ID
	username := "АНОНИМ_ЛЕГИВОН"
	lastname := "Doe"
	if c.Sender().Username != "" {
		username = c.Sender().Username
	}
	if c.Sender().LastName != "" {
		lastname = c.Sender().LastName
	}
	firstname := c.Sender().FirstName
	_, err := serv.UserDefaultCheck(SenderID, username, firstname, lastname, "/rollcharhard")
	if err != nil {
		return c.Send(err.Error())
	}

	chel := dnd.RollChar()
	chel.Name = username
	//ach, title, desc, achID := DNDStatsAchievement(str, dex, con, inn, wis, cha)
	ach, title, desc, achID := DNDStatsAchievement(chel.CharStats())
	message2 := c.Message().Sender.Username + " твой перец(а/я/мы):\n"
	message2 += chel.Generation
	message2 += "Вооружон " + string(chel.Weapon.Name) + "\n"
	message2 += "Адет " + string(chel.Armor.Name) + "\n"

	if ach == "" {
		brain.Party[SenderID] = chel
		brain.Game.Party[SenderID] = chel
		return c.Send(message2)
	} else {
		message2 += "Твой титул: " + title + "\n"
		message2 += "И у тебя ачивка: " + ach + "\n"
		message2 += "Описание ачивки: " + desc + "\n"
		if achID != 0 {
			_, _ = serv.UserAchAdd(SenderID, achID, c.Chat().Title, c.Chat().ID)
		}
		chel.Title = title
		brain.Party[SenderID] = chel
		brain.Game.Party[SenderID] = chel
		return c.Send(message2)
	}
}

func GetRandomCard(c tele.Context) error {
	var deck bigcat.Deck
	deck.NewStack()
	card := deck.TopDeck()
	get := card.OneShot()

	return c.Send(get)
}

func dice3of4() (scrib string) {
	rand.Seed(time.Now().UnixNano())
	min := rand.Intn(6) + 1
	summ := min
	scrib += strconv.Itoa(min)
	for i := 0; i < 3; i++ {
		k := rand.Intn(6) + 1
		summ += k
		if k < min {
			min = k
		}
		scrib += " + " + strconv.Itoa(k)
	}
	scrib += " = " + strconv.Itoa(summ) + " => " + strconv.Itoa(summ-min)
	return scrib
}
func dice3of4i() (res int) {
	rand.Seed(time.Now().UnixNano())
	min := rand.Intn(6) + 1
	summ := min
	for i := 0; i < 3; i++ {
		k := rand.Intn(6) + 1
		summ += k
		if k < min {
			min = k
		}
	}
	return summ - min
}

func butifulMumbers(i int) string {
	switch i {
	case 1:
		return "\u0031\uFE0F\u20E3"
	case 2:
		return "\u0032\uFE0F\u20E3"
	case 3:
		return "\u0033\uFE0F\u20E3"
	case 4:
		return "\u0034\uFE0F\u20E3"
	case 5:
		return "\u0035\uFE0F\u20E3"
	case 6:
		return "\u0036\uFE0F\u20E3"
	case 7:
		return "\u0037\uFE0F\u20E3"
	case 8:
		return "\u0038\uFE0F\u20E3"
	case 9:
		return "\u0039\uFE0F\u20E3"
	default:
		return strconv.Itoa(i)
	}
}

func DNDStatsAchievement(str, dex, con, inn, wis, cha int) (achive, title, desc string, achID int) {
	arr := [6]int{str, dex, con, inn, wis, cha}
	// brackets
	p18 := 0
	p16 := 0
	p12 := 0
	b14_15 := 0
	b12_15 := 0
	b10_14 := 0
	p14 := 0
	m12 := 0
	m10 := 0
	m8 := 0
	m6 := 0
	m4 := 0
	equals := 0

	for k, i := range arr {
		for j, z := range arr {
			if j == k {
				continue
			}
			if i == z {
				equals++
			}
		}
		if i >= 18 {
			p18++
		}
		if i >= 16 {
			p16++
		}
		if i >= 14 {
			p14++
		}
		if i >= 12 {
			p12++
		}
		if i <= 12 {
			m12++
		}
		if 15 >= i && i >= 14 {
			b14_15++
		}
		if 15 >= i && i >= 12 {
			b12_15++
		}
		if 14 >= i && i >= 10 {
			b10_14++
		}

		if i <= 10 {
			m10++
		}
		if i <= 8 {
			m8++
		}
		if i <= 6 {
			m6++
		}
		if i <= 4 {
			m4++
		}
	}
	// Perfect Balance
	if str == 18 && dex == 18 && con == 18 && inn == 18 && wis == 18 && cha == 18 {
		return "Перфектный Баланц", "Парагон", "Все статы по 18", 6
	}
	// The Mastermind
	if (p16 == 4 && b14_15 == 2) || (p16 == 5 && b14_15 == 1) {
		return "Мастерум", "Виртуоззо", "4/5 статов 16+ и 2/1 14-15", 7
	}
	// Well-rounded
	if p16 == 3 && b14_15 == 3 {
		return "Хорошо-круглый", "Герой ренисанса", "3 статы 16+ и 3 14-15", 8
	}
	if p16 == 3 {
		return "Высокий роллер", "Удачливая легенда", "3 статы 16+", 9
	}
	if str >= 16 && dex >= 16 && con >= 16 {
		return "Физический мощнодом", "Неостановимый джаггернихт", "сила, ловкость и конста 16+", 10
	}
	if inn >= 16 && wis >= 16 && cha >= 16 {
		return "Ментальное мастерство", "Саг", "интеллект, мудрость и харизма 16+", 11
	}
	if p18 == 1 && b12_15 == 5 {
		return "Спецальный", "Специлист", "одна характеристика 18 и остальные 12-15", 12
	}
	if p18 == 1 && p12 == 6 {
		return "Удачный прорыв", "Счастливая находка", "одна характеристика 18 и остальные 12+", 13
	}
	if p18 >= 1 && m6 >= 1 {
		return "Разбалансированный зверь", "Разбалансированная подсобака", "одна характеристика 18 и одна 6-", 14
	}
	if m10 == 1 && p14 == 5 {
		return "Небалансный", "Дикая карта", "одна характеристика меньше 10 остальные 14+", 15
	}
	if m10 == 6 {
		return "Низские начала", "Подсобака", "все характеристики 10-", 16
	}
	if equals == 0 {
		return "Уникальный снежок", "Энигма", "значения всех характеристик различны", 3
	}
	if b12_15 == 6 {
		return "Золотая середина", "Балансировочное лезвие", "все характеристики 12-15", 17
	}
	if b10_14 == 6 {
		return "Чудо середины дороги", "Средний приключенец", "все характеристики 10-14", 18
	}
	if inn >= 16 && wis >= 16 {
		return "Мозговая Правда", "Интелектуарный гигант", "интелект и мудрость 16+", 19
	}
	if cha == 18 && b12_15 == 5 {
		return "Харизматичный лидер", "Народный чемпион", "харизма 18 и остальные 12-15", 20
	}
	// the brawler - the unyielding warrior
	if str >= 16 && con >= 16 {
		return "Браулер", "Неелдующий воин", "сила и конста 16+", 21
	}
	if dex == 18 && str <= 6 {
		return "Подвижный якорь", "Быстрый страйкер", "ловкость 18 и сила 6-", 22
	}
	if dex == 18 && b12_15 == 5 {
		return "Подвижный туз", "Быстрый страйкер со слабым пятном", "ловкость 18 и остальные 12-15", 23
	}
	if p16 == 0 && m10 >= 3 {
		return "Жека без торговлей", "Борцующийся выживальщик", "ни одной характеристики выше 16 и три 10-", 4
	}
	if p18 == 1 && m10 == 5 {
		return "Непохоже что гирой", "Случайный чемпион", "одна характеристика 18 и остальные 10-", 24
	}
	if str == 18 && con == 18 {
		return "Скульптурный", "Недвижимый объект", "сила и конста 18", 25
	}
	if str == 18 && con == 18 && dex <= 6 {
		return "Скульптурный спотыкач", "Неуклюжий гигант", "сила и конста 18 а ловкость 6-", 26
	}
	if inn >= 16 && dex >= 16 {
		return "Церебральный убийца", "Кунирующий киллер", "интеллект и ловкость 16+", 27
	}
	if inn >= 16 && wis >= 16 && cha <= 6 {
		return "Церебральный спотыкач", "Нехаризматичный гений", "интеллекти и мудрость 16+ и харизма 6-", 28
	}
	if m8 == 6 {
		return "Неудачливая душа", "Проклятый путещественник", "все характеристики меньше 8", 29
	}
	if str <= 6 && m10 == 5 {
		return "Слабый боец", "Жалкий боксёр", "сила 6- и остальные характеристики 10-", 30
	}
	if dex <= 6 && m10 == 5 {
		return "Неуклюжий придурок", "Храбрый с пальцами", "ловкость 6- и остальные характеристики 10-", 31
	}
	if inn <= 6 && m10 == 5 {
		return "Тусклая лампочка", "Тупой сорвиголова", "интеллект 6- и остальные характеристики 10-", 32
	}
	if wis <= 6 && m10 == 5 {
		return "Наивный новичок", "Доверчивый игрок", "мудрость 6- и остальные характеристики 10-", 33
	}
	if cha <= 6 && m10 == 5 {
		return "Невдохновляющий лидер", "Непопулярный герой", "харизма 6- и остальные характеристики 10-", 34
	}
	if con <= 6 && m10 == 5 {
		return "Хрупкий цветочек", "Деликатный дорогуша", "конста 6- и остальные характеристики 10-", 35
	}
	if p12 == 0 && m8 >= 3 {
		return "Всеобще ужасный", "Полный провал", "нет характеристик выше 12 и три или больше характеристик 8-", 36
	}
	if m6 == 6 {
		return "Бездонный искатель приключений", "Жалкий исследователь", "все характеристики ниже 6", 37
	}
	if str <= 4 && m10 == 6 {
		return "Бессильный бродяга", "Бессильный боксёр", "сила 4- и остальные характеристики 10-", 38
	}
	if dex <= 4 && m10 == 6 {
		return "Неуклюжая катастрофа", "Каламутный скалолаз", "ловкость 4- и остальные характеристики 10-", 39
	}
	if inn <= 4 && m10 == 6 {
		return "Интелектуально вызванный", "Глупенький савант", "интеллект 4- и остальные характеристики 10-", 40
	}
	if wis <= 4 && m10 == 6 {
		return "Бездумный бродяга", "Наивный кочевник", "мудрость 4- и остальные характеристики 10-", 41
	}
	if cha <= 4 && m10 == 6 {
		return "Стрёмный неудачник", "Непопулярный недостигатор", "харизма 4- и остальные характеристики 10-", 42
	}
	if con <= 4 && m10 == 6 {
		return "Хрупкая неудача", "Ломаный храбрец", "конста 4- и остальные характеристики 10-", 43
	}
	if dex <= 6 && str <= 6 {
		return "Нескоординированный придурок", "Неуклюжий сгусток", "ловкость и сила 6-", 44
	}
	if inn <= 6 && wis <= 6 {
		return "Дуплет тупизны", "Глупый собиратель", "интеллект и мудрость 6-", 45
	}
	if inn <= 6 && wis <= 6 && cha <= 6 {
		return "Триплет тупизны", "Непопулярный контрдостигатель", "интеллект мудрость и харизма 6-", 46
	}
	if con <= 6 && str <= 6 && dex <= 6 && inn <= 6 {
		return "Четверной хрупок", "Мегаломающийся храбрый", "конста сила ловкость и интеллект 6-", 47
	}
	if m4 == 6 {
		return "Шестикратная угроза", "Полная катастрофа", "все характеристики 4-", 48
	}
	if m8 == 5 && m4 == 1 {
		return "Невезучая подсобака", "Неудачный контрдостигатель", "одна характеристика 4- и остальные 8-", 49
	}
	if m6 == 5 && m4 == 1 {
		return "Горестный негодяй", "Жалкий исследователь", "одна характеристика 4- и остальные 6-", 50
	}
	if m6 == 6 {
		return "Эксперт горячего замеса", "Каламутный крестодёр", "все характеристики 6-", 51
	}
	if m10 == 5 && str <= 4 {
		return "Кирпичная стена", "Неходящий объект", "сила 4- и остальные характеристики 10-", 52
	}
	if m10 == 5 && dex <= 4 {
		return "Неуклюжий кинг", "Неуклюжий чемпион", "ловкость 4- и остальные характеристики 10-", 53
	}
	if m10 == 5 && inn <= 4 {
		return "Мозговой слив", "Тускломуный дипломат", "интелект 4- и остальные характеристики 10-", 54
	}
	if m10 == 5 && wis <= 4 {
		return "Космический кадет", "Глупый собиратель", "мудрость 4- и остальные характеристики 10-", 55
	}
	if m10 == 5 && cha <= 4 {
		return "Социальный изгой", "Непопулярный контрдостигатель", "харизма 4- и остальные характеристики 10-", 56
	}
	if m10 == 5 && con <= 4 {
		return "Стеклянная пушка", "Мегаломающийся храбрый", "конста 4- и остальные характеристики 10-", 57
	}
	if m10 == 3 && m6 == 3 {
		return "Трио поездокрушения", "Каламутный крестодёр", "три характеристики 10- и три другие 6-", 58
	}
	if con <= 6 && str <= 6 && dex <= 6 {
		return "Несвятая троица", "Раскоординированный монстр", "конста сила и ловкость 6-", 59
	}
	if inn <= 6 && wis <= 6 && cha <= 6 {
		return "Подрыв доверия к мозгу", "Непопулярный контрдостигатель", "интеллект мудрость и харизма 6-", 60
	}
	if m4 == 6 {
		return "Пятикратная угроза себе", "Полная катастрофа", "все характеристики 4-", 61
	}
	if m8 == 6 && m4 == 1 {
		return "Попытка всплытия", "Полная катастрофа", "одна характеристика 4- и остальные 8-", 62
	}
	if p16 >= 1 && m8 >= 1 {
		return "Высокий и низкий герой", "Непредсказуемая подсобака", "одна характеристика 16+ и одна 8-", 5
	}
	if p18 >= 1 && m6 >= 1 {
		return "Разбалансированный зверь", "Разбалансированная подсобака", "одна характеристика 18 и одна 6-", 63
	}
	if inn == 18 && m10 == 5 {
		return "Нестабильный гений", "Безумный учоный", "интеллект 18 и остальные характеристики 10-", 64
	}
	if cha == 18 && m10 == 5 {
		return "Харизматичный шансовщик", "Серебрянноязыкий дъиявол", "харизма 18 и остальные характеристики 10-", 65
	}
	if dex == 18 && m10 == 5 {
		return "Подвижный акробат", "Бысрый ударщик", "ловкость 18 и остальные характеристики 10-", 66
	}
	if str <= 4 && m12 == 5 {
		return "Бессильное чудо", "Дрыщеватый боец", "сила 4- и остальные характеристики 12-", 67
	}
	if dex <= 4 && m12 == 6 {
		return "Неуклюжая катастрофа", "Щас-что-то-будет", "ловкость 4- и остальные характеристики 12-", 68
	}
	if p16 == 3 && m6 == 1 {
		return "Высокий роллер проклятый", "неУдачливая легенда", "три характеристики 16+ и одна 6-", 69
	}

	return "", "", "", 0
}
