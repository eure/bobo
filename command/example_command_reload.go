package command

// ReloadCommand is an example command.
// This command reloads bot engine.
var ReloadCommand = BasicCommandTemplate{
	Help:           "Reload Bot Engine",
	MentionCommand: "reload",
	GenerateFn: func(d CommandData) Command {
		c := Command{}
		c.Add(NewReplyEngineTask(d.Engine, d.Channel, "Reloading..."))
		c.Add(NewReloadEngineTask(d.Engine))
		return c
	},
}
