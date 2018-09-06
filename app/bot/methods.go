package bot

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/Syfaro/telegram-bot-api"
	"go.uber.org/zap"
)

func (s *Service) Job(ctx context.Context) {
	wg := new(sync.WaitGroup)
	wg.Add(s.workers)

	for i := 0; i < s.workers; i++ {
		go func(num int) {
			s.worker(ctx, wg, num)
		}(i)
	}

	<-ctx.Done()
	wg.Wait()
}

func (s *Service) worker(ctx context.Context, wg *sync.WaitGroup, num int) {
	log := s.log.With(
		zap.Int("worker", num))

loop:
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case update := <-s.updates:
			if update.Message == nil || update.Message.Chat == nil {
				continue loop
			}

			// --- Process message --- //

			msg := update.Message

			log := log.With(
				zap.String("from", msg.From.String()),
				zap.Int64("chat_id", msg.Chat.ID),
				zap.String("message", msg.Text))

			switch {
			case msg.IsCommand():
				// commands:
				switch msg.Text {
				case "/start":
					txt := fmt.Sprintf(HelloTpl, msg.From.FirstName)
					s.Send(log, update, txt)
				case "/help":
					txt := Markdown(HelpTpl)
					s.Send(log, update, txt)
				default:
					s.Send(log, update, fmt.Sprintf("%s: unknown command", msg.Text))
				}
			default:
				items := SearchRe.FindAllStringSubmatch(msg.Text, -1)
				if len(items) < 1 || len(items[0]) < 2 {
					txt := Markdown(HelpTpl)
					s.Send(log, update, txt)
					log.Info("bad query")
					continue loop
				}

				message := fmt.Sprintf(
					SearchAnswerTpl,
					items[0][1],
				)

				log.Info("search for query",
					zap.String("query", items[0][1]))

				s.Send(log, update, message)

				var (
					query   = QueryRe.ReplaceAllString(items[0][1], "+")
					buffer  = new(bytes.Buffer)
					results *SearchResults
					err     error
				)

				if results, err = SearchInDocuments(query); err != nil {
					log.Error("search results error",
						zap.Error(err))

					s.Send(log, update, "⛔️ Ошибка запроса")
					continue loop
				}

				if err = s.search.Execute(buffer, results); err != nil {
					log.Error("template error",
						zap.Error(err))
					s.Send(log, update, "⛔️ Ошибка ответа")
					return
				}

				s.Send(log, update, Markdown(buffer.String()))
			}
		}
	}
}

func (s *Service) Send(log *zap.Logger, update tgbotapi.Update, msg interface{}) {
	var content tgbotapi.Chattable
	switch data := msg.(type) {
	case string:
		content = tgbotapi.NewMessage(update.Message.Chat.ID, data)
	case Sticker:
		content = tgbotapi.NewStickerShare(update.Message.Chat.ID, string(data))
	case Markdown:
		txt := tgbotapi.NewMessage(update.Message.Chat.ID, "\n"+string(data)+"\n")
		txt.ParseMode = tgbotapi.ModeMarkdown
		content = txt
	}

	if _, err := s.bot.Send(content); err != nil {
		log.Error("can't send message",
			zap.Error(err))
	}
}
