package command

import (
	"regexp"

	"github.com/eure/bobo/log"
)

type CommandTemplate interface {
	Exec(CommandData)

	GetMentionCommand() string
	GetRegexp() *regexp.Regexp
	GetHelp() string
	HasHelp() bool
}

// BasicCommandTemplate is a basic command template.
type BasicCommandTemplate struct {
	MentionCommand string
	Regexp         *regexp.Regexp
	GenerateFn     func(CommandData) Command

	Help   string // usage text
	NoHelp bool   // don't show usage text
}

func (c BasicCommandTemplate) Exec(d CommandData) {
	comm := c.GenerateFn(d)
	comm.Exec()
}

func (c BasicCommandTemplate) GetMentionCommand() string {
	return c.MentionCommand
}

func (c BasicCommandTemplate) GetRegexp() *regexp.Regexp {
	return c.Regexp
}

func (c BasicCommandTemplate) GetHelp() string {
	return c.Help
}

func (c BasicCommandTemplate) HasHelp() bool {
	return !c.NoHelp
}

// Command contains task list to execute.
type Command struct {
	tasks  []Task
	logger log.Logger
}

// Add adds a task.
func (c *Command) Add(t Task) {
	c.tasks = append(c.tasks, t)
}

// Exec executes all of tasks.
func (c Command) Exec() {
	for _, task := range c.tasks {
		err := task.Run()
		if err != nil && c.logger != nil {
			c.logger.Errorf("[ERROR] task=[%s], error=%s", task.GetName(), err.Error())
			return
		}
	}
}
