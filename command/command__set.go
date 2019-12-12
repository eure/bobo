package command

// CommandSet is a collection of bot commands.
type CommandSet struct {
	List       []CommandTemplate
	mentionMap map[string]CommandTemplate
	regexpList []CommandTemplate
}

func NewCommandSet(templateList ...CommandTemplate) *CommandSet {
	s := &CommandSet{
		List: templateList,
	}
	s.Init()
	return s
}

func (s *CommandSet) Init() {
	s.mentionMap = make(map[string]CommandTemplate, len(s.List))
	s.regexpList = make([]CommandTemplate, 0, len(s.List))
	for _, c := range s.List {
		switch {
		case c.GetMentionCommand() != "":
			s.mentionMap[c.GetMentionCommand()] = c
		case c.GetRegexp() != nil:
			s.regexpList = append(s.regexpList, c)
		}
	}
}

func (s *CommandSet) Exec(d CommandData) {
	d.CommandSet = s

	if d.HasMyMention() {
		if comm, ok := s.getCommandForMention(d); ok {
			comm.Exec(d)
			return
		}
	}

	text := d.RawText
	for _, comm := range s.regexpList {
		if comm.GetRegexp().MatchString(text) {
			comm.Exec(d)
		}
	}
}

// functory method of a command for "@" mention.
func (s *CommandSet) getCommandForMention(d CommandData) (comm CommandTemplate, ok bool) {
	comm, ok = s.mentionMap[d.TextCommand]
	return comm, ok
}
