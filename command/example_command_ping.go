package command

// PingCommand is an example command.
// This command sais "PONG".
var PingCommand = BasicCommandTemplate{
	Help:           "reply PONG",
	MentionCommand: "ping",
	GenerateFn: func(d CommandData) Command {
		c := Command{}
		task := NewReplyEngineTask(d.Engine, d.Channel, "PONG")
		c.Add(task)
		return c
	},
}
