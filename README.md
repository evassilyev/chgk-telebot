# chgk-telebot
Telegram bot for "What? Where? When?" gaming.

Телеграм бот для игры в "Что? Где? Когда?".

### Description

The bot lets you to play "What? Where? When?" in the group Telegram chat in a text mode. 
It downloads the set of the questions from [the questions database](https://db.chgk.info) and shows you them one by one hiding answers until your request or a timer finish. 
You can configure the timer duration and the number of the questions in the set. All management and configuration provides by text commands in the Telegram chat.

### Installation

1. [Create your telegram bot](https://core.telegram.org/bots#3-how-do-i-create-a-bot) or use already created if you have it. Create the group chat for gaming and add your bot there. 
   You need to know the token of your bot and the group chat ID. 
   You can get the token by @BotFather. The chat ID you can get by using follow link: https:\//api.telegram.org/bot*YourBOTToken*/getUpdates. In the JSON responce you will find the ID of the group chat. Usually it is a negative integer number.
1. Prepare the configuration file. Download _config.example.json_, rename it to _config.json_. Then, write your bot's token (in the double quotes) and the chat ID (just like a number) there.
1. Download the executable file from the _Builds/v1.0/_ folder and place it in the same directory with the configuration file _config.json_.
1. Just run it as a console application. You can set the path to configuration file using the _-config_ console argument.

### Playing

You can get the list of all available commands by texting **/?** or **/HELP** in the chat.
Before the play you should load the set of questions by the comand **/get_packet** then text **/start** for starting the game.

##### List of commands

* __/?__ __/HELP__ - show the help
* __/get_packet N__ __/получить_пакет N__ - load the packet of N questions
* __/start__ __/начать__ - start the game
* __/next__ __/след__ - go to the next question
* __/prev__ __/пред__ - go to the previous question
* __/question N__ __/вопрос N__ - go to the quesiton number N
* __/answer__ __/ответ__ - show the answer
* __/info__ - show the additional question's information (author, sources etc.)
* __/set_timer__ - set the timer in minutes (fractional values available)

### Notes 

* In the first questions set loading you have to mandatory set the number of questions. In the next loading questions number will be the same as it was during the previous one by default. 
* Minimal timer value is 0.25 (15 seconds).
* The application loads the random set of questions for the all store period. All the questions in the set has a random complexity and just types "What? Where? When?" and "Brain-ring".

#Бот для игры в "Что? Где? Когда?"
### Описание

Телеграм бот позволяет вам играть в _"Что? Где? Когда?"_ прямо в групповом чате Телеграм. 
Для игры используются вопросы, которые загружаются напрямую из [базы вопросов](https://db.chgk.info) ЧГК.
Процесс игры представляет из себя отправку ботом вопроса с последующей отправкой ответа по таймеру либо по запросу пользователя.
Таймер и количество вопросов в пакете - настраиваемые параметры. Всё взаимодействие с ботом в процессе игры и настройки осуществляется посредством текстовых команд в чате.

### Установка

1. Для начала вам необходимо [создать Телеграм бота](https://core.telegram.org/bots#3-how-do-i-create-a-bot) или использовать уже имеющийся. Создайте групповой чат в Телеграм и добавьте туда своего бота.
   Вам необходимо узнать токен вашего бота и идентификатор группового чата. Токен бота можно получить при помощи служебного бота _@BotFather_ (Им вы уже пользовались при создании своего бота).
   Идентификатор группового чата вы можете получить, пройдя по ссылке https:\//api.telegram.org/bot*YourBOTToken*/getUpdates в которой вам необходимо предварительно заменить YourBOTToken на токен вашего бота.
   Ссылка вернет вам структуру в формате JSON, где вы сможете найти идентификатор группового чата, в который вы добавили бота. Обычно это целое отрицательное число.
1. Далее, подготовьте конфигурационный файл. Скачайте или скорпируйте из репозитория файл _config.example.json_, переименуйте его в _config.json_, затем, впишите в него в двойных кавычках токен вашего бота и идентификатор группового чата.
1. Загрузите готовый исполняемый файл из папки _Builds/v1.0/_ и разместите его в одной папке с конфигурационным файлом _config.json_.
1. Запустите исполняемый файл как консольное приложение. Так же, вы можете указать путь к конфигурационному файлу с помощью аргумента командной строки _-config_. 

### Процесс игры

Загрузите пакет написав в групповом чате команду **/get_packet** или **/получить_пакет**. После чего начните игру с помощью комады **/start** или **/начать**.

##### Список доступных команд

* __/?__ __/HELP__ - Показать справку по командам
* __/get_packet N__ __/получить_пакет N__ - загрузить пакет из N вопросов
* __/start__ __/начать__ - начать игру
* __/next__ __/след__ - следующий вопрос
* __/prev__ __/пред__ - предыдущий вопрос
* __/question N__ __/вопрос N__ - перейти к вопросу под номером N
* __/answer__ __/ответ__ - показать ответ
* __/info__ - показать информацию о вопросе (автор, источники и т.д.)
* __/set_timer__ - установить таймер в минутах (возможно установить дробные значения)

### Замечания

* При первой загрузке пакета после запуска приложения необходимо обязательно указать число вопросов в пакете, при последующих загрузках по умолчанию будет загружаться то число вопросов как и в предыдущей загрузке пакета.
* Минимально допустимое значение таймера 0.25 минут (15 секунд). После первой установки таймера полностью отключить его невозможно. Сымитиовать отключение таймера можно установив его на большой период времени.
* Из базы загружаются только случайные пакеты состоящие из вопросов за весь хранимый период (c 1990 года) произвольной сложности. В пакет включаются вопросы только типов "Что? Где? Когда?" и "Брейн-ринг".
