package bot

import (
	"context"

	"github.com/chapsuk/worker"
	"github.com/im-kulikov/helium"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type (
	app struct {
		log     *zap.Logger
		workers *worker.Group
	}

	appParams struct {
		dig.In

		Logger  *zap.Logger
		Workers *worker.Group
	}
)

func newApp(p appParams) helium.App {
	return &app{
		log:     p.Logger,
		workers: p.Workers,
	}
}

func (a *app) Run(ctx context.Context) error {
	a.log.Info("start workers...")
	a.workers.Run()

	a.log.Info("service started successful...")
	<-ctx.Done()

	a.log.Info("stop workers...")
	a.workers.Stop()

	a.log.Info("service stop successful...")
	return nil
}
