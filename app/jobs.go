package app

import (
	"github.com/chapsuk/worker"
	"github.com/im-kulikov/doc-bot/app/bot"
	"go.uber.org/dig"
)

type jobsParams struct {
	dig.In

	// workers:
	Bot *bot.Service
}

// Application jobs:
func newJobs(p jobsParams) map[string]worker.Job {
	return map[string]worker.Job{
		"bot": p.Bot.Job,
	}
}
