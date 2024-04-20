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
)

func CommandHandler(c tele.Context, serv *servitor.Servitor, flags *silly, comfig *config.AppConfig, logger *slog.Logger) error {
	logger.Info("incomingtext message",
		slog.Int64("chatID:", c.Chat().ID),
		slog.String("message:", c.Message().Text))
	msgText := strings.Split(c.Message().Text, " ")

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
	message += "Сила: " + dice3of4() + "\n"
	message += "Ловкость: " + dice3of4() + "\n"
	message += "Телосложение: " + dice3of4() + "\n"
	message += "Интеллект: " + dice3of4() + "\n"
	message += "Мудрость: " + dice3of4() + "\n"
	message += "Харя: " + dice3of4() + "\n"
	return c.Send(message)
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
