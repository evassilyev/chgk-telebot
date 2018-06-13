# chgk-telebot
Telegram bot for "What? Where? When?" gaming.

Телеграм бот для игры в "Что? Где? Когда?".

### Description

The bot lets you to play "What? Where? When?" in the group Telegram chat in a text mode. 
It downloads the set of the questions from [the questions database]( https://db.chgk.info) and shows you them one by one hiding answers until your request or a timer finish. 
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

#####List of commands

* __/?__ __/HELP__ - show the help
* __/get_packet N__ __/получить_пакет N__ - load the packet of N questions
* __/start__ __/начать__ - start the game
* __/next__ __/след__ - go to the next question
* __/prev__ __/пред__ - go to the previous question
* __/question N__ __/вопрос N__ - go to the quesiton number N
* __/answer__ __/ответ__ - show the answer
* __/info__ - show the additional question's information (author, sources etc.)
* __/set_timer__ - set the timer in minutes (fractional values available)
