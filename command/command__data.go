package command

import (
	"strings"

	"github.com/eure/bobo/engine"
	"github.com/eure/bobo/library"
)

// CommandData is used to create and execute command.
type CommandData struct {
	Engine engine.Engine

	RawText     string // "@bot <command> <other>" or "<command> <other>"
	Text        string // "@bot <command> <other>"
	TextMention string // @bot
	TextCommand string // <command>
	TextOther   string // <other>
	IsDM        bool   // it's on DM or not

	BotID      string
	SenderID   string
	SenderName string

	// Engine's data
	Channel         string
	ThreadTimestamp string

	// for help usecase
	CommandSet *CommandSet
}

func (d CommandData) HasMyMention() bool {
	text := library.TrimSigns(d.TextMention)
	if !strings.HasPrefix(text, "@") {
		return false
	}

	text = strings.TrimPrefix(text, "@")
	return text == d.BotID
}

func (d CommandData) GetMentionMap() map[string]CommandTemplate {
	return d.CommandSet.mentionMap
}

func (d CommandData) GetRegexpList() []CommandTemplate {
	return d.CommandSet.regexpList
}
