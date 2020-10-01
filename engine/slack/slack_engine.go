package slack

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/slack-go/slack"

	"github.com/eure/bobo/command"
	"github.com/eure/bobo/engine"
	"github.com/eure/bobo/errorcode"
	"github.com/eure/bobo/library"
	"github.com/eure/bobo/log"
)

// SlackEngine is a client for Slack.
type SlackEngine struct {
	logger   log.Logger
	commands *command.CommandSet

	slackClose      chan struct{}
	slackClient     *slack.Client
	slackRTM        *slack.RTM
	slackBotID      string
	slackBotName    string
	slackBotMention string // "<@SlackBotID>"

	usersCacheMu sync.RWMutex
	usersCache   map[string]engine.User // key=UserID
	emojiCacheMu sync.RWMutex
	emojiCache   []string

	execChan   chan command.CommandData
	closeChan  chan int
	reloadChan chan struct{}
}

// Init initializes slack engine with Config.
func (e *SlackEngine) Init(conf engine.Config) error {
	c, ok := conf.(Config)
	if !ok {
		return errors.New("Incompatible config type for SlackEngine")
	}

	var opts []slack.Option
	logger := c.GetLogger()
	if logger != nil {
		opts = append(opts, slack.OptionLog(SlackLogger{logger}))
	}

	cli := slack.New(getToken(c), opts...)

	// Test connection to Slack using Token.
	resp, err := cli.AuthTest()
	if err != nil {
		return err
	}
	e.logger = logger
	e.usersCache = make(map[string]engine.User)
	e.commands = c.GetCommandSet()
	e.slackClient = cli
	e.slackRTM = cli.NewRTM()
	e.slackBotID = resp.UserID
	e.slackBotName = resp.User
	e.slackBotMention = fmt.Sprintf("<@%s>", e.slackBotID)
	e.logDebug("slackBotID: [%s]", e.slackBotID)
	e.logDebug("slackBotName: [%s]", e.slackBotName)
	e.execChan = make(chan command.CommandData, c.GetMaxRunningCommands())
	e.slackClose = make(chan struct{}, 1)
	e.closeChan = make(chan int, 1)
	e.reloadChan = make(chan struct{}, 1)
	return nil
}

// KeepAlive keeps the connection in case of disconnection for long-running process.
func (e *SlackEngine) KeepAlive() {
	go e.keepAliveSlack()
}

func (e *SlackEngine) keepAliveSlack() {
	cli := e.slackRTM
	ticker := time.NewTicker(time.Second * 300)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// send empty message.
			cli.SendMessage(cli.NewTypingMessage(""))
		case <-e.slackClose:
			if err := cli.Disconnect(); err != nil {
				e.logDebug("keepAliveSlack slackClose error:[%s]", err.Error())
			}
			return
		}
	}
}

// Close closes the engine.
func (e *SlackEngine) Close(errCode int) {
	e.logDebug("Close")
	e.slackClose <- struct{}{}
	e.closeChan <- errCode
}

// Reload reloads the engine.
func (e *SlackEngine) Reload() {
	e.logDebug("Reload")
	e.reloadChan <- struct{}{}
}

// Run the SlackBot main logic loop.
func (e *SlackEngine) Run() int {
	e.logDebug("Run")
	rtm := e.slackRTM
	go rtm.ManageConnection()
	go e.execCommandDaemon()

	for {
		select {
		case errCode := <-e.closeChan:
			return errCode
		case <-e.reloadChan:
			e.exitNoError()
			return errorcode.None
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				// Event of reaction to user's message
				e.processMessage(ev)
			case *slack.FileSharedEvent:
				// Event of reaction to file upload.
				e.processFileShared(ev)
			case *slack.DisconnectedEvent:
				e.logging("DisconnectedEvent", "intentional=%t", ev.Intentional)
			case *slack.LatencyReport:
				// e.logDebug("LatencyReport: [%d]", ev.Value)
			case *slack.RTMError:
				e.logging("RTMError", "Error:%s", ev.Error())
			case *slack.ConnectionErrorEvent:
				if ev.Error() == "not_authed" && ev.Attempt > 2 {
					e.logging("ConnectionErrorEvent", "Auth Error")
					return errorcode.InvalidAuth
				}
			case *slack.InvalidAuthEvent:
				e.logging("InvalidAuthEvent", "Invalid credentials")
				return errorcode.InvalidAuth
			default:
				// e.logDebug("Unknown Event: [%T] [%+v] ", ev, msg)
			}
		}
	}
}

