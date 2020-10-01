package slack

import (
	"io"

	"github.com/slack-go/slack"

	"github.com/eure/bobo/engine"
)

// ReplyThread sends a message in the channel.
func (e *SlackEngine) Reply(channel, text string) error {
	reply := e.slackRTM.NewOutgoingMessage(text, channel)
	e.slackRTM.SendMessage(reply)
	return nil
}

// ReplyThread sends a message in the thread of threadTimestamp.
func (e *SlackEngine) ReplyThread(channel, text, threadTimestamp string) error {
	reply := e.slackRTM.NewOutgoingMessage(text, channel)
	reply.ThreadTimestamp = threadTimestamp
	e.slackRTM.SendMessage(reply)
	return nil
}

// ReactEmoji adds reaction emoji to a message.
func (e *SlackEngine) ReactEmoji(channel, emoji, msgTimestamp string) error {
	return e.slackClient.AddReaction(emoji, slack.ItemRef{
		Channel:   channel,
		Timestamp: msgTimestamp,
	})
}

// FileUpload uploads file.
func (e *SlackEngine) FileUpload(channel string, file io.Reader, filename string) error {
	return e.FileUploadWithType(channel, file, filename, "")
}

// FileUploadWithType uploads file with file type.
// erf: https://api.slack.com/types/file#file_types
func (e *SlackEngine) FileUploadWithType(channel string, file io.Reader, filename, filetype string) error {
	_, err := e.slackClient.UploadFile(slack.FileUploadParameters{
		Reader:   file,
		Filename: filename,
		Filetype: filetype,
		Channels: []string{channel},
	})
	return err
}

// GetUserByID gets user by given userID.
func (e *SlackEngine) GetUserByID(userID string) (engine.User, error) {
	if !e.hasUser(userID) {
		err := e.fetchUser(userID)
		if err != nil {
			return engine.User{}, err
		}
	}

	return e.getUser(userID), nil
}

// GetEmojiByRandom gets emoji randomly.
func (e *SlackEngine) GetEmojiByRandom() (string, error) {
	if len(e.emojiCache) == 0 {
		err := e.fetchAllEmoji()
		if err != nil {
			return "", err
		}
	}

	return e.getEmojiByRandom(), nil
}
