bobo
----

[![GoDoc][1]][2] [![License: MIT][3]][4] [![Release][5]][6] [![Build Status][7]][8] [![Co decov Coverage][11]][12] [![Go Report Card][13]][14] [![Code Climate][19]][20] [![BCH compliance][21]][22] [![Downloads][15]][16]

[1]: https://godoc.org/github.com/eure/bobo?status.svg
[2]: https://godoc.org/github.com/eure/bobo
[3]: https://img.shields.io/badge/License-MIT-blue.svg
[4]: LICENSE.md
[5]: https://img.shields.io/github/release/eure/bobo.svg
[6]: https://github.com/eure/bobo/releases/latest
[7]: https://travis-ci.org/eure/bobo.svg?branch=master
[8]: https://travis-ci.org/eure/bobo
[9]: https://coveralls.io/repos/eure/bobo/badge.svg?branch=master&service=github
[10]: https://coveralls.io/github/eure/bobo?branch=master
[11]: https://codecov.io/github/eure/bobo/coverage.svg?branch=master
[12]: https://codecov.io/github/eure/bobo?branch=master
[13]: https://goreportcard.com/badge/github.com/eure/bobo
[14]: https://goreportcard.com/report/github.com/eure/bobo
[15]: https://img.shields.io/github/downloads/eure/bobo/total.svg?maxAge=1800
[16]: https://github.com/eure/bobo/releases
[17]: https://img.shields.io/github/stars/eure/bobo.svg
[18]: https://github.com/eure/bobo/stargazers
[19]: https://codeclimate.com/github/eure/bobo/badges/gpa.svg
[20]: https://codeclimate.com/github/eure/bobo
[21]: https://bettercodehub.com/edge/badge/eure/bobo?branch=master
[22]: https://bettercodehub.com/



Slack Bot template for Golang.


# Install

```bash
$ go get -u github.com/eure/bobo
```

# Build

```bash
$ make build
```

for Raspberry Pi

```bash
$ make build-arm6
```

# Run

```bash
SLACK_RTM_TOKEN=xoxb-0000... ./bin/bobo
```

## Environment variables

|Name|Description|
|:--|:--|
| `SLACK_RTM_TOKEN` | [Slack Bot Token](https://slack.com/apps/A0F7YS25R-bots) |
| `SLACK_BOT_TOKEN` | [Slack Bot Token](https://slack.com/apps/A0F7YS25R-bots) |
| `SLACK_TOKEN` | [Slack Bot Token](https://slack.com/apps/A0F7YS25R-bots) |
| `BOBO_DEBUG` | Flag for debug logging. Set [boolean like value](https://golang.org/pkg/strconv/#ParseBool). |

# How to build your original bot

At first, create your own command.

```go
import (
	"github.com/eure/bobo/command"
)

// EchoCommand is an example command.
// This command says same text.
var EchoCommand = command.BasicCommandTemplate{
	Help:           "reply same text",
	MentionCommand: "echo",
	GenerateFn: func(d command.CommandData) command.Command {
		c := command.Command{}
		if d.TextOther == "" {
			return c
		}

		text := fmt.Sprintf("<@%s> %s", d.SenderID, d.TextOther)
		task := command.NewReplyEngineTask(d.Engine, d.Channel, text)
		c.Add(task)
		return c
	},
}
```

Then create `main.go` and add the command,

```go
package main

import (
	"github.com/eure/bobo"
	"github.com/eure/bobo/command"
	"github.com/eure/bobo/engine/slack"
	"github.com/eure/bobo/log"
)

// Entry Point
func main() {
	bobo.Run(bobo.RunOption{
		Engine: &slack.SlackEngine{},
		Logger: &log.StdLogger{
			IsDebug: bobo.IsDebug(),
		},
		CommandSet: command.NewCommandSet(
			// defalt example commands
			command.PingCommand,
			command.HelpCommand,
			// add your original commands
			EchoCommand,
		),
	})
}
```

And run it with Slack Token,

```bash
SLACK_RTM_TOKEN=xoxb-0000... go run ./main.go
```

## Supported tasks

- Slack
    - Reply message
    - Reply message as a thread
    - Add reaction
    - Upload file
- [GoogleHome](https://github.com/eure/bobo-googlehome)

## Experimental Commands

- [evalphobia/bobo-experiment](https://github.com/evalphobia/bobo-experiment)