// processMessage handles message and switch logics depending on user's message.
func (e *SlackEngine) processMessage(ev *slack.MessageEvent) {
	switch {
	case ev.Text == "",
		ev.Hidden,
		e.slackBotID == ev.Msg.User:
		return
	}

	userID := ev.User
	if !e.hasUser(userID) {
		// Get ans save user's name from Slack User ID via API.
		err := e.fetchUser(userID)
		if err != nil {
			return
		}
	}
	// Get user's name from in-memory cache.
	user := e.getUser(userID)

	rawText := library.TrimSpaces(ev.Text)
	text, isDM := e.addMentionIfDM(rawText, ev.Channel)
	e.logDebug("processMessage rawText=[%s] text=[%s] user=[%s | %s]", rawText, text, user.Name, userID)
	mention, commText, otherText := library.SplitTextForCommand(text)

	// wait for executing
	e.execChan <- command.CommandData{
		Engine:          e,
		SenderID:        userID,
		SenderName:      user.Name,
		RawText:         rawText,
		Text:            text,
		TextMention:     mention,
		TextCommand:     commText,
		TextOther:       otherText,
		Channel:         ev.Channel,
		ThreadTimestamp: ev.Timestamp,
		BotID:           e.slackBotID,
		IsDM:            isDM,
	}
}

// processFileShared handles file shared event.
func (e *SlackEngine) processFileShared(ev *slack.FileSharedEvent) {
	if ev.FileID == "" {
		return
	}

	file, err := e.fetchFile(ev.FileID)
	if err != nil {
		return
	}

	userID := file.User
	if !e.hasUser(userID) {
		// Get ans save user's name from Slack User ID via API.
		err := e.fetchUser(userID)
		if err != nil {
			return
		}
	}
	// Get user's name from in-memory cache.
	user := e.getUser(userID)

	channel := ""
	isDM := false
	switch {
	case len(file.Groups) > 0:
		channel = file.Groups[0]
	case len(file.IMs) > 0:
		channel = file.IMs[0]
		isDM = true
	}
	e.logDebug("processFileShared fileID=[%s] fileName=[%s] user=[%s | %s]", file.ID, file.Name, user.Name, userID)

	// wait for executing
	e.execChan <- command.CommandData{
		Engine:          e,
		SenderID:        userID,
		SenderName:      user.Name,
		Channel:         channel,
		ThreadTimestamp: ev.EventTimestamp,
		BotID:           e.slackBotID,
		IsDM:            isDM,
		IsFile:          true,
		File: command.File{
			ID:                file.ID,
			Name:              file.Name,
			Title:             file.Title,
			Mimetype:          file.Mimetype,
			ImageExifRotation: file.ImageExifRotation,
			Filetype:          file.Filetype,
			PrettyType:        file.PrettyType,
			Size:              file.Size,
			URL:               file.URLPrivate,
			IsPublic:          file.IsPublic,
			Permalink:         file.Permalink,
		},
	}
}

// execute a command without blocking.
func (e *SlackEngine) execCommandDaemon() {
	for ch := range e.execChan {
		// execute a command depending on user's message.
		go e.commands.Exec(ch)
	}
}

// getUser gets user from user's name in in-memory cache.
func (e *SlackEngine) getUser(user string) engine.User {
	e.usersCacheMu.RLock()
	defer e.usersCacheMu.RUnlock()
	return e.usersCache[user]
}

