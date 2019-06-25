package bot

import (
	"github.com/im-kulikov/doc-bot/modules/query"
	"github.com/im-kulikov/helium/grace"
	"github.com/im-kulikov/helium/logger"
	"github.com/im-kulikov/helium/module"
	"github.com/im-kulikov/helium/settings"
	"github.com/im-kulikov/helium/workers"
)

var Module = module.Module{
	{Constructor: newApp},
	{Constructor: newJobs},
}.Append(
	grace.Module,    // graceful context
	settings.Module, // settings
	logger.Module,   // logger
	workers.Module,  // workers
	query.Module,    // telegram bot
)
