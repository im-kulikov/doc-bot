package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/im-kulikov/doc-bot/internal"
	"github.com/rs/zerolog/log"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	self *tb.Bot

	helpMsg  string
	helloMsg string

	search *template.Template
}

func (b *Bot) sendMessage(sender *tb.User, message string) {
	if _, err := b.self.Send(sender, message, tb.ModeMarkdown); err != nil {
		log.Error().Msgf("Error: %v", err)
	}
}

func (b *Bot) handler(m *tb.Message) {
	var (
		err     error
		query   string
		message string
		results *internal.SearchResults
		buffer  = bytes.NewBuffer(nil)
		items   = internal.SearchRe.FindAllStringSubmatch(m.Text, -1)
	)

	if len(items) == 0 {
		b.help(m)
		return
	}

	message = fmt.Sprintf(
		internal.SearchAnswerTpl,
		items[0][1],
	)

	log.Info().Msgf("search for query '%s'", items[0][1])

	b.sendMessage(m.Sender, message)

	query = internal.QueryRe.ReplaceAllString(items[0][1], "+")

	if results, err = internal.SearchInDocuments(query); err != nil {
		log.Error().Msgf("search results error: %v", err)

		b.sendMessage(m.Sender, "Ошибка запроса")
		return
	}

	if err = b.search.Execute(buffer, results); err != nil {
		log.Error().Msgf("template error: %v", err)
		b.sendMessage(m.Sender, "Ошибка ответа")
		return
	}

	b.sendMessage(m.Sender, buffer.String())
}

func (b *Bot) hello(m *tb.Message) {
	var message = fmt.Sprintf(internal.HelloTpl, m.Sender.FirstName)
	b.sendMessage(m.Sender, message)
}

func (b *Bot) help(m *tb.Message) {
	b.sendMessage(m.Sender, internal.HelpTpl)
}

var (
	BuildTime string
	Version   string
)

func main() {
	var (
		err      error
		bot      *tb.Bot
		instance *Bot
	)

	log.Output(os.Stdout)
	log.Info().Msgf(internal.InfoTpl, Version, BuildTime)

	if bot, err = tb.NewBot(tb.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}); err != nil {
		log.Panic().Msgf("Error on connect: %v", err)
	}

	instance = new(Bot)
	instance.self = bot

	if instance.search, err = internal.SearchTemplate(); err != nil {
		log.Panic().Msgf("Error on template: %v", err)
	}

	bot.Handle("/start", instance.hello)
	bot.Handle("/help", instance.help)
	bot.Handle(tb.OnText, instance.handler)

	bot.Start()
}
