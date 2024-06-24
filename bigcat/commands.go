package bigcat

import (
	bigcat "Guenhwyvar/bigcat/games"
	"Guenhwyvar/config"
	"Guenhwyvar/entities"
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
	HelpCmd     = "/help"
	HelpCmdFull = "/help@GuenhwyvarBot"
	StartCmd    = "/start"
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
	// card stuff
	Card = "/card"
	// weather
	WeatherForecastDay = "/wday"
	WeatherCurrent     = "/weather"
	// steam
	GetFreeSteamGames = "/steam"
	// police
	MetatronChatAdd     = "/chatadd"
	MetatronChatList    = "/chatlist"
	MetatronChatSend    = "/chatsend"
	MetatronChatForward = "/chatforward"
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
			logger.Info("user found:", c.Sender().ID)
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
			logger.Info("user not found:", c.Sender().ID)
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
		return DnDRollChar(c)
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

func FreeMawRep(c tele.Context, serv *servitor.Servitor) (err error) {
	return c.Send("report is not awalablash")
}

func DnDRollChar(c tele.Context) error {
	rand.Seed(time.Now().UnixNano())
	gender := []string{"Спермобак", "Вагинокапиталист", "Мунгендер", "Агендер", "Гендервой", "Нонбайори"}
	racces := []string{"Дворф", "Халфлинг", "Хуманс", "Эльф", "Гнум", "Драгонборн", "Полуорк", "Полуэльф", "Тифлинг"}
	clases := []string{"Бесполезный", "Барбариан", "Солдат", "Визард", "Друль", "Жрец", "Колдун", "Монк", "ПаллАдин", "Шельма", "Следопыт", "Военный Замок"}
	message := c.Message().Sender.Username + " твой перец(а/я/мы):\n"
	message += "Гендир: " + gender[rand.Intn(len(gender))] + "\n"
	message += "Расса: " + racces[rand.Intn(len(racces))] + "\n"
	message += "Клас: " + clases[rand.Intn(len(clases))] + "\n"
	str := dice3of4i()
	message += "Сила: " + strconv.Itoa(str) + "\n"
	dex := dice3of4i()
	message += "Ловкость: " + strconv.Itoa(dex) + "\n"
	con := dice3of4i()
	message += "Телосложение: " + strconv.Itoa(con) + "\n"
	inn := dice3of4i()
	message += "Интеллект: " + strconv.Itoa(inn) + "\n"
	wis := dice3of4i()
	message += "Мудрость: " + strconv.Itoa(wis) + "\n"
	cha := dice3of4i()
	message += "Харя: " + strconv.Itoa(cha) + "\n"
	ach, title := DNDStatsAchievement(str, dex, con, inn, wis, cha)
	if ach == "" {
		return c.Send(message)
	} else {
		message += "Твой титул: " + title + "\n"
		message += "И у тебя ачивка: " + ach + "\n"
		return c.Send(message)
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

func DNDStatsAchievement(str, dex, con, inn, wis, cha int) (achive, title string) {
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
		if i >= 12 {
			p12++
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
		if i >= 14 {
			p14++
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
		return "Перфентный Баланц", "Парагон"
	}
	// The Mastermind
	if (p16 == 4 && b14_15 == 2) || (p16 == 5 && b14_15 == 1) {
		return "Мастерум", "Виртуоззо"
	}
	// Well-rounded
	if p16 == 3 && b14_15 == 3 {
		return "Хорошо-круглый", "Герой ренисанса"
	}
	if p16 == 3 {
		return "Высокий роллер", "Удачливая легенда"
	}
	if str >= 16 && dex >= 16 && con >= 16 {
		return "Физический мощнодом", "Неостановимый джаггернихт"
	}
	if inn >= 16 && wis >= 16 && cha >= 16 {
		return "Ментальное мастерство", "Саг"
	}
	if p18 == 1 && b12_15 == 5 {
		return "Спецальный", "Специлист"
	}
	if p18 == 1 && p12 == 5 {
		return "Удачный прорыв", "Счастливая находка"
	}
	if p18 >= 1 && m6 >= 1 {
		return "Разбалансированный зверь", "Разбалансированная подсобака"
	}
	if m10 == 1 && p14 == 5 {
		return "Небалансный", "Дикая карта"
	}
	if m10 == 6 {
		return "Низские начала", "Подсобака"
	}
	if equals == 0 {
		return "Уникальный снежок", "Энигма"
	}
	if b12_15 == 6 {
		return "Золотая середина", "Балансировочное лезвие"
	}
	if b10_14 == 6 {
		return "Чудо середины дороги", "Средний приключенец"
	}
	if inn >= 16 && wis >= 16 {
		return "Мозговая Правда", "Интелектуарный гигант"
	}
	if cha == 18 && b12_15 == 5 {
		return "Харизматичный лидер", "Народный чемпион"
	}
	// the brawler - the unyielding warrior
	if str >= 16 && con >= 16 {
		return "Браулер", "Неелдующий воин"
	}
	if dex == 18 && str <= 6 {
		return "Подвижный якорь", "Быстрый страйкер"
	}
	if dex == 18 && b12_15 == 5 {
		return "Подвижный туз", "Быстрый страйкер со слабым пятном"
	}
	if p16 == 0 && m10 >= 3 {
		return "Жека без торговлей", "Борцующийся выживальщик"
	}
	if p18 == 1 && m10 == 5 {
		return "Непохоже что гирой", "Случайный чемпион"
	}
	if str == 18 && con == 18 {
		return "Скульптурный", "Недвижимый объект"
	}
	if str == 18 && con == 18 && dex <= 6 {
		return "Скульптурный спотыкач", "Неуклюжий гигант"
	}
	if inn >= 16 && dex >= 16 {
		return "Церебральный убийца", "Кунирующий киллер"
	}
	if inn >= 16 && wis >= 16 && cha <= 6 {
		return "Церебральный спотыкач", "Нехаризматичный гений"
	}
	if m8 == 6 {
		return "Неудачливая душа", "Проклятый путещественник"
	}
	if str <= 6 && m10 == 5 {
		return "Слабый боец", "Жалкий боксёр"
	}
	if dex <= 6 && m10 == 5 {
		return "Неуклюжий придурок", "Храбрый с пальцами"
	}
	if inn <= 6 && m10 == 5 {
		return "Тусклая лампочка", "Тупой сорвиголова"
	}
	if wis <= 6 && m10 == 5 {
		return "Наивный новичок", "Доверчивый игрок"
	}
	if cha <= 6 && m10 == 5 {
		return "Невдохновляющий лидер", "Непопулярный герой"
	}
	if con <= 6 && m10 == 5 {
		return "Хрупкий цветочек", "Деликатный дорогуша"
	}
	if m12 == 0 && m8 >= 3 {
		return "Всеобще ужасный", "Полный провал"
	}
	if m6 == 6 {
		return "Бездонный искатель приключений", "Жалкий исследователь"
	}
	if str >= 4 && m10 == 5 {
		return "Бессильный бродяга", "Бессильный боксёр"
	}
	if dex >= 4 && m10 == 5 {
		return "Неуклюжая катастрофа", "Каламутный скалолаз"
	}
	if inn >= 4 && m10 == 5 {
		return "Интелектуально вызванный", "Глупенький савант"
	}
	if wis >= 4 && m10 == 5 {
		return "Бездумный бродяга", "Наивный кочевник"
	}
	if cha >= 4 && m10 == 5 {
		return "Стрёмный неудачник", "Непопулярный недостигатор"
	}
	if con >= 4 && m10 == 5 {
		return "Хрупкая неудача", "Ломаный храбрец"
	}
	if dex <= 6 && str <= 6 {
		return "Нескоординированный придурок", "Неуклюжий сгусток"
	}
	if inn <= 6 && wis <= 6 {
		return "Дуплет тупизны", "Глупый собиратель"
	}
	if inn <= 6 && wis <= 6 && cha <= 6 {
		return "Триплет тупизны", "Непопулярный контрдостигатель"
	}
	if con <= 6 && str <= 6 && dex <= 6 && inn <= 6 {
		return "Четверной хрупок", "Мегаломающийся храбрый"
	}
	if m4 == 6 {
		return "Шестикратная угроза", "Полная катастрофа"
	}
	if m8 == 5 && m4 == 1 {
		return "Невезучая подсобака", "Неудачный контрдостигатель"
	}
	if m6 == 5 && m4 == 1 {
		return "Горестный негодяй", "Жалкий исследователь"
	}
	if m6 == 6 {
		return "Эксперт горячего замеса", "Каламутный крестодёр"
	}
	if m10 == 5 && str <= 4 {
		return "Кирпичная стена", "Неходящий объект"
	}
	if m10 == 5 && dex <= 4 {
		return "Неуклюжий кинг", "Неуклюжий чемпион"
	}
	if m10 == 5 && inn <= 4 {
		return "Мозговой слив", "Тускломуный дипломат"
	}
	if m10 == 5 && wis <= 4 {
		return "Космический кадет", "Глупый собиратель"
	}
	if m10 == 5 && cha <= 4 {
		return "Социальный изгой", "Непопулярный контрдостигатель"
	}
	if m10 == 5 && con <= 4 {
		return "Стеклянная пушка", "Мегаломающийся храбрый"
	}
	if m10 == 3 && m6 == 3 {
		return "Трио поездокрушения", "Каламутный крестодёр"
	}
	if con <= 6 && str <= 6 && dex <= 6 {
		return "Несвятая троица", "Раскоординированный монстр"
	}
	if inn <= 6 && wis <= 6 && cha <= 6 {
		return "Подрыв доверия к мозгу", "Непопулярный контрдостигатель"
	}
	if m4 == 6 {
		return "Пятикратная угроза себе", "Полная катастрофа"
	}
	if m8 == 5 && m4 == 1 {
		return "Пятикратная угроза себе", "Полная катастрофа"
	}
	if p16 >= 1 && m8 >= 1 {
		return "Высокий и низкий герой", "Непредсказуемая подсобака"
	}
	if p18 >= 1 && m6 >= 1 {
		return "Разбалансированный зверь", "Разбалансированная подсобака"
	}
	if inn == 18 && m10 == 5 {
		return "Нестабильный гений", "Безумный учоный"
	}
	if cha == 18 && m10 == 5 {
		return "Харизматичный шансовщик", "Серебрянноязыкий дъиявол"
	}
	if dex == 18 && m10 == 5 {
		return "Подвижный акробат", "Бысрый ударщик"
	}
	if str <= 4 && m12 == 5 {
		return "Бессильное чудо", "Дрыщеватый боец"
	}
	if dex <= 4 && m12 == 5 {
		return "Неуклюжая катастрофа", "Щас-что-то-будет"
	}
	if p16 == 3 && m6 == 1 {
		return "Высокий роллер проклятый", "неУдачливая легенда"
	}

	return "", ""
}