// hasUser checks if given user name exists in in-memory cache or not.
func (e *SlackEngine) hasUser(user string) bool {
	e.usersCacheMu.RLock()
	defer e.usersCacheMu.RUnlock()

	_, ok := e.usersCache[user]
	return ok
}

// fetchUser fetches user via Slack API and caches it to in-memory cache.
func (e *SlackEngine) fetchUser(user string) error {
	resp, err := e.slackClient.GetUserInfo(user)
	switch {
	case err != nil:
		e.logging("fetchUser", "Error on `GetUserInfo`: %s", err.Error())
		return err
	case resp == nil:
		err := errors.New("response is nil")
		e.logging("fetchUser", "Error on `GetUserInfo`: %s", err.Error())
		return err
	}

	e.usersCacheMu.Lock()
	defer e.usersCacheMu.Unlock()
	e.usersCache[user] = engine.User{
		ID:    user,
		Name:  resp.Name,
		Email: resp.Profile.Email,
		Phone: resp.Profile.Phone,
	}
	return nil
}

// fetchFile fetches file meta data via Slack API.
func (e *SlackEngine) fetchFile(fileID string) (slack.File, error) {
	resp, _, _, err := e.slackClient.GetFileInfo(fileID, 0, 0)
	switch {
	case err != nil:
		e.logging("fetchFile", "Error on `GetFileInfo`: %s", err.Error())
		return slack.File{}, err
	case resp == nil:
		err := errors.New("response is nil")
		e.logging("fetchFile", "Error on `GetFileInfo`: %s", err.Error())
		return slack.File{}, err
	}
	return *resp, nil
}

// addMentionIfDM adds bot mention into the message when it's on DM.
func (e *SlackEngine) addMentionIfDM(rawText, channel string) (string, bool) {
	switch {
	case !isDM(channel):
		return rawText, false
	case strings.HasPrefix(rawText, e.slackBotMention):
		return rawText, true
	}

	return fmt.Sprintf("%s %s", e.slackBotMention, rawText), true
}

// isDM checks if the message is "Direct Message" or not.
func isDM(channel string) bool {
	// ref: https://stackoverflow.com/questions/41111227/how-can-a-slack-bot-detect-a-direct-message-vs-a-message-in-a-channel
	return strings.HasPrefix(channel, "D")
}

// getEmojiByRandom gets emoji from in-memory cache.
func (e *SlackEngine) getEmojiByRandom() string {
	e.emojiCacheMu.RLock()
	defer e.emojiCacheMu.RUnlock()

	list := e.emojiCache
	if len(list) == 0 {
		return ""
	}

	return list[rand.Intn(len(list))]
}

// fetch all emoji via Slack API.
func (e *SlackEngine) fetchAllEmoji() error {
	e.emojiCacheMu.Lock()
	defer e.emojiCacheMu.Unlock()

	emojiMap, err := e.slackClient.GetEmoji()
	switch {
	case err != nil:
		e.logging("fetchAllEmoji", "Error on `ListReactions`: %s", err.Error())
		return err
	case len(emojiMap) == 0:
		err := errors.New("response is nil")
		e.logging("fetchAllEmoji", "Error on `ListReactions`: %s", err.Error())
		return err
	}

	emojiList := make([]string, 0, len(emojiMap))
	for key := range emojiMap {
		emojiList = append(emojiList, key)
	}

	e.emojiCache = emojiList
	return nil
}

func (e *SlackEngine) logging(typ, msg string, v ...interface{}) {
	msg = "[" + typ + "]\t" + msg
	e.logger.Infof("Slack", fmt.Sprintf(msg, v...))
}

func (e *SlackEngine) logDebug(msg string, v ...interface{}) {
	switch {
	case len(v) == 0:
		e.logger.Debugf("Slack", msg)
	default:
		e.logger.Debugf("Slack", fmt.Sprintf(msg, v...))
	}
}

func (e *SlackEngine) exitNoError() {
	e.logging("exitNoError", "closing...")
	e.Close(errorcode.None)
}
