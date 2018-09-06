package bot

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"text/template"

	"github.com/Syfaro/telegram-bot-api"
	"github.com/im-kulikov/helium/module"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
)

type (
	Service struct {
		bot     *tgbotapi.BotAPI
		log     *zap.Logger
		workers int
		search  *template.Template
		updates tgbotapi.UpdatesChannel
	}

	Config struct {
		Token   string
		Client  *http.Client
		Workers int
	}

	params struct {
		dig.In

		Config *Config
		Logger *zap.Logger
	}

	Sticker  string
	Markdown string

	StringsField [][]string
)

var Module = module.Module{
	{Constructor: newConfig},
	{Constructor: newService},
}

func newClient(socks5 string) (*http.Client, error) {
	client := &http.Client{}
	if len(socks5) > 0 {
		proxyURL, err := url.Parse(socks5)
		if err != nil {
			return nil, errors.WithMessage(err, socks5)
		}

		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			return nil, errors.WithMessage(err, socks5)
		}

		client.Transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}
	}

	return client, nil
}

func newConfig(v *viper.Viper) (*Config, error) {
	var (
		err    error
		cli    *http.Client
		token  string
		socks5 string
		count  int
	)

	socks5 = v.GetString("telegram.proxy")
	if cli, err = newClient(socks5); err != nil {
		return nil, err
	}

	if token = v.GetString("telegram.token"); token == "" {
		return nil, errors.New("`telegram.token` can't be empty")
	}

	if count = v.GetInt("telegram.workers_count"); count == 0 {
		count = 10
	}

	return &Config{
		Token:   token,
		Client:  cli,
		Workers: count,
	}, nil
}

func newService(p params) (*Service, error) {
	bot, err := tgbotapi.NewBotAPIWithClient(p.Config.Token, p.Config.Client)
	if err != nil {
		return nil, err
	}

	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Timeout: 60,
	})
	if err != nil {
		return nil, err
	}

	tpl, err := SearchTemplate()
	if err != nil {
		return nil, err
	}

	return &Service{
		bot:     bot,
		workers: p.Config.Workers,
		updates: updates,
		search:  tpl,
		log: p.Logger.With(
			zap.String("caller", "bot")),
	}, nil
}
