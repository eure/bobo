package command

import (
	"github.com/eure/bobo/engine"
)

type reactionEmojiEngineTask struct {
	engine    engine.Engine
	channel   string
	emoji     string
	timestamp string
}

// NewReactionEmojiEngineTask is a task to add reaction to a message.
func NewReactionEmojiEngineTask(e engine.Engine, channel, emoji, timestamp string) *reactionEmojiEngineTask {
	return &reactionEmojiEngineTask{
		engine:    e,
		channel:   channel,
		emoji:     emoji,
		timestamp: timestamp,
	}
}

func (reactionEmojiEngineTask) GetName() string {
	return "reaction_emoji_engine_task"
}

func (t reactionEmojiEngineTask) Run() error {
	return t.engine.ReactEmoji(t.channel, t.emoji, t.timestamp)
}
