package bigcat

const msgHelp = `Покасьто я могу так:

/progress - как много работы надо мной (и не только) было сегодня
/changelog - последние обновления миня
/meme - мем с непонятливым котиком
/hold - мем с мощным сдерживанием
/joke - помочь с разгадкой шутки дня от бати
/getopening - получить файл с опенингами, доступно, навельное, не всем
/setopening - я знаю где в гитхабе лежит файл с опенингами, мозет залью себе
/mawlist - все наши варики
/animemaw - случайный опенинг
/grobmaw - случайный гроб
/maw - случайный мав из общака
/mawadd - добавить мав в общак описание + ссылка
/timewo - cписьёк инцедентив
/newtimewo - добавить инцидентъ
/weather - температуря в городи
/wday - прогнозь для городя
/vectorgame N - сыграть в викторину с указанным номером
/vectoraddtype - добавить викторину
/vectortypes - посмотреть номера викторин
/vectoraddq - добавить вопрос в викторину (через запятую) номер_викторины, ссылка_на_картинку, вопрос, ответ1,...,ответN
`

//settask - поставь себе задачу и я её для тебя запомню 👨‍❤️‍👨
//gettask - узнай что ты задумол но так и не зделол 🫄
//all - посмотреть что по задачам у всей бригады

const msgHello = "Здарова, заебал!\n\n" + msgHelp

const changelog = "добавление в базу викторин\n"

const (
	msgUnknownCommand = "Незнаю такую командоцку 🤔🥴"
	msgNoSavedPages   = "У васъ нету страницек 😅"
	msgSaved          = "Схорониля ❤️"
	msgTaskSaved      = "Таргет дезигнейтед"
	msgAlreadyExists  = "Узе есь такая ссилоцка 😅"
)
