package command

import (
	"github.com/eure/bobo/engine"
)

type replyThreadEngineTask struct {
	engine    engine.Engine
	channel   string
	text      string
	timestamp string
}

// NewReplyThreadEngineTask is a task to reply text on a thread.
func NewReplyThreadEngineTask(e engine.Engine, channel, text, timestamp string) *replyThreadEngineTask {
	return &replyThreadEngineTask{
		engine:    e,
		channel:   channel,
		text:      text,
		timestamp: timestamp,
	}
}

func (replyThreadEngineTask) GetName() string {
	return "reply_thread_engine_task"
}

func (t replyThreadEngineTask) Run() error {
	return t.engine.ReplyThread(t.channel, t.text, t.timestamp)
}
