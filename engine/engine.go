package engine

import "io"

// Engine is a bot engine for a platform.
// (e.g. Slack).
type Engine interface {
	Init(conf Config) error
	Run() (errCode int)
	Reload()
	Close(errCode int)

	// for command and task
	GetUserByID(userID string) (User, error)
	GetEmojiByRandom() (string, error)
	Reply(channel, text string) error
	ReplyThread(channel, text, threadTimestamp string) error
	ReactEmoji(channel, emoji, threadTimestamp string) error
	FileUpload(channel string, file io.Reader, filename string) error
	FileUploadWithType(channel string, file io.Reader, filename, filetype string) error
}

type Config interface{}

// User is a user existed in any platform.
type User struct {
	ID    string
	Name  string
	Email string
	Phone string
}
