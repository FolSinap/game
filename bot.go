package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

var buttons = []tgbotapi.KeyboardButton{
	tgbotapi.KeyboardButton{Text: "осмотреться"},
}

func initBot() {
	bot, err := tgbotapi.NewBotAPI("5321154532:AAFHkjPuu9y9c_chf0UmCNcz3KDvZpMTRXQ")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)
	updatesChan := make(chan []tgbotapi.Update)
	updateIdChan := make(chan int)
	go getUpdates(bot, updatesChan, updateIdChan)
	go respond(bot, updatesChan, updateIdChan)
	updateIdChan <- 43802767
}

func processMessages(player *Player, chat *tgbotapi.Chat,bot *tgbotapi.BotAPI, die chan bool) {
	var message tgbotapi.MessageConfig
	var msg string
	for {
		select {
		case msg = <- player.output:
			message = tgbotapi.NewMessage(chat.ID, msg)
			message.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)
			bot.Send(message)
			log.Println("send: ", message.Text)
			log.Println("to: ", player.name)
		case <- die:
			removePlayer(player)
			return
		}
	}
}

func timer(player *Player, die chan bool) {
	var msg string
	for {
		select {
		case msg = <-player.output:
			player.output <- msg
		case <-time.After(time.Minute * 15):
			player.HandleOutput("Вы удалены из игры. Причина - афк 15 минут")
			die <- true
			return
		}
	}
}

func respond(bot *tgbotapi.BotAPI, c chan []tgbotapi.Update, lastIdChan chan int) {
	var updates []tgbotapi.Update
	var player *Player
	var lastId int
	for {
		updates = <- c
		for _, update := range updates {
			player = findPlayer(update.Message.From.FirstName)
			if player == nil {
				player = NewPlayer(update.Message.From.FirstName)
				addPlayer(player)
				die := make(chan bool)
				go processMessages(player, update.Message.Chat, bot, die)
				go timer(player, die)
			}

			log.Println("received: ", update.Message.Text)
			log.Println("from: ", player.name)

			player.HandleInput(update.Message.Text)
			lastId = update.UpdateID + 1
		}
		lastIdChan <- lastId
	}
}

func getUpdates(bot *tgbotapi.BotAPI, c chan []tgbotapi.Update, lastId chan int) {
	var conf tgbotapi.UpdateConfig
	conf.AllowedUpdates = []string{"message"}
	for {
		conf.Offset = <- lastId
		updates, _ := bot.GetUpdates(conf)
		c <- updates
	}
}
