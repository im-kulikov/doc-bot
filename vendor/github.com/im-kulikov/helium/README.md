# Helium

![Codecov](https://img.shields.io/codecov/c/github/im-kulikov/helium.svg?style=flat-square)
[![Build Status](https://travis-ci.com/im-kulikov/helium.svg?branch=master)](https://travis-ci.com/im-kulikov/helium)
[![Report](https://goreportcard.com/badge/github.com/im-kulikov/helium)](https://goreportcard.com/report/github.com/im-kulikov/helium)
[![GitHub release](https://img.shields.io/github/release/im-kulikov/helium.svg)](https://github.com/im-kulikov/helium)
![GitHub](https://img.shields.io/github/license/im-kulikov/helium.svg?style=popout)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fim-kulikov%2Fhelium.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fim-kulikov%2Fhelium?ref=badge_shield)

<img src="./.github/helium.jpg" width="350" alt="logo">

# Documentation

* [Why Helium](#why-helium)
* [About Helium and modules](#about-helium-and-modules)
    * [Logger](#logger-module)
    * [NATS](#nats-module)
    * [PostgreSQL](#postgresql-module)
    * [Redis](#redis-module)
    * [Settings](#settings-module)
    * [Web](#web-module)
* [Examples of code](https://github.com/go-helium/demos)
* [Project Examples](#project-examples)
* [Example](#example)
* [Supported Go versions](#supported-go-versions)
* [Contribute](#contribute)
* [Credits](#credits)
* [License](#license)

## Why Helium

When building a modern application or prototype proof of concept, the last thing you want to worry about is boilerplate code or pass dependencies.
All this is a routine that each of us faces.
Helium lets you get rid of this routine and focus on your code.
Helium provides to you Convention over configuration.

The modular structure of helium consists of the following concepts
- **Module** is a set of providers
- **Provider** is what you put in the DI container. It teaches the container how to build values of one or more types and expresses their dependencies. It consists of two components of the **constructor** and **options**.
- **Constructor** is a function that accepts zero or more parameters and returns one or more results. The function may optionally return an error to indicate that it failed to build the value. This function will be treated as the constructor for all the types it returns. This function will be called AT MOST ONCE when a type produced by it, or a type that consumes this function's output, is requested via Invoke. If the same types are requested multiple times, the previously produced value will be reused. In addition to accepting constructors that accept dependencies as separate arguments and produce results as separate return values, Provide also accepts constructors that specify dependencies as dig.In structs and/or specify results as dig.Out structs.
- **Options** modifies the default behavior of dig.Provide

## About Helium and modules

*Helium is a small, simple, modular constructor with some pre-built components for your convenience.*
 
It contains the following components for rapid prototyping of your projects:
- Grace - [context](https://golang.org/pkg/context/) that helps you gracefully shutdown your application
- Logger - [zap](https://go.uber.org/zap) is blazing fast, structured, leveled logging in Go
- DI - based on [DIG](https://go.uber.org/dig). A reflection based dependency injection toolkit for Go.
- Module - set of tools for working with the DI component
- [NATS](https://github.com/go-helium/nats) - [nats](https://github.com/nats-io/go-nats) and [NSS](https://github.com/nats-io/nats-streaming-server), client for the cloud native messaging system
- [PostgreSQL](https://github.com/go-helium/postgres) - client module for [ORM](https://github.com/go-pg/pg) with focus on PostgreSQL features and performance
- [Redis](https://github.com/go-helium/redis) - module for type-safe [Redis](https://github.com/go-redis/redis) client for Golang  
- Settings - based on [Viper](https://github.com/spf13/viper). A complete configuration solution for Go applications including 12-Factor apps. It is designed to work within an application, and can handle all types of configuration needs and formats
- Web - [see more](#web-module)
- Workers - are tools to run goroutines and do arbitrary work on a schedule along with a mechanism to safely stop each one. Based on [chapsuk/worker](https://github.com/chapsuk/worker)

### Logger module

Module provides you with the following things:
- `*zap.Logger` instance of [Logger](https://godoc.org/go.uber.org/zap#Logger)

A Logger provides fast, leveled, structured logging. All methods are safe for concurrent use.
The Logger is designed for contexts in which every microsecond and every allocation matters, so its API intentionally favors performance and type safety over brevity. For most applications, the SugaredLogger strikes a better balance between performance and ergonomics.
- `*zap.SugaredLogger` instance of [SugaredLogger](https://godoc.org/go.uber.org/zap#SugaredLogger)

Unlike the Logger, the SugaredLogger doesn't insist on structured logging. For each log level, it exposes three methods: one for loosely-typed structured logging, one for println-style formatting, and one for printf-style formatting. For example, SugaredLoggers can produce InfoLevel output with Infow ("info with" structured context), Info, or Infof.
A SugaredLogger wraps the base Logger functionality in a slower, but less verbose, API. Any Logger can be converted to a SugaredLogger with its Sugar method.
- `logger.StdLogger` provides simple interface that pass calls to **zap.SugaredLogger**

```
StdLogger interface {
    Fatal(v ...interface{})
    Fatalf(format string, v ...interface{})
    Print(v ...interface{})
    Printf(format string, v ...interface{})
}
```

Logger levels:
- DebugLevel logs are typically voluminous, and are usually disabled in production
- InfoLevel is the default logging priority
- WarnLevel logs are more important than Info, but don't need individual human review
- ErrorLevel logs are high-priority. If an application is running smoothly, it shouldn't generate any error-level logs
- DPanicLevel logs are particularly important errors. In development the logger panics after writing the message
- PanicLevel logs a message, then panics.
- FatalLevel logs a message, then calls os.Exit(1)

Logger formats:
- console:
```
2019-02-19T20:22:28.239+0300	info	web/servers.go:80	Create metrics http server, bind address: :8090	{"app_name": "Test", "app_version": "dev"}
2019-02-19T20:22:28.239+0300	info	web/servers.go:80	Create pprof http server, bind address: :6060	{"app_name": "Test", "app_version": "dev"}
2019-02-19T20:22:28.239+0300	info	app.go:26	init	{"app_name": "Test", "app_version": "dev"}
2019-02-19T20:22:28.239+0300	info	app.go:28	run workers	{"app_name": "Test", "app_version": "dev"}
2019-02-19T20:22:28.239+0300	info	app.go:31	run web-servers	{"app_name": "Test", "app_version": "dev"}
```
- json:
```
{"level":"info","msg":"Failed to fetch URL.","url":"http://example.com","attempt":3,"backoff":"1s"}
{"level":"info","msg":"Failed to fetch URL: http://example.com"}
{"level":"info","msg":"Failed to fetch URL.","url":"http://example.com","attempt":3,"backoff":"1s"}
```

Configuration for logger
- yaml example
```yaml
debug: true

logger:
    format: console
    level: info
    trace_level: fatal
    no_disclaimer: false
    color: true
    full_caller: true
    sampling:
      initial: 100
      thereafter: 100
```
- env example
```
DEBUG=true
LOGGER_NO_DISCLAIMER=true
LOGGER_COLOR=true
LOGGER_FULL_CALLER=true
LOGGER_FORMAT=console
LOGGER_LEVEL=info
LOGGER_TRACE_LEVEL=fatal
LOGGER_SAMPLING_INITIAL=100
LOGGER_SAMPLING_THEREAFTER=100
```

- `debug` - with this option you can enable `zap.DevelopmentConfig()`
- `logger.no_disclaimer` - with this option, you can disable `app_name` and `app_version` for any reason (not recommended in production)
- `logger.trace_level` - configures the Logger to record a stack trace for all messages at or above a given level
- `logger.color` - serializes a Level to an all-caps string and adds color
- `logger.full_caller` - serializes a caller in /full/path/to/package/file:line format
- `logger.sampling.initial` and `logger.sampling.thereafter` to setup [logger sampling](https://godoc.org/go.uber.org/zap#SamplingConfig). SamplingConfig sets a sampling strategy for the logger. Sampling caps the global CPU and I/O load that logging puts on your process while attempting to preserve a representative subset of your logs. Values configured here are per-second. See [zapcore.NewSampler](https://godoc.org/go.uber.org/zap/zapcore#NewSampler) for details.

## NATS Module

[Module](https://github.com/go-helium/nats) provides you with the following things:
- [`*nats.Conn`](https://godoc.org/github.com/nats-io/go-nats#Conn) represents a bare connection to a nats-server. It can send and receive []byte payloads
- [`stan.Conn`](https://godoc.org/github.com/nats-io/go-nats-streaming#Conn) represents a connection to the NATS Streaming subsystem. It can Publish and Subscribe to messages within the NATS Streaming cluster.

Configuration:
- yaml example
```yaml
nats:
  url: nats://<host>:<port>
  cluster_id: string
  client_id: string
  servers: [...server slice...]
  no_randomize: bool
  name: string
  verbose: bool
  pedantic: bool
  secure: bool
  allow_reconnect: bool
  max_reconnect: int
  reconnect_wait: duration
  timeout: duration
  flusher_timeout: duration
  ping_interval: duration
  max_pings_out: int
  reconnect_buf_size: int
  sub_chan_len: int
  user: string
  password: string
  token: string
```
- env example
```
NATS_URL=nats://<host>:<port>
NATS_CLUSTER_ID=string
NATS_CLIENT_ID=string
NATS_SERVERS=[...server slice...]
NATS_NO_RANDOMIZE=bool
NATS_NAME=string
NATS_VERBOSE=bool
NATS_PEDANTIC=bool
NATS_SECURE=bool
NATS_ALLOW_RECONNECT=bool
NATS_MAX_RECONNECT=int
NATS_RECONNECT_WAIT=duration
NATS_TIMEOUT=duration
NATS_FLUSHER_TIMEOUT=duration
NATS_PING_INTERVAL=duration
NATS_MAX_PINGS_OUT=int
NATS_RECONNECT_BUF_SIZE=int
NATS_SUB_CHAN_LEN=int
NATS_USER=string
NATS_PASSWORD=string
NATS_TOKEN=string
```

## PostgreSQL Module

[Module](https://github.com/go-helium/postgres) provides you connection to PostgreSQL server
- `*pg.DB` is a database handle representing a pool of zero or more underlying connections. It's safe for concurrent use by multiple goroutines

Configuration:
- yaml example
```yaml
posgres:
    address: string
    username: string
    password: string
    database: string
    debug: bool
    pool_size: int
```
- env example
```
POSTGRES_ADDRESS=string
POSTGRES_USERNAME=string
POSTGRES_PASSWORD=string
POSTGRES_DATABASE=string
POSTGRES_DEBUG=bool
POSTGRES_POOL_SIZE=int
```

## Redis Module

[Module](https://github.com/go-helium/redis) provides you connection to Redis server
- `*redis.Client` is a Redis client representing a pool of zero or more underlying connections. It's safe for concurrent use by multiple goroutines

Configuration:
- yaml example
```yaml
redis:
  address: string
  password: string
  db: int
  max_retries: int
  min_retry_backoff: duration
  max_retry_backoff: duration
  dial_timeout: duration
  read_timeout: duration
  write_timeout: duration
  pool_size: int
  pool_timeout: duration
  idle_timeout: duration
  idle_check_frequency: duration
```
- env example
```
REDIS_ADDRESS=string
REDIS_PASSWORD=string
REDIS_DB=int
REDIS_MAX_RETRIES=int
REDIS_MIN_RETRY_BACKOFF=duration
REDIS_MAX_RETRY_BACKOFF=duration
REDIS_DIAL_TIMEOUT=duration
REDIS_READ_TIMEOUT=duration
REDIS_WRITE_TIMEOUT=duration
REDIS_POOL_SIZE=int
REDIS_POOL_TIMEOUT=duration
REDIS_IDLE_TIMEOUT=duration
REDIS_IDLE_CHECK_FREQUENCY=duration
```

## Settings module

Module provides you [`*viper.Viper`](https://godoc.org/github.com/spf13/viper#Viper), a complete configuration solution for Go applications including 12-Factor apps. It is designed to work within an application, and can handle all types of configuration needs and formats.

Viper is a prioritized configuration registry. It maintains a set of configuration sources, fetches values to populate those, and provides them according to the source's priority. The priority of the sources is the following:
1. overrides
2. flags
3. env. variables
4. config file
5. key/value store 6. defaults

## Web Module

- bind - simple replacement for echo.Binder
- validate - simple replacement for echo.Validate
- logger - provides echo.Logger that pass calls to **zap.Logger**

- `ServersModule` puts into container [multi-server](https://github.com/chapsuk/mserv):
    - [pprof](https://golang.org/pkg/net/http/pprof/) endpoint
    - [metrics](https://github.com/prometheus/client_golang) enpoint (by Prometheus)
    - You can pass `profile_handler` and/or `metric_handler`, that will be embedded into common handler,
     and will be available to call them
    - **api** endpoint by passing http.Handler from DI
- [`echo.Module`](https://github.com/go-helium/echo) boilerplate that preconfigures echo.Engine for you
    - with custom Binder / Logger / Validator / ErrorHandler

Configuration:
- yaml example
```yaml
pprof:
  address: :6060
  shutdown_timeout: 10s

metrics:
  address: :8090
  shutdown_timeout: 10s

api:
  address: :8080
  shutdown_timeout: 10s
```
- env example
```
PPROF_ADDRESS=string
PPROF_SHUTDOWN_TIMEOUT=duration
METRICS_ADDRESS=string
METRICS_SHUTDOWN_TIMEOUT=duration
API_ADDRESS=string
API_SHUTDOWN_TIMEOUT=duration
```

**Possible options**:
- `address` - (string) host and port
- `disabled` - (bool) to disable server
- `read_timeout` - (duration) is the maximum duration for reading the entire request, including the body
- `read_header_timeout` - (duration) is the amount of time allowed to read request headers
- `write_timeout` - (duration) is the maximum duration before timing out writes of the response
- `idle_timeout` - (duration) is the maximum amount of time to wait for the next request when keep-alives are enabled
- `max_header_bytes` - (int) controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line
- `shutdown_timeout` - (duration) context timeout for stopping server 

## Workers Module

Simple abstraction for control background jobs.

[Worker module](https://github.com/chapsuk/worker) adding the abstraction layer around background jobs, allows make a job periodically, observe execution time and to control concurrent execution.

Group of workers allows to control jobs start time and wait until all runned workers finished when we need stop all jobs.


Workers config example:
```yaml
workers:
  ## example: ##
  ## job_name:
  ##   disabled: false   # to disable worker
  ##   immediately: true # run worker when start application
  ##   ticker: 30s       # run job every 30s
  ##   timer: 30s        # run job every 30s and reset timer
  ##   cron: * * * * *   # run job by crontab specs, e.g. "* * * * *"
  ##   lock:             # use lock by redis
  ##     key: myKey
  ##     ttl: 150s
  ##     retry:          # retry get lock to run job
  ##       count:   5
  ##       timeout: 30s
```

**Features**
- Scheduling, use one from existing worker.By* schedule functions. Supporting cron schedule spec format by robfig/cron parser.
- Control concurrent execution around multiple instances by worker.With* lock functions. Supporting redis locks by go-redis/redis and bsm/redis-lock pkgs.
- Observe a job execution time duration with worker.SetObserever. Friendly for prometheus/client_golang package.
- Graceful stop, wait until all running jobs was completed.

## Project Examples

- [Atlant.io Test Task](https://github.com/im-kulikov/atlantio-task) 
- [Simplinic Test Task](https://github.com/im-kulikov/simplinic-task) 
- [Golang documentation Telegram Bot](https://github.com/im-kulikov/doc-bot)
- [Potter](https://github.com/im-kulikov/potter) is simple fixture based API service
- [Image processing service](https://github.com/archaron/secimage)

## Example

**config.yml**
```yaml
api:
  address: :8080
  debug: true
  shutdown_timeout: 10s

logger:
  level: debug
  format: console
```

**main.go**
```go
package main

import (
	"context"
	"net/http"

	"github.com/chapsuk/mserv"
	"github.com/im-kulikov/helium"
	"github.com/im-kulikov/helium/grace"
	"github.com/im-kulikov/helium/logger"
	"github.com/im-kulikov/helium/module"
	"github.com/im-kulikov/helium/settings"
	"github.com/im-kulikov/helium/web"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func main() {
	h, err := helium.New(&helium.Settings{
		Name:         "demo3",
		Prefix:       "DM3",
		File:         "config.yml",
		BuildVersion: "dev",
	}, module.Module{
		{Constructor: handler},
	}.Append(
		grace.Module,
		settings.Module,
		web.ServersModule,
		logger.Module,
	))
	err = dig.RootCause(err)
	helium.Catch(err)
	err = h.Invoke(runner)
	err = dig.RootCause(err)
	helium.Catch(err)
}

func handler() http.Handler {
	h := http.NewServeMux()
	h.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	return h
}

func runner(s mserv.Server, l *zap.Logger, ctx context.Context) {
	l.Info("run servers")
	s.Start()

	l.Info("application started")
	<-ctx.Done()

	l.Info("stop servers")
	s.Stop()

	l.Info("application stopped")
}
```

## Supported Go versions

Helium is available as a [Go module](https://github.com/golang/go/wiki/Modules).
- 1.11+

## Contribute

**Use issues for everything**

- For a small change, just send a PR.
- For bigger changes open an issue for discussion before sending a PR.
- PR should have:
  - Test case
  - Documentation
  - Example (If it makes sense)
- You can also contribute by:
  - Reporting issues
  - Suggesting new features or enhancements
  - Improve/fix documentation

## Credits

- [Evgeniy Kulikov](https://github.com/im-kulikov) - Author
- [Alexander Tischenko](https://github.com/archaron) - Consultant
- [Contributors](https://github.com/im-kulikov/helium/graphs/contributors)

## License

[MIT](LICENSE)

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fim-kulikov%2Fhelium.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fim-kulikov%2Fhelium?ref=badge_large)
