package command

import (
	"github.com/eure/bobo/engine"
)

type replyEngineTask struct {
	engine  engine.Engine
	channel string
	text    string
}

// NewReplyEngineTask is a task to reply text.
func NewReplyEngineTask(e engine.Engine, channel, text string) *replyEngineTask {
	return &replyEngineTask{
		engine:  e,
		channel: channel,
		text:    text,
	}
}

func (replyEngineTask) GetName() string {
	return "reply_engine_task"
}

func (t replyEngineTask) Run() error {
	return t.engine.Reply(t.channel, t.text)
}
