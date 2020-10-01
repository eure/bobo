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

	IsFile bool // it's file upload event or not
	File   File

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

// File contains metadeta of a file uploaded by user.
type File struct {
	ID                string
	Name              string
	Title             string
	Mimetype          string
	ImageExifRotation int
	Filetype          string
	PrettyType        string
	Size              int
	URL               string
	IsPublic          bool
	Permalink         string
}
